
package orbits

import (
	"math"
	
   "go-orbits/pkg/io"
)


// kepler law to get binary separation from orbital period
func PtoA (p float64, m1 float64, m2 float64) float64 {

   return math.Pow(StandardCgrav * (m1 + m2) * math.Pow(p/(2*math.Pi),2), 1/3)

}


// kepler law to get orbital period from binary separation
func AtoP (a float64, m1 float64, m2 float64) float64 {

   return (2*math.Pi) * math.Pow(math.Pow(a,3) / (StandardCgrav * (m1 + m2)),0.5)

}


// input should be in Msun / Rsun / Lsun and so on.. here we change it to CGS
func (b *Binary) ConvertoCGS () {

   if b.LogLevel != "none"{
      io.LogInfo("ORBITS - orbits.go - ConvertCGS", "converting to CGS units")
   }

   b.M1 = b.M1 * Msun
   b.M2 = b.M2 * Msun
   b.Separation = b.Separation * Rsun
   b.Period = b.Period * 24 * 3600.0
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
   b.Period = b.Period / 24 / 3600.0
   b.MCO = b.MCO / Msun

   for k, w := range b.W {
      b.W[k] = w / km2cm
   }

   for k,w := range b.WBounded {
      b.WBounded[k] = w / km2cm
      b.SeparationBounded[k] = b.SeparationBounded[k] / Rsun
      b.PeriodBounded[k] = b.PeriodBounded[k] / 24 / 3600.0
   }

}
