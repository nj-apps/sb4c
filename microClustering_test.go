package microClustering

import (
	"fmt"
	"testing"
)

func TestClustering(t *testing.T) {

	data := [][]float64{{2.0, 2.0}, {1.0, 3.0}, {2.0, 8.0}, {2.0, 9.0}, {3, 8}, {4, 6}, {4, 7}, {4, 9}, {5, 7}, {5, 8}, {5, 9}, {6, 4}, {7, 5}, {9, 4}}

	radius := 2.0
	min := 2

	SetDistanceFunction("manhattan")
	c := NewClusterer(radius, min, 1, 2)

	c.Add(data)
	c.RandomDelete(0.08, 0.01)
	c.PrintMicroClusters()
	//data = c.Generate(10000)
	c.Stats()

}

func TestMarshall(t *testing.T) {

	data := [][]float64{{2.0, 2.0}, {1.0, 3.0}, {2.0, 8.0}, {2.0, 9.0}, {3, 8}, {4, 6}, {4, 7}, {4, 9}, {5, 7}, {5, 8}, {5, 9}, {6, 4}, {7, 5}, {9, 4}}

	radius := 2.0
	min := 2

	SetDistanceFunction("manhattan")
	c := NewClusterer(radius, min, 1, 2)

	c.Add(data)
	c.RandomDelete(0.08, 0.01)
	c.PrintMicroClusters()
	//data = c.Generate(10000)
	c.Stats()

	js, err := c.ToJson()
	fmt.Println("err:", err)
	fmt.Println("js:", string(js))

	c2, err := NewClustererFromJson(js)
	fmt.Println("err:", err)
	c2.PrintMicroClusters()
}
