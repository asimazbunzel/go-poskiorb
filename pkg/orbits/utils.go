
package orbits

import (
	"math"
	
   "github.com/asimazbunzel/go-orbits/pkg/io"
)


// kepler law to get binary separation from orbital period
func PtoA (p float64, m1 float64, m2 float64) float64 {

   return math.Pow(StandardCgrav * (m1 + m2) * math.Pow(p/(2.0*math.Pi),2.0), 1.0/3.0)

}


// kepler law to get orbital period from binary separation
func AtoP (a float64, m1 float64, m2 float64) float64 {

   return (2.0*math.Pi) * math.Pow(math.Pow(a,3.0) / (StandardCgrav * (m1 + m2)),0.5)

}


// input should be in Msun / Rsun / Lsun and so on.. here we change it to CGS
func (b *Binary) ConvertoCGS () {

   if b.LogLevel != "none"{
      io.LogInfo("ORBITS - orbits.go - ConvertCGS", "converting to CGS units")
   }

   b.M1 = b.M1 * Msun
   b.M2 = b.M2 * Msun
   b.Separation = b.Separation * Rsun
   b.Period = b.Period * 24.0 * 3600.0
   b.MCO = b.MCO * Msun

   for k, w := range b.W {
      b.W[k] = w * km2cm
   }

}


// let's go back from CGS to astro units
func (b *Binary) ConvertoAstro () {

   if b.LogLevel != "none"{
      io.LogInfo("ORBITS - orbits.go - ConvertoAstro", "converting to Astro units (Msun, Rsun, etc)")
   }

   b.M1 = b.M1 / Msun

   b.M2 = b.M2 / Msun
   b.Separation = b.Separation / Rsun
   b.Period = b.Period / 24.0 / 3600.0
   b.MCO = b.MCO / Msun

   for k, w := range b.W {
      b.W[k] = w / km2cm
   }

   for k,w := range b.WBounded {
      b.WBounded[k] = w / km2cm
      b.SeparationBounded[k] = b.SeparationBounded[k] / Rsun
      b.PeriodBounded[k] = b.PeriodBounded[k] / 24.0 / 3600.0
   }

   for k, _ := range b.PeriodGrid {
      b.PeriodGrid[k] = b.PeriodGrid[k] / 24.0 / 3600.0
      b.SeparationGrid[k] = b.SeparationGrid[k] / Rsun
   }

}


// linspace function
func LinSpace (xi float64, xf float64, num int) []float64 {

   if num <= 1 {
      io.LogError("ORBITS - orbits.go - LinSpace", "`num` must be greater than 1")
   }

   xstep := (xf - xi) / float64(num-1)
   x := make([]float64, num)
	x[0] = xi
	for k := 1; k < num; k++ {
		x[k] = xi + float64(k) * xstep
	}
   x[num-1] = xf
	
   return x

}


// logspace function
func LogSpace (xi float64, xf float64, num int, base float64) []float64 {

   // first, get power in linspace
   xpower := LinSpace(xi, xf, num)

   // now loop over array and compute its power
   x := make([]float64, num)
   for k := 0; k <= len(xpower)-1; k++ {
      x[k] = math.Pow(base, xpower[k])
   }

   return x
}


// return the number of digits of an integer
func CountDigits (number int) int {

   if number < 10 {
      return 1
   } else {
      return 1 + CountDigits(number / 10)
   }
}
