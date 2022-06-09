// Package main provides ...
package main

import (
	"flag"

	"github.com/asimazbunzel/go-orbits/pkg/io"
	"github.com/asimazbunzel/go-orbits/pkg/orbits"
)

func main () {

   // store name of config file from command line argument
   var configFilename string
   flag.StringVar(&configFilename, "config-file", "config.yaml", "Specify name of configuration file")
   flag.StringVar(&configFilename, "C", "config.yaml", "Specify name of configuration file")
   flag.Parse()

   // get binary configuration previous to kick study
   b := orbits.InitBinary(configFilename)
   
   // starting logging message
   if b.LogLevel != "none" {
      io.LogInfo("MAIN - main.go - main", "starting orbits study")
   }

   // compute kicks
   b.ComputeKicks()

   // use CGS units
   b.ConvertoCGS()

   // orbit configurations after momentum kick
   b.OrbitsAfterKicks()

   // compute grid of orbital parameters
   b.GridOfOrbits()

   // go back to astro units
   b.ConvertoAstro()

   // saves to files
   if b.StoreKicks {
      b.SaveKicks(b.KicksFilename)
   }
   if b.StoreOrbits {
      b.SaveBoundedOrbits(b.BoundedBinariesFilename)
   }
   if b.StoreGrid {
      b.SaveGridOrbits(b.GridFilename)
   }

   // end of computation
   if b.LogLevel != "none" {
      io.LogInfo("MAIN - main.go - main", "exit code with success")
   }

}
