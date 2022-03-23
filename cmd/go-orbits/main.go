
// Package main provides ...
package main

import (

   "go-orbits/pkg/io"
   "go-orbits/pkg/orbits"

)

func main () {

   // starting logging message
   io.LogInfo("MAIN - main.go - main", "starting orbits study")

   // get binary configuration previous to kick study
   b := orbits.InitBinary("test/config.yaml")

   // compute kicks
   b.ComputeKicks()

   // use CGS units
   b.ConvertoCGS()

   // orbit configurations after momentum kick
   b.OrbitsAfterKicks(true, false)

   // save to file
   b.SaveKicks("test_kicks.data")


   // go back to astro units
   b.ConvertoAstro()

   // save orbits to a file
   b.SaveBoundedOrbits("test_orbits.data")

   // end of computation
   io.LogInfo("MAIN - main.go - main", "exit code with success")

}
