package jsonfmt

import (
  "io"
  "strings"
  "testing"
)

const testJSON = "{\"foo\":[1, [-23,false, \"hi 🍷🍷🍷\"], 0.23, [], {}, 2.3e-23],\"bar\":null}"

type stringWriter string

func (w *stringWriter) Write(p []byte) (int, error) {
  *w = stringWriter(string(*w) + string(p))
  return len(p), nil
}


func TestColorFormatter(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewConsoleColorProcessor())
  err := f.Process(strings.NewReader(testJSON), &output)

  expectedOutput :=
  "{\"foo\":[\033[0;33m1\033[0m,[\033[0;33m-23\033[0m,\033[1;35m"+
  "false\033[0m,\033[0;36m\"hi 🍷🍷🍷\"\033[0m],\033[0;33m0.23\033[0m,[]"+
  ",{},\033[0;33m2.3e-23\033[0m],\"bar\":\033[1;31mnull\033[0m}"

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expectedOutput, string(output))
}

func TestIndentFormatter(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewIndentProcessor("", "  "))
  err := f.Process(strings.NewReader(testJSON), &output)

  expectedOutput :=
`{
  "foo": [
    1,
    [
      -23,
      false,
      "hi 🍷🍷🍷"
    ],
    0.23,
    [],
    {},
    2.3e-23
  ],
  "bar": null
}`
  
  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expectedOutput, string(output))
}

func TestIndentFormatterPrefix(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewIndentProcessor(">>", " "))
  err := f.Process(strings.NewReader(testJSON), &output)

  expectedOutput :=
`>>{
>> "foo": [
>>  1,
>>  [
>>   -23,
>>   false,
>>   "hi 🍷🍷🍷"
>>  ],
>>  0.23,
>>  [],
>>  {},
>>  2.3e-23
>> ],
>> "bar": null
>>}`
  
  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expectedOutput, string(output))
}

func TestCombinedFormatter(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewConsoleColorProcessor(), NewIndentProcessor(">>", " "))
  err := f.Process(strings.NewReader(testJSON), &output)

  expectedOutput :=
">>{\n"+
">> \"foo\": [\n"+
">>  \033[0;33m1\033[0m,\n"+
">>  [\n"+
">>   \033[0;33m-23\033[0m,\n"+
">>   \033[1;35mfalse\033[0m,\n"+
">>   \033[0;36m\"hi 🍷🍷🍷\"\033[0m\n"+
">>  ],\n"+
">>  \033[0;33m0.23\033[0m,\n"+
">>  [],\n"+
">>  {},\n"+
">>  \033[0;33m2.3e-23\033[0m\n"+
">> ],\n"+
">> \"bar\": \033[1;31mnull\033[0m\n"+
">>}"

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expectedOutput, string(output))
}

func TestStreamFormatter(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewConsoleColorProcessor(), NewIndentProcessor(">>", " "))
  err := f.Process(strings.NewReader(testJSON+" "+testJSON), &output)

  expectedOutput :=
">>{\n"+
">> \"foo\": [\n"+
">>  \033[0;33m1\033[0m,\n"+
">>  [\n"+
">>   \033[0;33m-23\033[0m,\n"+
">>   \033[1;35mfalse\033[0m,\n"+
">>   \033[0;36m\"hi 🍷🍷🍷\"\033[0m\n"+
">>  ],\n"+
">>  \033[0;33m0.23\033[0m,\n"+
">>  [],\n"+
">>  {},\n"+
">>  \033[0;33m2.3e-23\033[0m\n"+
">> ],\n"+
">> \"bar\": \033[1;31mnull\033[0m\n"+
">>}\n\n"+
">>{\n"+
">> \"foo\": [\n"+
">>  \033[0;33m1\033[0m,\n"+
">>  [\n"+
">>   \033[0;33m-23\033[0m,\n"+
">>   \033[1;35mfalse\033[0m,\n"+
">>   \033[0;36m\"hi 🍷🍷🍷\"\033[0m\n"+
">>  ],\n"+
">>  \033[0;33m0.23\033[0m,\n"+
">>  [],\n"+
">>  {},\n"+
">>  \033[0;33m2.3e-23\033[0m\n"+
">> ],\n"+
">> \"bar\": \033[1;31mnull\033[0m\n"+
">>}"

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expectedOutput, string(output))
}

func TestBadJsonFormatter(t *testing.T) {
  output := stringWriter("")
  f := NewFormatter(NewConsoleColorProcessor(), NewIndentProcessor(">>", " "))
  err := f.Process(strings.NewReader("{123: 123"), &output)
  expectedOutput := ">>{\n>> "
  testAssertEqual(t, "Unexpected character '1' in map key", err.Error())
  testAssertEqual(t, expectedOutput, string(output))
}
