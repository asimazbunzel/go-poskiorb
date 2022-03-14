package orbits

import (
   "io/ioutil"
   "go-orbits/pkg/io"

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
}


// read information on binary and kicks from YAML file
func (b *Binary) parseYAML (filename string) error {
  
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

   // load binary into memory
   var binary Binary
   err := binary.parseYAML(filename)
   if err != nil {
      io.LogError("ORBITS - orbits.go - start", "unable to parse YAML file at start")
   }

   // use CGS for this
   binary.convertoCGS()

   return binary
}


// input should be in Msun / Rsun / Lsun and so on.. here we change it to CGS
func (b *Binary) convertoCGS () {

   b.M1 = b.M1 * Msun
   b.M2 = b.M2 * Msun
   b.Separation = b.Separation * Rsun
   b.Period = b.Period * 24 * 3600.0
   b.MCO = b.MCO * Msun

}
