package microClustering

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkManhattan(b *testing.B) {

	// Initialise deux vecteurs de test
	n := 32
	x := make([]float64, n)
	y := make([]float64, n)

	for i := range x {
		x[i] = rand.Float64() * 100.0
		y[i] = rand.Float64() * 100.0
	}

	for i := 0; i < b.N; i++ {
		ManhattanDistance(x, y)
	}
}

func BenchmarkCosinus(b *testing.B) {

	// Initialise deux vecteurs de test
	n := 32
	x := make([]float64, n)
	y := make([]float64, n)

	for i := range x {
		x[i] = rand.Float64() * 100.0
		y[i] = rand.Float64() * 100.0
	}

	for i := 0; i < b.N; i++ {
		CosinusSimilarity(x, y)
	}
}

func BenchmarkEuclidian(b *testing.B) {

	// Initialise deux vecteurs de test
	n := 32
	x := make([]float64, n)
	y := make([]float64, n)

	for i := range x {
		x[i] = rand.Float64() * 100.0
		y[i] = rand.Float64() * 100.0
	}

	for i := 0; i < b.N; i++ {
		EuclidianDistance(x, y)
	}
}

func BenchmarkMinkowski(b *testing.B) {

	// Initialise deux vecteurs de test
	n := 32
	x := make([]float64, n)
	y := make([]float64, n)

	for i := range x {
		x[i] = rand.Float64() * 100.0
		y[i] = rand.Float64() * 100.0
	}

	for i := 0; i < b.N; i++ {
		MinkowskiDistance(x, y)
	}
}

func BenchmarkEisen(b *testing.B) {

	// Initialise deux vecteurs de test
	n := 32
	x := make([]float64, n)
	y := make([]float64, n)

	for i := range x {
		x[i] = rand.Float64() * 100.0
		y[i] = rand.Float64() * 100.0
	}

	//test
	for i := 0; i < b.N; i++ {
		EisenDistance(x, y)
	}
}

type vecteur struct {
	label    string
	data     []float64
	required string
}
type measure struct {
	label    string
	function DistanceFunc
	maxDist  float64
}

func TestShapeAnalysis(t *testing.T) {

	// création de plusieurs vecteurs de référence de 32 bytes
	reference := []vecteur{
		vecteur{label: "incr", data: []float64{15, 11, 15, 13, 18, 17, 15, 14, 16, 20, 30, 33, 29, 37, 38, 40, 52, 47, 55, 68, 72, 69, 73, 77, 71, 80, 66, 74, 71, 69, 73, 80}},
		vecteur{label: "spike", data: []float64{15, 11, 15, 13, 18, 17, 15, 14, 16, 20, 30, 33, 57, 53, 73, 80, 74, 58, 53, 49, 39, 41, 29, 27, 20, 17, 19, 16, 21, 27, 23, 12}},
		vecteur{label: "flat", data: []float64{69, 67, 63, 55, 68, 72, 69, 73, 77, 71, 80, 66, 74, 71, 69, 73, 80, 67, 55, 68, 72, 69, 73, 77, 71, 80, 66, 74, 71, 69, 73, 80}},
		vecteur{label: "dip", data: []float64{85, 89, 85, 87, 82, 83, 85, 86, 84, 80, 70, 67, 43, 47, 27, 20, 26, 42, 47, 51, 61, 59, 71, 73, 80, 83, 81, 84, 79, 73, 77, 88}},
	}

	test := []vecteur{
		vecteur{label: "test increase", required: "incr", data: []float64{17, 10, 12, 9, 14, 15, 18, 13, 19, 19, 25, 30, 34, 32, 41, 36, 51, 49, 57, 69, 68, 68, 72, 78, 75, 82, 67, 79, 70, 71, 70, 85}},
		vecteur{label: "test spike", required: "spike", data: []float64{17, 14, 12, 9, 16, 19, 15, 13, 18, 22, 27, 28, 61, 55, 70, 77, 75, 60, 56, 51, 44, 42, 25, 23, 18, 21, 21, 15, 18, 24, 23, 13}},
		vecteur{label: "test flat", required: "flat", data: []float64{65, 68, 61, 55, 70, 71, 69, 78, 82, 68, 78, 68, 77, 71, 65, 74, 85, 68, 54, 70, 69, 67, 76, 73, 67, 79, 66, 71, 73, 73, 76, 79}},
		vecteur{label: "test dip", required: "dip", data: []float64{82, 86, 90, 84, 83, 88, 83, 84, 86, 82, 74, 68, 47, 47, 25, 17, 23, 40, 44, 51, 58, 56, 67, 77, 79, 85, 84, 88, 77, 72, 82, 85}},
		vecteur{label: "test zeros", required: "", data: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		vecteur{label: "test incr1", required: "incr", data: []float64{19, 19, 25, 30, 34, 32, 41, 36, 51, 49, 57, 69, 68, 68, 72, 78, 75, 82, 67, 79, 70, 71, 70, 85, 75, 82, 67, 79, 70, 71, 70, 85}},
		vecteur{label: "test incr2", required: "incr", data: []float64{17, 10, 12, 9, 14, 15, 18, 13, 17, 10, 12, 9, 14, 15, 18, 13, 19, 19, 25, 30, 34, 32, 41, 36, 51, 49, 57, 69, 68, 68, 72, 78}},
		vecteur{label: "test fincr", required: "incr", data: []float64{1, 2, 1, 2, 3, 2, 1, 3, 5, 1, 2, 3, 1, 4, 2, 60, 57, 49, 57, 69, 68, 68, 72, 78, 75, 82, 67, 79, 70, 71, 70, 85}},
		vecteur{label: "test fincr", required: "incr", data: []float64{5, 1, 2, 3, 1, 4, 2, 60, 57, 49, 57, 69, 68, 68, 72, 78, 75, 82, 67, 79, 70, 71, 70, 85, 75, 82, 67, 79, 70, 71, 70, 85}},
		vecteur{label: "test fincr", required: "incr", data: []float64{1, 2, 1, 2, 3, 2, 1, 3, 1, 2, 1, 2, 3, 2, 1, 3, 5, 1, 2, 3, 1, 4, 2, 60, 57, 49, 57, 69, 68, 68, 72, 78}},
	}

	distances := []measure{
		measure{label: "manhattan", function: ManhattanDistance, maxDist: 1000},
		measure{label: "cosinus", function: CosinusSimilarity, maxDist: 1},
		measure{label: "euclidian", function: EuclidianDistance, maxDist: 200},
		measure{label: "eisen", function: EisenDistance, maxDist: 1},
	}

	for _, m := range distances {
		fmt.Println("Measure : ", m.label, ":")
		// comparatifs de distances pour des vecteurs similaires
		fmt.Printf("\t\t\t%s\t%s\t%s\t%s\n", reference[0].label, reference[1].label, reference[2].label, reference[3].label)
		for _, v := range test {
			d := make([]float64, len(reference))
			dmin := -1.0
			refMin := ""

			for i, ref := range reference {
				d[i] = m.function(v.data, ref.data)
				if d[i] < m.maxDist && (refMin == "" || d[i] < dmin) {
					dmin = d[i]
					refMin = ref.label
				}
			}

			fmt.Printf("%s\t\t%0.2f\t%0.2f\t%0.2f\t%0.2f => %s %v\n", v.label, d[0], d[1], d[2], d[3], refMin, v.required == refMin)
		}
	}
}
