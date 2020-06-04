package microClustering

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

/*
  Clustering online

  Principe :
  - Attribution des points de mesure à des microclusters (µC)
  - chaque microcluster est une sphere à N dimensions caractérisée par son centre, un poids (nombre de mesures dans le cluster)
  - Paramètres :
      - rayon d'un microcluster (précision de clustering)
      - seuil de prise en compte
  - Pour l'ajout :
      - identifier un microcluster existant et incrémenter son poids
      - TODO : recalculer le centre du µC pour le centrer au mieux par rapport à l'historique de mesures
      - si la mesure n'appartient à aucun microcluster, création d'un nouveau µC
  - Un pourcentage d'oubli peut être appliqué à intervalles régulier dans le cas de jeu de données live pour adapter la segmentation à
    l'évolution de la population dans le temps. Les µC dont la taille passe en dessous du seuil de prise en compte sont supprimés.
  - Generation de jeu de données : génération d'un jeu de données de taille fixe respectant la typologie des points de mesure injectés.

  - Algorithmes exploitant les µC:
      - kMeans : implémentation du clustering utilisant les µC au lieu des données brutes.
      - kNN : Version µC de l'Algorithme de classification kNN
      - Algo de clustering basé sur l'aggrégation de µC

*/

// DistanceFunc represents a function for measuring distance
// between n-dimensional vectors.
type DistanceFunc func([]float64, []float64) float64

// Une mesure est un vecteur de float64 à N dimensions
//type Measurement []float64

type microcluster struct {
	Center []float64 `json:"center"` // centre
	Zones  []int     `json:"zones"`  // nombre de mesures par zone
	Weight int       `json:"weight"` // poids du cluster

	//kmeanId int       // numéro du clusters en clusterisation kmean

}

func (mc microcluster) String() string {
	return fmt.Sprintf("MC : center=%v zones=%v weight=%d", mc.Center, mc.Zones, mc.Weight)
}

type Clusterer struct {
	//Paramètres
	mcRadius float64 // rayon du cluster
	minSize  int     // nombre minimum d'éléments composant un micro cluster pour qu'il soit pris en compte pour la génération du jeu de données
	zones    int     // nombre de zones concentriques pour la répartition
	// statistiques sur les µC
	mediumSize       float64
	sigmaSize        float64
	outlierThreshold float64 // un µC est considéré comme outlier si Weight < mediumSize-outlierThreshold*sigmaSize

	//structures du cluster
	vectorSize   int
	distFunction string          //nom de la fonction distance utilisée
	distance     DistanceFunc    // fonction utilisée pour évaluer les distances
	mc           []*microcluster // liste de tous les microclusters créés
}

func (c *Clusterer) CountMC() int {
	return len(c.mc)
}

//IsOutlier renvoie true si le point n'appartient a aucun µCluster représentatif
func (c *Clusterer) IsOutlier(x []float64) bool {
	threshold := c.mediumSize - c.outlierThreshold*c.sigmaSize
	for _, mc := range c.mc { // Pour chaque µC
		if float64(mc.Weight) > threshold { // S'il est représentatif
			if c.distance(x, mc.Center) <= c.mcRadius { // tant que le point généré n'est pas dans la sphere
				return false
			}
		}
	}
	return true
}

// Generate génére un jeu de données de 'size' éléments aléatoire respectant la distribution des µC représentatifs
// Le jeu de données généré peut être légèrement plus grand que la taille demandée si la difference de taille entre les plus grands
// et les plus petits clusters est très importante
func (c *Clusterer) Generate(size int) (data [][]float64) {
	totalSize := 0
	//calcule le nombre d'elements
	for _, mc := range c.mc {
		if mc.Weight >= c.minSize {
			totalSize += mc.Weight
		}
	}
	//fmt.Println("Original size :", totalSize)

	rand.Seed(time.Now().UnixNano())

	//génération du jeu de données
	for _, mc := range c.mc { // Pour chaque µC

		//fmt.Println("µC weight : ", mc.Weight, " minSize=", c.minSize)
		if mc.Weight >= c.minSize { // S'il est représentatif

			coeff := float64(mc.Weight) / float64(totalSize)
			nbToGenerate := int(coeff * float64(size)) // calcule le nombre d'éléments à générer pour ce µC
			//fmt.Println("coeff=", coeff, " nbToGenerate=", nbToGenerate, " (", coeff*float64(size), ")")
			if nbToGenerate == 0 { // il doit y avoir au moins un point par µC représentatif
				nbToGenerate = 1
			}
			if nbToGenerate > 0 {
				data = append(data, mc.Generate(nbToGenerate, c.mcRadius, c.distance)...)
			}
		}
	}

	// Si le nombre de points générés est inférieur au nombre de points demandé, ajoute autant de points que nécessaire
	for len(data) < size {
		mcid := rand.Intn(len(c.mc))
		if c.mc[mcid].Weight >= c.minSize {
			data = append(data, c.mc[mcid].Generate(1, c.mcRadius, c.distance)...)
		}
	}

	return data
}

/*func (c *Clusterer) Generate(size int) (data [][]float64) {
	totalSize := 0
	//calcule le nombre d'elements
	for _, mc := range c.mc {
		if mc.Weight >= c.minSize {
			totalSize += mc.Weight
		}
	}
	//fmt.Println("Original size :", totalSize)

	rand.Seed(time.Now().Unix())

	//génération du jeu de données
	for _, mc := range c.mc { // Pour chaque µC
		//fmt.Println("µC weight : ", mc.Weight, " minSize=", c.minSize)
		if mc.Weight >= c.minSize { // S'il est représentatif

			coeff := float64(mc.Weight) / float64(totalSize)
			nbToGenerate := int(coeff * float64(size)) // calcule le nombre d'éléments à générer pour ce µC
			if nbToGenerate == 0 {                     // il doit y avoir au moins un point par µC représentatif
				nbToGenerate = 1
			}

			for i := 0; i < nbToGenerate; i++ {
				vector := []float64{}
				distCentre := c.mcRadius + 1
				for distCentre >= c.mcRadius { // tant que le point généré n'est pas dans la sphere
					vector = []float64{}
					for d := 0; d < c.vectorSize; d++ { // Génère un point aléatoirement
						v := mc.Center[d] - c.mcRadius + 2*c.mcRadius*rand.Float64()
						vector = append(vector, v)
					}
					distCentre = c.distance(vector, mc.Center) // calcule la distance au centre du mc
				}
				data = append(data, vector)
			}
		}
	}

	// Si le nombre de points générés est inférieur au nombre de points demandé, ajoute autant de points que nécessaire
	for len(data) < size {
		mcid := rand.Intn(len(c.mc))
		if c.mc[mcid].Weight >= c.minSize {
			vector := []float64{}
			distCentre := c.mcRadius + 1
			for distCentre >= c.mcRadius { // tant que le point généré n'est pas dans la sphere
				vector = []float64{}
				for d := 0; d < c.vectorSize; d++ { // Génère un point aléatoirement
					v := c.mc[mcid].Center[d] - c.mcRadius + 2*c.mcRadius*rand.Float64()
					vector = append(vector, v)
				}
				distCentre = c.distance(vector, c.mc[mcid].Center) // calcule la distance au centre du mc
			}
			data = append(data, vector)
		}
	}

	return data
}
*/

func NewClusterer(radius float64, minSize int, zones int, outlierThreshold float64) (clusterer *Clusterer) {
	clusterer = new(Clusterer)

	clusterer.mcRadius = radius
	clusterer.minSize = minSize
	clusterer.distFunction = distanceName
	clusterer.distance = Distance
	clusterer.outlierThreshold = outlierThreshold
	clusterer.zones = zones
	return clusterer
}

func (c *Clusterer) Stats() {
	fmt.Println("nb µClusters : ", len(c.mc))

	moy := 0
	min := math.MaxInt64
	max := 0
	for _, m := range c.mc {
		moy += m.Weight
		if m.Weight > max {
			max = m.Weight
		}
		if m.Weight < min {
			min = m.Weight
		}
	}
	fmt.Println("mean weight : ", moy/len(c.mc), " max=", max, " min=", min)
	fmt.Println("weighted radius : moy=", float64(c.mcRadius)*math.Log(float64(moy/len(c.mc))), " max=", float64(c.mcRadius)*math.Log(float64(max)), " min=", float64(c.mcRadius)*math.Log(float64(min)))

}

// recherche un cluster pour chaque point
func (c *Clusterer) Add(m [][]float64) {
	if c.vectorSize == 0 {
		c.vectorSize = len(m[0])
	}
	for i := range m {
		if m[i] == nil {
			continue
		}
		clusterFound := false
		for mc := range c.mc { // recherche dans les cluster nouvellement créés
			distance := c.distance(c.mc[mc].Center, m[i])
			if distance <= c.mcRadius {
				c.mc[mc].add(m[i], distance, c.mcRadius)
				clusterFound = true
				break
			}
		}
		if !clusterFound { //création d'un nouveau microcluster
			newMc := microcluster{Center: m[i], Weight: 1}
			newMc.Zones = make([]int, c.zones)
			newMc.Zones[0] = 1
			c.mc = append(c.mc, &newMc)
		}
	}
	//fmt.Println("MC : ", len(c.mc))
}

//Add ajoute une mesure dans un microcluster
//Add décale la position du centre du cluster vers le nouveau point ajouté en prennant en compte la pondération du µC
// Si le µC ne contient qu'un seul point alors le nouveau centre sera à mi-distance entre le centre actuel et le nouveau point
// par contre si le µC contient déjà 100 points alors le nouveau centre sera 100x plus proche du centre actuel que du nouveau point
func (mc *microcluster) add(m []float64, dist float64, radius float64) {
	for i := range m {
		mc.Center[i] = (float64(mc.Weight)*mc.Center[i] + m[i]) / (float64(mc.Weight) + 1)
	}
	//	fmt.Printf("ADD %v dist=%0.2f ", mc.Zones, dist)
	for z := 0; z < len(mc.Zones); z++ {
		//fmt.Printf("z=%d (r=%0.2f) ", z, (float64(z+1)*radius)/float64(len(mc.Zones)))
		if dist <= (float64(z+1)*radius)/float64(len(mc.Zones)) {
			//		fmt.Print(": z=", z)
			mc.Zones[z]++
			break
		}
	}
	//fmt.Println(mc.Zones)
	mc.Weight++
}

func (c *Clusterer) PrintMicroClusters() {
	for i, mc := range c.mc {
		fmt.Println(i, " - ", mc.Center, " weight=", mc.Weight)
	}
}

// RandomDelete supprime pct mesures dans les mc
// chaque mesure à la probabilité p d'être supprimée
func (c *Clusterer) RandomDelete(pct float64, p float64) {
	// compte les mc
	totalMc := 0
	for i := range c.mc {
		totalMc += c.mc[i].Weight
	}

	// supprime des mesures
	proba := 100.0 * p

	toDelete := int(pct * float64(totalMc))

	//	fmt.Println("random delete : ", toDelete)

	for toDelete > 0 {
		for i := range c.mc {
			if c.mc[i].Weight > 0 { // si le mc contient encore des mesures
				p := rand.Float64() * 100
				if p <= proba { // proba de supprimer une mesure
					c.mc[i].Weight--
					z := 0
					for { // sélectionne aléatoirement la zone dans laquelle supprimer le point
						z = rand.Intn(c.zones)
						if c.mc[i].Zones[z] > 0 {
							break
						}
					}
					c.mc[i].Zones[z]--

					toDelete--
					if toDelete == 0 {
						break
					}
				}
			}
		}
	}

	// supprime les mc vides
	length := len(c.mc)
	nbMCdeleted := 0
	for i := 0; i < length; i++ {
		if c.mc[i].Weight == 0 {
			copy(c.mc[i:], c.mc[i+1:])
			c.mc[len(c.mc)-1] = nil
			c.mc = c.mc[:len(c.mc)-1]
			length = len(c.mc)
			i--
			nbMCdeleted++
		}
	}
}

func (mc *microcluster) generateVector(radius float64, distance DistanceFunc) (result []float64) {
	// Tirage aléatoire de l'ordre de génération des axes
	nb := len(mc.Center)
	X := make([]float64, nb)

	//	dist := distuv.UnitNormal // Distribution

	//fmt.Println("distNormal : ", dist.Rand())

	// génération X1 à Xn entre 0..1 (mean=0 et variance=1)
	for i := 0; i < nb; i++ {
		X[i] = rand.NormFloat64()
	}
	max := math.Abs(X[0])
	for _, x := range X {
		max = math.Max(max, math.Abs(x))
	}
	for i := range X {
		X[i] /= max
	}

	fmt.Println("X=", X)

	// U entre 0..1
	u := rand.Float64()
	fmt.Println("u=", u)

	//Calcule les points
	sum := 0.0
	for _, x := range X {
		sum += math.Pow(x, 2)
	}

	d := math.Sqrt(sum)

	//  point sur la sphere unitaire
	for i := range X {
		X[i] /= d
	}

	fmt.Println("Sur Sphere unitaire : ", X)

	//	coeff := (math.Pow(radius*u, 1/float64(nb))) / math.Sqrt(sum)
	coeff := radius * math.Pow(u, 1/float64(nb))
	fmt.Println("coeff=", coeff)

	result = make([]float64, nb)

	for i := range X {
		result[i] = mc.Center[i] + coeff*X[i]
	}
	//fmt.Println("X=", result)

	//Zero := make([]float64, nb)

	fmt.Printf("RADIUS=%0.001f  Manhattan=%0.001f Euclidien=%0.001f\n", radius, distance(mc.Center, X), EuclidianDistance(mc.Center, X))
	return result
}

func (mc *microcluster) generateVectorOld(radius float64, distance DistanceFunc) (result []float64) {
	// Tirage aléatoire de l'ordre de génération des axes
	nb := len(mc.Center)
	order := []int{}
	for nb > 0 {
		pos := rand.Intn(len(mc.Center))
		found := false
		for _, v := range order {
			if v == pos {
				found = true
				break
			}
		}
		if !found {
			order = append(order, pos)
			nb--
		}
	}
	//	fmt.Println("order = ", order)
	// generation
	result = make([]float64, len(mc.Center))
	copy(result, mc.Center)
	maxLen := radius
	var dist float64
	for i, pos := range order {
		result[pos] = mc.Center[pos] - maxLen + 2*maxLen*rand.Float64()
		//	fmt.Println("maxlen=", maxLen, " --> ", result)
		dist = distance(mc.Center, result)
		//fmt.Print("dist=", dist)
		maxLen = radius - dist
		fmt.Println("i=", i, " dist=", dist, " radius=", radius, " max=", maxLen)
		if maxLen < 0 {
			//fmt.Println("maxlen<0 : return => ", result)
			return result
		}

	}
	//fmt.Println("finish : return => ", result)

	return result
}

//Generate crée nb points aleatoires dans le cluster de rayon "radius" en respectant la répartition dans les zones
func (mc *microcluster) Generate(nb int, radius float64, distance DistanceFunc) (data [][]float64) {
	totalGenerated := 0
	for z, zone := range mc.Zones {

		nbZone := int(math.Round(float64(zone*nb) / float64(mc.Weight)))
		totalGenerated += nbZone
		//calcule le rayon de la zone en fonction du nombre de zones
		radiusZone := radius * float64(z+1) / float64(len(mc.Zones))
		radiusPrevZone := radius * float64(z) / float64(len(mc.Zones))
		for nbZone > 0 {
			r := radiusPrevZone + rand.Float64()*(radiusZone-radiusPrevZone) // génére aléatoirement un rayon dans la zone
			vector := nSphere(mc.Center, r)
			data = append(data, vector)
			nbZone--
		}
	}

	// Complète en rajoutant des points au hazard s'il en manque
	manque := nb - totalGenerated

	for manque > 0 {
		r := rand.Float64() * radius // génére aléatoirement un rayon dans la sphere
		vector := nSphere(mc.Center, r)
		data = append(data, vector)
		manque--
	}

	return data
}

type neighbor struct {
	distance float64
	weight   int
	class    int
}

type neighborList []neighbor

func (m neighborList) Len() int { return len(m) }
func (m neighborList) Less(i, j int) bool {
	return m[i].distance < m[j].distance
}
func (m neighborList) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//KNN renvoie les k µC les plus proches
func (c *Clusterer) KNN(x []float64, k int) (mc []neighbor) {
	var nb neighborList

	// mesure la distance à chaque µc
	for _, mc := range c.mc {

		dist := Distance(x, mc.Center)
		w := float64(mc.Weight)
		newNeighbor := neighbor{distance: dist / w, weight: mc.Weight}
		//newNeighbor := neighbor{distance: Distance(x, mc.Center) / float64(mc.Weight), weight: mc.Weight}
		nb = append(nb, newNeighbor)
	}

	// trie la liste
	sort.Sort(nb)

	if k < len(nb) {
		toK := 0
		lastValue := -1.0
		for i := range nb {
			if nb[i].distance != lastValue {
				lastValue = nb[i].distance
				toK++
			}
			if toK == k {
				return nb[0 : i+1]
			}
		}
		return nb[0:k]
	} else {
		return nb
	}

}

// NNDistance calcul la plus petite distance entre deux mesures de data
func NNDistance(id int, data [][]float64 /*, distance DistanceFunc*/) (min float64) {
	min = math.MaxFloat64

	sample := false
	if len(data) > 5000 {
		sample = true
	}

	j := 0
	for j = id; j == id; j = rand.Intn(len(data)) {

	}

	min = Distance(data[id], data[j])
	for i := range data {
		if i != id {
			if !sample || rand.Float64() < 0.01 {
				d := Distance(data[id], data[i])
				if d < min {
					min = d
				}
			}
		}
	}
	return min
}

// MeanNN calcule la distance moyenne et l'écart type moyen entre deux mesures.
func MeanNN(data [][]float64 /*, distance DistanceFunc*/) (mean, sd float64) {
	var (
		distances []float64
		sample    bool = false
	)
	if len(data) > 5000 {
		sample = true
	}

	if sample { // prends 100 points aléatoirement
		for nb := 0; nb < 100; nb++ {
			i := rand.Intn(len(data))
			min := math.MaxFloat64
			for j := range data {
				if i != j {
					min = math.Min(min, Distance(data[i], data[j]))
				}
			}
			distances = append(distances, min)
		}
		_, maxDist := Minmax(distances)
		_, std := EcartType(distances)
		return maxDist, std
	}

	// Sinon calcule la distance la plus courte en parcourant toutes les combinatoires
	for i := range data {
		distances = append(distances, NNDistance(i, data /*, distance*/))
	}

	mean, sd = EcartType(distances)

	return mean, sd
}

func (c *Clusterer) Size() int {
	totalSize := 0
	//calcule le nombre d'elements
	for _, mc := range c.mc {
		if mc.Weight >= c.minSize {
			totalSize += mc.Weight
		}
	}
	return totalSize
}
