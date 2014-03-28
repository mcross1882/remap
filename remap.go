// remap Takes a CSV file and a map file then generates
// a output file as defined within the map file
//
// \author  Matthew Cross <matthew@pmg.co>
// \package remap
// \version 1.0
package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "encoding/csv"
)

// Global constants
const DefaultBufferSize  = 4096
const DefaultHeaderCount = 128

// Reads a newline delimited map file that is used when outputing the results
func readMapFile(filename string) (mappedHeaders []string, errorMessage string) {
  mappedHeaders = make([]string, DefaultHeaderCount)
  
  file, err := os.Open(filename)
  if err != nil {
    errorMessage = "Cannot open mapped header file."
    return
  }
  defer file.Close()
  
  reader := bufio.NewReaderSize(file, DefaultBufferSize)
  
  index := 0
  line, isPrefix, err := reader.ReadLine()
  for err == nil && !isPrefix {
    if (index >= DefaultBufferSize) {
      errorMessage = "Maxim amount of headers reached."
      return
    }
    mappedHeaders[index] = string(line)
    index++
    line, isPrefix, err = reader.ReadLine()
  }
  return mappedHeaders[:index], ""
}

// Converts two arrays into one associative array (logical equivelent to PHPs' array_combine() function)
//
// \since  1.0
// \param  mappedHeaders a slice of headers retrieved from the map file
// \param  columns an associative array of columns for the current input file line
// \return an array of column values mapped from mappedHeaders
func combineHeaders(mappedHeaders []string, columns map[string]string) (records []string) {
  records = make([]string, len(mappedHeaders))
  for index, header := range mappedHeaders {
    records[index] = columns[header]
  }
  return
}

// Read an input file and a map file then generate an output file
//
// \since  1.0
// \param  inputName the input filename
// \param  outputName the output filename
// \param  mapName the map filename
// \return void
func readFile(inputName string, outputName string, mapName string) {
  inputFile, err := os.Open(inputName)
  if err != nil {
    panic(err)
  }
  defer inputFile.Close()
  
  outputFile, err := os.Create(outputName)
  if err != nil {
    panic(err)
  }
  defer outputFile.Close()
  
  mappedHeaders, mappingError := readMapFile(mapName)
  if mappingError != "" {
    panic(err)
  }
  
  reader := csv.NewReader(inputFile)
  writer := csv.NewWriter(outputFile)
  
  fileHeaders, err := reader.Read()
  columns := make(map[string]string, len(fileHeaders))
  
  writer.Write(mappedHeaders)
  for err == nil {
    records, err := reader.Read()
    if err != nil {
      return
    }
    
    for index, header := range fileHeaders {
      columns[header] = records[index]
    }
    
    err = writer.Write(combineHeaders(mappedHeaders, columns))
    if err != nil {
      fmt.Println(err)
    }
    writer.Flush()
  }
}

// Main entry point
func main() {
  flag.Parse()
  args := flag.Args()
  if len(args) < 3 {
      fmt.Println("remap [input file] [output file] [map file]")
      os.Exit(1)
  }
  readFile(args[0], args[1], args[2])
  fmt.Println("Finished!")
  os.Exit(0)
}
