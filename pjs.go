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
  deleteEmptyMatch bool
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

  hideEmptyMatchEnv := os.Getenv("PJS_HIDE_EMPTY")
  options.deleteEmptyMatch = hideEmptyMatchEnv == "true"

  name := "pjs"
  version := "1.0.0"

  flagset := flag.NewFlagSet(name, flag.ExitOnError)
  flagset.BoolVar(&options.color, "c", options.color, "\tOutput in colors")
  flagset.UintVar(&options.indent, "i", options.indent, "\tIndent size")
  flagset.BoolVar(&options.deleteEmptyMatch, "m", options.deleteEmptyMatch, "\tOnly print if a path matches")

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
    env_usage := `
Env Variables:
  PJS_COLOR       true/false   Set default for color output
  PJS_INDENT      NUM          Set default indent number
  PJS_HIDE_EMPTY  true/false   Set default behavior for empty matches
`
    fmt.Fprintf(os.Stderr, "\n%s-%s - %s", name, version, usage)
    flagset.PrintDefaults()
    fmt.Println(env_usage)
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

  stat, _ := os.Stdin.Stat()
  if (stat.Mode() & os.ModeCharDevice) == 0 {
    // Data is being piped into stdin even though we have a file
    if file != os.Stdin {
      file.Close()
      errorAndExit(1, "Simultaneous input from pipe and file not supported")
    }
  } else if file == os.Stdin {
    // Nothing is being passed to stdin
    errorAndExit(1, "No input provided. See pjs -h for usage.")
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
  encoder := jsonfmt.NewEncoder()
  for {
    var js interface{}
    if err := dec.Decode(&js); err == io.EOF {
      break
    } else if err != nil {
      errorAndExit(2, err.Error())
    }

    var itValue iterator.Value
    if len(options.paths) > 0 {
      var matches paths.Matches
      for _, path := range options.paths {
        matches = append(matches, path.FindMatches(js)...)
      }
      if len(matches) == 0 && options.deleteEmptyMatch { continue }

      itValue = matches.IteratorValue()
      if itValue == nil {
        switch js.(type) {
        case map[string]interface{}:
          itValue = iterator.NewDataValue(map[string]interface{}{}, true)
        case []interface{}:
          itValue = iterator.NewDataValue([]interface{}{}, true)
        }
      }
    } else {
      itValue = iterator.NewDataValue(js, true)
    }

    encoder.Queue(itValue)
    err := formatter.Process(encoder, os.Stdout)
    if err != nil && err != io.EOF {
      os.Stdout.WriteString("\n")
      errorAndExit(3, err.Error())
    }
    os.Stdout.WriteString("\n")
  }
}
