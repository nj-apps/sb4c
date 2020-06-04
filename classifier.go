package microClustering

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type Classifier struct {
	classes       map[int]*Clusterer //map avec un clusterer par classe. Le label de la classe est obligatoirement un int
	Radius        float64
	threshold     int     //Seuls les µC dont la taille dépasse le seuil de prise en compte seront utilisés pour la génération du jeu de données
	outlier       float64 //"nombre de sigmas, un cluster est considéré comme outlier si weight < µ(weights)-outliers*stddev(weights)"
	labelID       int
	Verbose       int
	zones         int
	CheckOutliers bool // TODO: si un point est un outlier pour l'ensemble des classe alors renvoie une classe -1
}

func NewClassifier(labelId int, radius float64, threshold int, zones int, outlier float64) *Classifier {
	newClassifier := Classifier{}
	newClassifier.classes = make(map[int]*Clusterer)
	newClassifier.labelID = labelId
	newClassifier.zones = zones
	newClassifier.threshold = threshold
	newClassifier.outlier = outlier
	newClassifier.Radius = radius
	return &newClassifier
}

var (
	labelId int
)

//Définition d'un type byLabel permettant de trier les données sur la colonne label
type sortByLabel [][]float64

func (m sortByLabel) Len() int { return len(m) }
func (m sortByLabel) Less(i, j int) bool {
	return m[i][labelId] < m[j][labelId]
}
func (m sortByLabel) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//Fit réalise l'apprentissage des données 'X' dont le label est précisé dans 'Y'
func (c *Classifier) FitXY(X [][]float64, Y []int) error {
	if len(X) != len(Y) {
		return fmt.Errorf("data and label mismatch")
	}

	data := [][]float64{}

	//positionne le libellé sur la dernière colonne
	c.labelID = len(X[0])

	//construit le jeu de données
	for i, v := range X {
		data = append(data, append(v, float64(Y[i])))
	}

	c.Fit(data)

	return nil
}

//Fit réalise l'apprentissage des données 'data' dont le label est en colonnes 'labelId'
func (c *Classifier) Fit(data [][]float64) {
	var (
		classData [][]float64
		stdDev    []float64
		moyenne   []float64
	)
	labelId = c.labelID

	// trie les données sur la colonne label
	m := sortByLabel(data)
	sort.Sort(&m)

	// traite les données par label
	labelStart := 0
	currentLabel := data[0][c.labelID]

	if c.Radius == 0 {
		// calcule les distances moyennes par classe
		for i, v := range data {
			if v[c.labelID] != currentLabel || i == len(data)-1 {
				if i < len(data)-1 {
					classData = data[labelStart:i]
				} else {
					classData = data[labelStart:]
				}
				labelStart = i
				// calcule des statistiques de distance entre points pour estimer automatiquement la bonne taille de cluster

				distances := []float64{}

				for nb := 0; nb < 100; nb++ {
					i := rand.Intn(len(classData))
					min := math.MaxFloat64
					for j := range classData {
						if i != j {
							min = math.Min(min, Distance(classData[i], classData[j]))
						}
					}
					distances = append(distances, min)
				}
				_, mean := Minmax(distances)
				_, std := EcartType(distances)
				moyenne = append(moyenne, mean)
				stdDev = append(stdDev, std)
				currentLabel = v[c.labelID]
			}
		}

		// calcule le rayon en fonction des distances moyennes de chaque classe
		mean, _ := EcartType(moyenne)
		_, maxStdDev := Minmax(stdDev)
		c.Radius = mean + 2*maxStdDev //0.5*minStdDev
		if c.Verbose > 0 {
			fmt.Println("=> radius=", c.Radius, " seuil=", c.threshold)
		}
	}

	// Classification
	labelStart = 0
	currentLabel = data[0][c.labelID]
	for i, v := range data { // Parcours les données
		if v[c.labelID] != currentLabel || i == len(data)-1 { //Fin d'une classe
			if i < len(data)-1 {
				classData = data[labelStart:i]
			} else {
				classData = data[labelStart:]
			}
			labelStart = i
			dataFragment := [][]float64{}

			for _, d := range classData { // parcours le jeu de données
				vector := []float64{}
				for i, v := range d { // parcours les colonnes
					if i != c.labelID { // ajoute toutes les colonnes sauf celle contenant le label
						vector = append(vector, v)
					}
				}
				dataFragment = append(dataFragment, vector)
			}

			cl, exists := c.classes[int(currentLabel)]
			if !exists {
				cl = NewClusterer(c.Radius, c.threshold, c.zones, c.outlier)
				c.classes[int(currentLabel)] = cl
			}
			cl.Add(dataFragment)

			if c.Verbose > 0 {
				fmt.Println(" Fit ", len(dataFragment))
			}
			currentLabel = v[c.labelID]
		}
	}
	//Affiche les statistiques des classes
	if c.Verbose > 1 {
		for key, cl := range c.classes {
			fmt.Println("classe : ", key, " µC=", cl.CountMC())
		}
	}
}

//Knn renvoie les libellés des classes les plus proches en utilisant l'algorithme k-Nearest Neightbors
func (c *Classifier) KNN(x [][]float64, k int) (y []int) {
	y = []int{}
	// traite chaque vecteur de données
	for _, vector := range x {
		nearestNeighbors := neighborList{} // liste des µC les plus proches
		// recherche les k NN de chaque classe
		for key, clusterer := range c.classes {
			nb := clusterer.KNN(vector, k)
			for i := range nb {
				nb[i].class = key
				//nb[i].distance = nb[i].distance /* / float64(nb[i].weight)*/ //pondère la distance en fonction du poids du cluster
				nearestNeighbors = append(nearestNeighbors, nb[i])
			}
		}

		//trie par distance
		sort.Sort(nearestNeighbors)

		// recherche la k ieme distance (plusieurs points peuvent être à la même distance)
		toK := 0
		kId := 0
		lastValue := -1.0
		for i := range nearestNeighbors {
			if nearestNeighbors[i].distance != lastValue {
				lastValue = nearestNeighbors[i].distance
				toK++
			}
			if toK == k {
				kId = i + 1
				break
			}
		}
		// recherche la classe la plus représentée
		nearestClass := 0
		nbMax := -1.0
		for key := range c.classes {
			nb := 0.0

			for _, v := range nearestNeighbors[:kId] {
				if v.class == key {
					nb += 1 //v.distance / float64(v.weight)
				}
			}
			if nb > nbMax {
				nbMax = nb
				nearestClass = key
			}
		}

		y = append(y, nearestClass)
	}
	return y
}
