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
  "strings"
  "encoding/csv"
)

// Global constants
const DefaultBufferSize  = 4096
const DefaultHeaderCount = 128
const DefaultFilterCount = 24

var filterSet = make([]Filter, DefaultFilterCount)
var filterSetIndex = 0

type Filter struct {
  field string
  operation string
  value string
}

func (f *Filter) parseOperation(testValue string) bool {
  switch strings.ToLower(f.operation) {
    case "=":
      return f.value == testValue
      
    case "!=":
      return f.value != testValue
      
    case "<=":
      return f.value <= testValue
      
    case ">=":
      return f.value >= testValue
      
    case "<":
      return f.value < testValue
      
    case ">":
      return f.value > testValue
      
    case "like":
      return strings.Contains(f.value, testValue)
      
    case "notlike":
      return !strings.Contains(f.value, testValue)
  }
  return false
}

func (f *Filter) Apply(columns map[string]string) bool {
  return f.parseOperation(columns[f.field])
}


// Adds a filter to the filterset if one exists
//
// \since  1.0
// \param  line the extracted file line
// \return void
func addFilter(line string) string {
  if fields := strings.Split(line, " "); len(fields) == 3 {
    filterSet[filterSetIndex].field     = fields[0]
    filterSet[filterSetIndex].operation = fields[1]
    filterSet[filterSetIndex].value     = fields[2]
    filterSetIndex++
    return fields[0]
  }
  return line
}

// Checks to see if a line contains filter operations
//
// \since  1.0
// \param  filters the filters to apply
// \param  column the line parsed into columns
// \return bool true if the line passes all the filters
func checkFilters(filters []Filter, columns map[string]string) (bool) {
  for i := 0; i < filterSetIndex; i++ {
    if !filters[i].Apply(columns) {
      return false
    }
  }
  return true
}

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
    mappedHeaders[index] = addFilter(string(line))
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
    
    if checkFilters(filterSet, columns) {
      err = writer.Write(combineHeaders(mappedHeaders, columns))
      if err != nil {
        fmt.Println(err)
      }
      writer.Flush()
    }
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
