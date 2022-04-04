package orbits

import (
	"fmt"
	"go-orbits/pkg/io"
	"math"
	"sort"
	"strconv"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// structure with binary configuration
type Binary struct {
   M1 float64 `yaml:"m1"`
   M2 float64 `yaml:"m2"`
   Separation float64 `yaml:"separation"`
   Period float64 `yaml:"period"`
   
   MCO float64 `yaml:"compact_object_mass"`

   KickStrengthDistribution string `yaml:"kick_distribution"`
   KickDirection string `yaml:"kick_direction"`
   ReduceByFallback bool `yaml:"reduce_by_fallback"`

   SigmaStrength float64 `yaml:"kick_sigma"`
   MinKickStrength float64 `yaml:"min_kick_value"`
   MaxKickStrength float64 `yaml:"max_kick_value"`

   MinPhi float64 `yaml:"min_phi"`
   MaxPhi float64 `yaml:"max_phi"`

   MinTheta float64 `yaml:"min_theta"`
   MaxTheta float64 `yaml:"max_theta"`
   
   Seed uint64 `yaml:"seed"`

   NumberOfCases int `yaml:"number_of_cases"`
   
   LogLevel string `yaml:"log_level"`

   StoreKicks bool `yaml:"save_kicks"`
   StoreOrbits bool `yaml:"save_bounded_orbits"`
   StoreGrid bool `yaml:"save_grid_of_orbits"`

   KicksFilename string `yaml:"kicks_filename"`
   BoundedBinariesFilename string `yaml:"bounded_orbits_filename"`
   GridFilename string `yaml:"grid_of_orbits_filename"`

   PQuantileMin float64 `yaml:"period_quantile_min"`
   PQuantileMax float64 `yaml:"period_quantile_max"`
   EQuantileMin float64 `yaml:"eccentricity_quantile_min"`
   EQuantileMax float64 `yaml:"eccentricity_quantile_max"`
   PNum int `yaml:"number_of_periods"`
   ENum int `yaml:"number_of_eccentricities"`
   MinProb float64 `yaml:"minimum_probability_for_grid"`

   W []float64
   Phi []float64
   Theta []float64

   IndexBounded []int
   WBounded []float64
   ThetaBounded []float64
   PhiBounded []float64
   SeparationBounded []float64
   EccentricityBounded []float64
   PeriodBounded []float64

   PeriodGrid []float64
   SeparationGrid []float64
   EccentricityGrid []float64
   ProbabilityGrid []float64

}


// initialize structure Binary with the info from a binary system that will then be analyze in
// different conditions due to asymmetric momentum kicks
// it returns the Binary object
func InitBinary (filename string) Binary {

   // load binary into memory
   var binary Binary
   err := binary.parseYAML(filename)
   if err != nil {
      io.LogError("ORBITS - orbits.go - InitBinary", "unable to parse YAML file at start")
   }

   return binary
}


// create slices of asymmetric kicks following a given probability density function
func (b *Binary) ComputeKicks () {

   if b.LogLevel != "none" {
      io.LogInfo("ORBITS - orbits.go - ComputeKicks", "computing momentum kicks")
   }

   // random seed
   src := rand.New(rand.NewSource(b.Seed))

   // Strength of kick based on config option
   if b.KickStrengthDistribution == "Maxwell" {
      // Maxwell distribution is just a chi-squared distribution with 3 d.o.f., k=3
      // therefore, just use inverse sampling for the chi-squared and then correct values with
      // normalization constant
      maxwell := distuv.ChiSquared{3, src}
      for k := 0; k < b.NumberOfCases; k++ {
         b.W = append(b.W, b.SigmaStrength * math.Sqrt(maxwell.Rand()))
      }
   } else if b.KickStrengthDistribution == "Uniform" {
      // Uniform distribution needs min & max values as input
      uniform := distuv.Uniform{b.MinKickStrength, b.MaxKickStrength, src}
      for k := 0; k < b.NumberOfCases; k++ {
         b.W = append(b.W, uniform.Rand())
      }
   } else {
      io.LogError("ORBITS - orbits.go - ComputeKicks", "unknown KickStrengthDistribution")
   }

   // Direction of kicks
   if b.KickDirection == "Uniform" {
      // phi distribution must be between 0 and 2pi
      uniform_phi := distuv.Uniform{b.MinPhi * math.Pi, b.MaxPhi * math.Pi, src}
      for k := 0; k < b.NumberOfCases; k++ {
         b.Phi = append(b.Phi, uniform_phi.Rand())
      }

      // theta distribution must be between 0 and pi, but remember that is modulated by cosine
      uniform_theta := distuv.Uniform{0, 1, src}
      for k := 0; k < b.NumberOfCases; k++ {
         b.Theta = append(b.Theta, math.Acos(2.0 * uniform_theta.Rand() - 1.0))
      }
   } else {
      io.LogError("ORBITS - orbits.go - ComputeKicks", "unknown KickDirection")
   }

   if b.LogLevel == "debug" {
      last_index := 0
      for k, _ := range b.PeriodGrid {
         last_index = k
      }
      digits := CountDigits(last_index)
      fmt.Printf("  id      w   theta   phi\n")
      for k := 0; k < b.NumberOfCases; k++ {
         fmt.Printf("  %0*d    %.2E     %.2E       %.2E\n", digits, k, b.W[k], b.Theta[k], b.Phi[k])
      }
   }

}


// compute orbital parameters assuming linear momentum conservation before and just after
// a momentum kick using Kalogera 1996
func (b *Binary) OrbitsAfterKicks () {

   if b.LogLevel != "none" {
      msg := "calculating post core-collapse orbits for: " + strconv.Itoa(b.NumberOfCases) + " kicks"
      io.LogInfo("ORBITS - orbits.go - OrbitAfterKicks", msg)
   }

   // velocity pre-SN
   vPre := math.Sqrt(StandardCgrav * (b.M1 + b.M2) / b.Separation)
   
   for k := 0; k < b.NumberOfCases; k++ {

      // kick velocity projected to (x,y,z)
      // wx := b.W[k] * math.Cos(b.Phi[k]) * math.Sin(b.Theta[k])
      wy := b.W[k] * math.Cos(b.Theta[k])
      wz := b.W[k] * math.Sin(b.Phi[k]) * math.Sin(b.Theta[k])

      // eqs (3), (4) & (5)
      apost := StandardCgrav * (b.MCO + b.M2) / (2.0 * StandardCgrav * (b.MCO + b.M2) / b.Separation - math.Pow(b.W[k],2.0) - math.Pow(vPre,2.0) - 2.0 * wy * vPre)
      epost := math.Sqrt(1.0 - (math.Pow(wz,2.0) + math.Pow(wy,2.0) + math.Pow(vPre,2.0) + 2.0 * wy * vPre) * math.Pow(b.Separation,2.0) / (StandardCgrav * (b.MCO + b.M2) * apost))

      if epost < 0 || epost > 1 {
         if b.LogLevel == "debug" {
            fmt.Printf("unbounded binary for case: id=%d, w=%.2E, theta=%.2f, phi=%.2f, a=%.2E, e=%.2f\n", k, b.W[k]/1e5, b.Theta[k], b.Phi[k], apost/Rsun, epost)
         }
      } else {

         b.IndexBounded = append(b.IndexBounded, k)
         b.WBounded = append(b.WBounded, b.W[k])
         b.ThetaBounded = append(b.ThetaBounded , b.Theta[k])
         b.PhiBounded = append(b.PhiBounded, b.Phi[k])

         b.SeparationBounded = append(b.SeparationBounded, apost)
         b.EccentricityBounded = append(b.EccentricityBounded, epost)
         // kepler needed here
         b.PeriodBounded = append(b.PeriodBounded, AtoP(apost, b.M1, b.M2))
         
         // if here, binary is bounded after momentum kick
         if b.LogLevel == "debug" {
            fmt.Printf("  bounded binary for case: id=%d, w=%.2E, theta=%.2f, phi=%.2f, a=%.2E, p=%.2E, e=%.2f\n", k, b.W[k]/1e5, b.Theta[k], b.Phi[k], apost/Rsun, AtoP(apost, b.M1, b.M2)/24.0/3600.0, epost)
         }
      }
   }


   if b.LogLevel == "info" || b.LogLevel == "debug" {
      nbounded := len(b.IndexBounded)
      nunbounded := b.NumberOfCases - len(b.IndexBounded)
      fmt.Println("\nSummary of momentum kicks:")
      fmt.Println("number of kicks:", b.NumberOfCases)
      fmt.Printf("fraction of binaries bounded: %d/%d (%f%%)\n", nbounded, b.NumberOfCases, 100*float64(nbounded)/float64(b.NumberOfCases))
      fmt.Printf("fraction of binaries unbounded: %d/%d (%f%%)\n\n", nunbounded, b.NumberOfCases, 100*float64(nunbounded)/float64(b.NumberOfCases))
   }

}


// divide orbital parameter in a grid
func (b *Binary) GridOfOrbits () {

   if b.LogLevel != "none" {
      msg := "calculating grid of orbits for: " + strconv.Itoa(len(b.IndexBounded)) + " cases"
      io.LogInfo("ORBITS - orbits.go - GridOfOrbits", msg)
   }

   // temporary arrays, stat.Quantile needs sorted arrays
   x := make([]float64, len(b.IndexBounded))
   y := make([]float64, len(b.IndexBounded))
   for k, _ := range b.IndexBounded {
      x[k] = b.PeriodBounded[k]
      y[k] = b.EccentricityBounded[k]
   }
   sort.Float64s(x)
   sort.Float64s(y)
   
   // find quantiles according to limits given
   pMin := stat.Quantile(b.PQuantileMin, 1, x, nil)
   pMax := stat.Quantile(b.PQuantileMax, 1, x, nil)
   eMin := stat.Quantile(b.EQuantileMin, 1, y, nil)
   eMax := stat.Quantile(b.EQuantileMax, 1, y, nil)
   
   
   if b.LogLevel != "none" {
      fmt.Println("\nGrid of orbits")
      fmt.Printf("period quantiles: %.2E, %.2E\n", pMin/24.0/3600.0, pMax/24.0/3600.0)
      fmt.Printf("eccentricity quantiles: %.2f, %.2f\n", eMin, eMax)
   }

   // borders in grid
   pBorders := LogSpace(math.Log10(pMin), math.Log10(pMax), b.PNum, 10.0)
   eBorders := LinSpace(eMin, eMax, b.ENum)

   // make grid using borders
   pGrid := make([]float64, b.PNum-1)
   for k := 1; k < len(pBorders); k++ {
      pGrid[k-1] = math.Sqrt(pBorders[k-1] * pBorders[k])
   }
   
   eGrid := make([]float64, b.ENum-1)
   for k := 1; k < len(eBorders); k++ {
      eGrid[k-1] = 0.5 * (eBorders[k-1] + eBorders[k])
   }

   // compute 2D-grid of probabilities
   nRows := len(eGrid)
   nCols := len(pGrid)
   probabilities := make([][]float64, nRows)
   for i := 0; i < nRows; i++ {
      probabilities[i] = make([]float64, nCols)
      for j := 0; j < nCols; j++ {
      probabilities[i][j] = 0.0
      }
   }

   // loop over each binary bounded after kick
   if b.LogLevel == "debug" {
      io.LogInfo("ORBITS - orbits.go - GridOfOrbits", "start loop over random binaries")
   }
   for k := 0; k < len(b.IndexBounded); k++ {
      // temporary vars
      p := b.PeriodBounded[k]
      e := b.EccentricityBounded[k]
      for i := 0; i < nRows; i++ {
         if e >= eBorders[i] && e < eBorders[i+1] {
            for j:= 0; j < nCols; j++ {
               if p >= pBorders[j] && e < pBorders[j+1] {
                  probabilities[i][j] += 1 / float64(len(b.IndexBounded))
                  // if b.LogLevel == "debug" {
                     // fmt.Println("lower < period < upper", pBorders[j]/24.0/3600.0, p/24.0/3600.0, pBorders[j+1]/24.0/3600.0)
                     // fmt.Println("lower < eccentricity < upper", eBorders[i], e, eBorders[i+1])
                  // }
               }
            }
         }
      }
   }
   // some more output for debugging mode
   if b.LogLevel == "debug" {
      for i := 0; i < nRows; i++ {
         fmt.Println("row, probability row:", i, probabilities[i])
      }
   }

   // now get values from grid that are above a minimum probability value
   for i := 0; i < nRows; i++ {
      for j:= 0; j < nCols; j++ {
         if probabilities[i][j] > b.MinProb {
            b.PeriodGrid = append(b.PeriodGrid, pGrid[j])
            b.EccentricityGrid = append(b.EccentricityGrid, eGrid[i])
            b.SeparationGrid = append(b.SeparationGrid, PtoA(pGrid[j], b.M1, b.M2))
            b.ProbabilityGrid = append(b.ProbabilityGrid, probabilities[i][j])
         }
      }
   }
   // output grid above probability minimum
   if b.LogLevel != "none" {
      fmt.Println("\nGrid of orbits above minimum probability")
      fmt.Printf("  id      period   separation   eccentricity\n")
      last_index := 0
      for k, _ := range b.PeriodGrid {
         last_index = k
      }
      digits := CountDigits(last_index)
      for k, _ := range b.PeriodGrid {
         fmt.Printf("  %0*d    %.2E     %.2E       %.2E\n", digits, k, b.PeriodGrid[k]/24.0/3600.0, b.SeparationGrid[k] / Rsun, b.EccentricityGrid[k])
      }
      fmt.Printf("\n")
   }

}
