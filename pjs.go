package main

import (
  "fmt"
  "flag"
  "os"
)


type optionSet struct {
  color bool
  indent int
  paths []string
}


func parseFlag() (*os.File, optionSet) {
  options := optionSet{}
  noColor := false
  name := "pjs"

  flagset := flag.NewFlagSet(name, flag.ExitOnError)
  flagset.BoolVar(&noColor, "nc", true, "\tDon't output colors")
  flagset.IntVar(&options.indent, "i", 2, "\tIndent size")

  flagset.Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s - Pretty print and manipulate JSON data\n\n[options]\n\n", name)
    flagset.PrintDefaults()
    fmt.Fprintf(os.Stderr, "\n")
    os.Exit(2)
  }

  flagset.Parse(os.Args[1:])

  options.color = !noColor

  var filename string
  processPaths := false

  for _, item := range flagset.Args() {
    if item == "--" {
      processPaths = true
    } else if processPaths {
      options.paths = append(options.paths, item) // TODO: Process paths
    } else if filename == "" {
      filename = item // TODO: Process file into IO
    } else {
      fmt.Println("Error: Only one file name may be specified")
      os.Exit(1)
    }
  }

  return os.Stdin, options
}


func main() {
  _, options := parseFlag()
  fmt.Printf("%v\n", options)
}