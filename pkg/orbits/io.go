package orbits

import (
   "fmt"
	"os"
   "strconv"
	"io/ioutil"
	
   "go-orbits/pkg/io"

	"gopkg.in/yaml.v3"
)


// read information on binary and kicks from YAML file
func (b *Binary) parseYAML (filename string) error {

   // read YAML data file into bytes 
   data, err := ioutil.ReadFile(filename)
   if err != nil {
      io.LogError("ORBITS - orbits.go - ParseYAML", "problem reading YAML file")
   }
   
   return yaml.Unmarshal(data, b)
}


// save kick info to file
func (b *Binary) SaveKicks (filename string) {

   if b.LogLevel != "none"{
      io.LogInfo("ORBITS - orbits.go - SaveKicks", "saving kicks information")
   }

   // create file
   f, err := os.Create(filename)
   if err != nil {
      io.LogError("error writing to file", "open file")
   }

   // remember to close the file
   defer f.Close()

   // header
   column_names := [4]string{"id", "w", "theta", " phi"}
   str := fmt.Sprintf("%20s", column_names[0]) 
   str += fmt.Sprintf("%20s", column_names[1])
   str += fmt.Sprintf("%20s", column_names[2])
   str += fmt.Sprintf("%20s\n", column_names[3])
   _, err = f.WriteString(str)
   if err != nil {
      io.LogError("ORBITS - orbits.go - SaveKicks", "error writing header to file")
   }

   // write rows of different natal kicks
   for k, w := range b.W {
      str := fmt.Sprintf("%20s", strconv.Itoa(k))
      str += fmt.Sprintf("%20s", strconv.FormatFloat(w, 'f', 5, 64))
      str += fmt.Sprintf("%20s",strconv.FormatFloat(b.Theta[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s\n",strconv.FormatFloat(b.Phi[k], 'f', 5, 64))
      _, err := f.WriteString(str)
      if err != nil {
         io.LogError("error writing to file", "write error")
      }
   }

}


// save orbits info to file
func (b *Binary) SaveBoundedOrbits (filename string) {

   if b.LogLevel != "none"{
      io.LogInfo("ORBITS - orbits.go - SaveBoundedOrbits", "saving orbits information")
   }

   // create file
   f, err := os.Create(filename)
   if err != nil {
      io.LogError("error writing to file", "open file")
   }

   // remember to close the file
   defer f.Close()

   // header
   column_names := [7]string{"id", "w", "theta", "phi", "period", "separation", "eccentricity"}
   str := fmt.Sprintf("%20s", column_names[0]) 
   str += fmt.Sprintf("%20s", column_names[1])
   str += fmt.Sprintf("%20s", column_names[2])
   str += fmt.Sprintf("%20s", column_names[3])
   str += fmt.Sprintf("%20s", column_names[4]) 
   str += fmt.Sprintf("%20s", column_names[5])
   str += fmt.Sprintf("%20s\n", column_names[6])
   _, err = f.WriteString(str)
   if err != nil {
      io.LogError("ORBITS - orbits.go - SaveKicks", "error writing header to file")
   }

   // write rows of different natal kicks
   for k, kb := range b.IndexBounded {
      str := fmt.Sprintf("%20s", strconv.Itoa(kb))
      str += fmt.Sprintf("%20s", strconv.FormatFloat(b.WBounded[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s", strconv.FormatFloat(b.ThetaBounded[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s", strconv.FormatFloat(b.PhiBounded[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s", strconv.FormatFloat(b.PeriodBounded[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s",  strconv.FormatFloat(b.SeparationBounded[k], 'f', 5, 64))
      str += fmt.Sprintf("%20s\n",  strconv.FormatFloat(b.EccentricityBounded[k], 'f', 5, 64))
      _, err := f.WriteString(str)
      if err != nil {
         io.LogError("error writing to file", "write error")
      }
   }

}
