package microClustering

import (
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func init() {
	distanceFunctions = make(map[string]DistanceFunc)
	distanceFunctions["euclidian"] = EuclidianDistance
	distanceFunctions["manhattan"] = ManhattanDistance
	distanceFunctions["chebyshev"] = ChebyshevDistance
	distanceFunctions["minkowski"] = MinkowskiDistance
	distanceFunctions["eisen"] = EisenDistance
	distanceFunctions["mahalanobis"] = MahalanobisDistance
	distanceFunctions["cosinus"] = CosinusSimilarity
	Distance = EuclidianDistance
	distanceName = "euclidian"
}

var (
	distanceFunctions map[string]DistanceFunc
	Distance          DistanceFunc
	distanceName      string

	ChebyshevDistance = func(a, b []float64) float64 {
		max := 0.0
		for i, _ := range a {
			max = math.Max(max, math.Abs(a[i]-b[i]))
		}
		return max
	}

	ManhattanDistance = func(a, b []float64) float64 {
		var (
			s float64
		)
		for i, _ := range a {
			s += math.Abs(a[i] - b[i])
		}
		return s

	}

	EuclidianDistance = func(a, b []float64) float64 {
		var (
			s float64
		)
		for i, _ := range a {
			s += math.Pow(a[i]-b[i], 2)
		}
		return math.Sqrt(s)

	}

	MinkowskiP        float64 = 4
	MinkowskiDistance         = func(a, b []float64) float64 {
		var (
			s float64
		)
		for i, _ := range a {
			s += math.Pow(math.Abs(a[i]-b[i]), MinkowskiP)
		}
		return math.Pow(s, 1/MinkowskiP)
	}

	EisenDistance = func(a, b []float64) float64 {
		var (
			s1, s2, s3 float64
		)
		for i, _ := range a {
			s1 += a[i] * b[i]
			s2 += math.Pow(a[i], 2)
			s3 += math.Pow(b[i], 2)
		}
		if s2*s3 == 0 {
			return 1
		}
		return 1 - math.Abs(s1)/math.Sqrt(s2*s3)
	}

	cholCovariance *mat.Cholesky

	MahalanobisDistance = func(a, b []float64) float64 {
		va := mat.NewVecDense(len(a), a)
		vb := mat.NewVecDense(len(b), b)

		return stat.Mahalanobis(va, vb, cholCovariance)
	}

	CosinusSimilarity = func(a, b []float64) float64 {
		na := norm(a)
		nb := norm(b)
		d := na * nb
		if d == 0 {
			if na == nb { // deux vecteurs nul sont similaires
				return 0.0
			}
			return 1.0 // si un seul vecteur est num --> dissimilaires
		}
		return 1.0 - produitVectoriel(a, b)/d

	}
)

func SetDistanceFunction(name string) {
	distanceName = name
	Distance = distanceFunctions[name]
}

func produitVectoriel(a, b []float64) float64 {
	p := float64(0)

	for i := range a {
		p += a[i] * b[i]
	}

	return p
}

func norm(a []float64) float64 {
	n := float64(0)

	for i := range a {
		n += a[i] * a[i]
	}

	return math.Sqrt(n)
}
