package microClustering

import (
	"math"
	"math/rand"
)

// unitySphere génère un point sur une sphere unitaire à N dimensions
func unitySphere(n int) (point []float64) {
	point = make([]float64, n)

	euclidian := distanceName == "euclidian"

	// génération X1 à Xn entre 0..1 (mean=0 et variance=1)
	sum := 0.0
	for i := range point {
		point[i] = rand.Float64()
		if euclidian {
			sum += math.Pow(point[i], 2)
		} else {
			sum += point[i]
		}
	}
	if euclidian {
		sum = math.Sqrt(sum)
		/*fmt.Println("Euclidian distance, sum=", sum)
		} else {
			fmt.Println("Manhattan distance, sum=", sum)*/
	}
	for i := range point {
		point[i] /= sum
	}
	return point
}

func nSphere(center []float64, radius float64) (point []float64) {
	unity := unitySphere(len(center))
	point = make([]float64, len(center))
	for i := range unity {
		point[i] = center[i] + radius*unity[i]
	}
	//fmt.Println("nSphere : dist(P,C)=", Distance(point, center), " radius=", radius)
	return point
}
