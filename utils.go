package microClustering

import "math"

// ecartType calcule la moyenne, l'écart-type
func EcartType(num []float64) (float64, float64) {

	var sum, mean, sd float64

	for i := 0; i < len(num)-2; i++ {
		sum += num[i]
	}
	mean = sum / float64(len(num))

	for j := 0; j < len(num)-2; j++ {
		// The use of Pow math function func Pow(x, y float64) float64
		sd += math.Pow(num[j]-mean, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	sd = math.Sqrt(sd / float64(len(num)))

	//fmt.Println("The Standard Deviation is : ", sd)
	return mean, sd

}

func Minmax(x []float64) (min float64, max float64) {
	if len(x) > 0 {
		tmpMin := x[0]
		tmpMax := x[0]

		for _, v := range x {
			if v < tmpMin {
				tmpMin = v
			}
			if v > tmpMax {
				tmpMax = v
			}
		}
		return tmpMin, tmpMax
	}

	return 0, 0
}
