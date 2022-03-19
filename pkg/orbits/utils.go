
package orbits

import (
	"math"
)


// kepler law to get binary separation from orbital period
func PtoA (p float64, m1 float64, m2 float64) float64 {

   return math.Pow(StandardCgrav * (m1 + m2) * math.Pow(p/(2*math.Pi),2), 1/3)

}


// kepler law to get orbital period from binary separation
func AtoP (a float64, m1 float64, m2 float64) float64 {

   return (2*math.Pi) * math.Pow(math.Pow(a,3) / (StandardCgrav * (m1 + m2)),0.5)

}
