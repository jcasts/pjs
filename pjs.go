package main

import (
  "fmt"
  "flag"
  "io"
  "encoding/json"
  "os"
  "strconv"
  "strings"
  "./paths"
  "./jsonfmt"
  "./iterator"
)


type optionSet struct {
  color bool
  indent uint
  paths []paths.Path
}


func errorAndExit(code int, msg string, args ...interface{}) {
  finalMessage := fmt.Sprintf(msg, args...)
  fmt.Fprintf(os.Stderr, "Error: %s\n", finalMessage)
  os.Exit(code)
}


func parseFlag() (*os.File, optionSet) {
  options := optionSet{}

  colorEnv := os.Getenv("PJS_COLOR")
  options.color = colorEnv != "false"

  indentEnv, _ := strconv.ParseUint(os.Getenv("PJS_INDENT"), 10, 0)
  if indentEnv <= 0 { indentEnv = 2 }
  options.indent = uint(indentEnv)

  name := "pjs"

  flagset := flag.NewFlagSet(name, flag.ExitOnError)
  flagset.BoolVar(&options.color, "c", options.color, "\tOutput in colors")
  flagset.UintVar(&options.indent, "i", options.indent, "\tIndent size")

  flagset.Usage = func() {
    usage := `Pretty print and manipulate JSON data

Usage:
  pjs [options] [filepath] [-- json paths]

Examples:
  pjs path/to/file.json
  pjs path/to/file.json -- **/username=foo
  curl api.twitter.com/1.1/notifications.json | pjs

Options:
`
    fmt.Fprintf(os.Stderr, "\n%s - %s", name, usage)
    flagset.PrintDefaults()
    fmt.Fprintf(os.Stderr, "\n")
    os.Exit(2)
  }

  flagset.Parse(os.Args[1:])

  var err error
  var file *os.File

  for i, item := range os.Args {
    if item == "--" && i != len(os.Args)-1 {
      for _, pathStr := range os.Args[i+1:] {
        pathItem, err := paths.NewPath(pathStr)
        if err != nil { errorAndExit(1, "Invalid path %s", pathStr) }
        options.paths = append(options.paths, pathItem)
      }
      break
    }
  }

  flagArgs := flagset.Args()
  pathLen := len(options.paths)
  argDiff := len(flagArgs) - pathLen
  if (pathLen > 0 && argDiff == 2) || (pathLen == 0 && argDiff == 1) {
    file, err = os.OpenFile(flagArgs[0], os.O_RDONLY, 0)
    if (err != nil) {
      errorAndExit(1, err.Error())
    }
  } else if argDiff == 0 {
    file = os.Stdin
  } else {
    errorAndExit(1, "Only one file name may be specified")
  }

  if file != os.Stdin {
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
      // Data is being piped into stdin even though we have a file
      file.Close()
      errorAndExit(1, "Simultaneous input from pipe and file not supported")
    }
  }

  return file, options
}


func main() {
  input, options := parseFlag()
  defer input.Close()

  indent := strings.Repeat(" ", int(options.indent))
  var formatter *jsonfmt.Formatter

  if options.color {
    formatter = jsonfmt.NewFormatter(
      jsonfmt.NewConsoleColorProcessor(),
      jsonfmt.NewIndentProcessor("", indent))
  } else {
    formatter = jsonfmt.NewFormatter(jsonfmt.NewIndentProcessor("", indent))
  }

  dec := json.NewDecoder(input)
  for {
    var js interface{}
    if err := dec.Decode(&js); err == io.EOF {
      break
    } else if err != nil {
      errorAndExit(2, err.Error())
    }

    if len(options.paths) > 0 {
      // TODO
    } else {
      itValue := iterator.NewDataValue(js, true)
      err := formatter.Process(jsonfmt.NewEncoder(itValue), os.Stdout)
      if err != nil && err != io.EOF {
        os.Stdout.WriteString("\n")
        errorAndExit(3, err.Error())
      }
    }
    os.Stdout.WriteString("\n")
  }
}
