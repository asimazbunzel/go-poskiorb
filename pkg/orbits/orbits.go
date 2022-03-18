package orbits

import (
	"fmt"
	"go-orbits/pkg/io"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat/distuv"
	"gopkg.in/yaml.v3"
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
   
   NumberOfCases int `yaml:"number_of_cases"`

   W []float64
   Phi []float64
   Theta []float64
}


// read information on binary and kicks from YAML file
func (b *Binary) parseYAML (filename string) error {

   io.LogInfo("ORBITS - orbits.go - parseYAML", "reading configuration from file")
  
   // read YAML data file into bytes 
   data, err := ioutil.ReadFile(filename)
   if err != nil {
      io.LogError("ORBITS - orbits.go - ParseYAML", "problem reading YAML file")
   }
   
   return yaml.Unmarshal(data, b)
}


// initialize structure Binary with the info from a binary system that will then be analyze in
// different conditions due to asymmetric momentum kicks
// it returns the Binary object
func InitBinary (filename string) Binary {

   io.LogInfo("ORBITS - orbits.go - InitBinary", "initializing binary")

   // load binary into memory
   var binary Binary
   err := binary.parseYAML(filename)
   if err != nil {
      io.LogError("ORBITS - orbits.go - InitBinary", "unable to parse YAML file at start")
   }

   return binary
}


// input should be in Msun / Rsun / Lsun and so on.. here we change it to CGS
func (b *Binary) ConvertoCGS () {

   io.LogInfo("ORIBTS - orbits.go - ConvertCGS", "converting to CGS units")

   b.M1 = b.M1 * Msun
   b.M2 = b.M2 * Msun
   b.Separation = b.Separation * Rsun
   b.Period = b.Period * 24 * 3600.0
   b.MCO = b.MCO * Msun

   for k, w := range b.W {
      b.W[k] = w * km2cm
   }

}


// create slices of asymmetric kicks following a given probability density function
func (b *Binary) ComputeKicks () {

   io.LogInfo("ORBITS - orbits.go - ComputeKicks", "computing momentum kicks")

   // Strength of kick based on config option
   if b.KickStrengthDistribution == "Maxwell" {
      // Maxwell distribution is just a chi-squared distribution with 3 d.o.f., k=3
      // therefore, just use inverse sampling for the chi-squared and then correct values with
      // normalization constant
      maxwell := distuv.ChiSquared{3, nil}
      for k := 1; k <= b.NumberOfCases; k++ {
         b.W = append(b.W, b.SigmaStrength * math.Sqrt(maxwell.Rand()))
      }
   } else if b.KickStrengthDistribution == "Uniform" {
      // Uniform distribution needs min & max values as input
      uniform := distuv.Uniform{b.MinKickStrength, b.MaxKickStrength, nil}
      for k := 1; k <= b.NumberOfCases; k++ {
         b.W = append(b.W, uniform.Rand())
      }
   } else {
      io.LogError("ORBITS - orbits.go - ComputeKicks", "unknown KickStrengthDistribution")
   }

   // Direction of kicks
   if b.KickDirection == "Uniform" {
      // phi distribution must be between 0 and 2pi
      uniform_phi := distuv.Uniform{b.MinPhi * math.Phi, b.MaxPhi * math.Phi, nil}
      for k := 1; k <= b.NumberOfCases; k++ {
         b.Phi = append(b.Phi, uniform_phi.Rand())
      }

      // theta distribution must be between 0 and pi, but remember that is modulated by cosine
      uniform_theta := distuv.UnitUniform
      for k := 1; k <= b.NumberOfCases; k++ {
         b.Theta = append(b.Theta, math.Acos(2 * uniform_theta.Rand() - 1))
      }
   } else {
      io.LogError("ORBITS - orbits.go - ComputeKicks", "unknown KickDirection")
   }

}


// compute orbital parameters assuming linear momentum conservation before and just after
// a momentum kick using Kalogera 1996
func (b *Binary) OrbitsAfterKicks (verbose bool) {

   msg := "calculating post core-collapse orbits for: " + strconv.Itoa(b.NumberOfCases) + " kicks"
   io.LogInfo("ORBITS - orbits.go - OrbitAfterKicks", msg)


   // velocity pre-SN
   vPre := math.Sqrt(StandardCgrav * (b.M1 + b.M2) / b.Separation)
   
   for k := 0; k <= b.NumberOfCases-1; k++ {

      // kick velocity projected to (x,y,z)
      // wx := b.W[k] * math.Cos(b.Phi[k]) * math.Sin(b.Theta[k])
      wy := b.W[k] * math.Cos(b.Theta[k])
      wz := b.W[k] * math.Sin(b.Phi[k]) * math.Sin(b.Theta[k])

      // eqs (3), (4) & (5)
      apost := StandardCgrav * (b.MCO + b.M2) / (2 * StandardCgrav * (b.MCO + b.M2) / b.Separation - math.Pow(b.W[k],2) - math.Pow(vPre,2) - 2*wy * vPre)
      epost := math.Sqrt(1 - (math.Pow(wz,2) + math.Pow(wy,2) + math.Pow(vPre,2) + 2*wy*vPre) * math.Pow(b.Separation,2) / (StandardCgrav * (b.MCO + b.M2) * apost))

      if epost < 0 || epost > 1 {
         fmt.Println("unbind binary for case ", k)
      } else {
         fmt.Println("bounded binary for case", k)
      }

   }

}


// save kick info to file
func (b *Binary) SaveKicks (filename string) {

   io.LogInfo("ORBITS - orbits.go - SaveKicks", "saving kicks information")

   // create file
   f, err := os.Create(filename)
   if err != nil {
      io.LogError("error writing to file", "open file")
   }

   // remember to close the file
   defer f.Close()

   // write rows of different natal kicks
   for k, w := range b.W {
      str := strconv.FormatFloat(w, 'f', 5, 64) + " " + strconv.FormatFloat(b.Theta[k], 'f', 5, 64) + " " + strconv.FormatFloat(b.Phi[k], 'f', 5, 64) + "\n"
      _, err := f.WriteString(str)
      if err != nil {
         io.LogError("error writing to file", "write error")
      }
   }


}
