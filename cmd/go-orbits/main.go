
// Package main provides ...
package main

import (
   "fmt"
   "go-orbits/pkg/io"
   "go-orbits/pkg/orbits"
)

func main () {

   // starting logging message
   io.LogInfo("MAIN - main.go - main", "starting orbits study")

   // get binary configuration previous to kick study
   b := orbits.InitBinary("test/config.yaml")

   fmt.Printf("%+v\n",b)

   // end of computation
   io.LogInfo("MAIN - main.go - main", "exit code with success")

}
