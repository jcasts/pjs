package jsonfmt

import (
  //"fmt"
  "io"
  "strings"
  "testing"
)

func testAssertToken(t *testing.T, value string, tt TokenType, depth int, inArray bool, tn *Token) {
  testAssertEqual(t, value, tn.Value)
  testAssertEqual(t, tt, tn.Type)
  testAssertEqual(t, depth, tn.Depth)
  testAssertEqual(t, inArray, tn.InArray)
}

func TestInteger(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader(" 123"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "123", IntegerLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("-123"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "-123", IntegerLiteralToken, 0, false, scan.Token())
}

func TestBadInteger(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("1-23"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '-' in integer", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("1lskd23"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character 'l' in integer", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestFloat(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("1.23"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "1.23", FloatLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("0.123"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "0.123", FloatLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("-0.123"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "-0.123", FloatLiteralToken, 0, false, scan.Token())
}

func TestBadFloat(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("023"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '2' in float", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("-023"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '2' in float", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("0.23.3"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '.' in float", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestScientific(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("1.23e3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "1.23e3", ScientificLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("1.12E3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "1.12E3", ScientificLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("-1.12E3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "-1.12E3", ScientificLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("1.12e+3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "1.12e+3", ScientificLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("1.12e-3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "1.12e-3", ScientificLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("112e+3"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "112e+3", ScientificLiteralToken, 0, false, scan.Token())
}

func TestBadScientific(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("1.23e--3"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '-' in scientific", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("1.23e-3."))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '.' in scientific", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestBoolean(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("true"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "true", BooleanLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("false"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "false", BooleanLiteralToken, 0, false, scan.Token())
}

func TestBadBoolean(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("truue"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character 'u' in boolean", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("falsee"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character 'e' in boolean", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestNull(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("null"))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "null", NullLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestBadNull(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("nulll"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character 'l' in null", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("nul"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestString(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("\"null\""))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"null\"", StringLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("\"quote \\\"thing\\\"\""))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "\"quote \\\"thing\\\"\"", StringLiteralToken, 0, false, scan.Token())

  scan = NewScanner(strings.NewReader("\"123\""))
  testAssertTrue(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertToken(t, "\"123\"", StringLiteralToken, 0, false, scan.Token())
}

func TestBadString(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("\"thing"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("\"thing\"more"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character 'm' in string", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestBadValue(t *testing.T) {
  var scan *Scanner

  badValues := []string{"]", "}", ":", "apple", "+", "=", "\\"}

  for _, val := range badValues {
    scan = NewScanner(strings.NewReader(val))
    testAssertFalse(t, scan.Next())
    testAssertEqual(t, "Unexpected character '"+string(val[0])+"' in JSON", scan.Error().Error())
    testAssertTrue(t, scan.Token() == nil)
  }
}

func TestValueStream(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("123\n\"foo\"\nfalse"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "123", IntegerLiteralToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "", StartNewJsonToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", StringLiteralToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "", StartNewJsonToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "false", BooleanLiteralToken, 0, false, scan.Token())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestBadValueStream(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("\"foo\",false"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", StringLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character ',' in data structure", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("\"foo\":false"))
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character ':' in string", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("\"foo\"]false"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", StringLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character ']' in data structure", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("\"foo\"}false"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", StringLiteralToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '}' in data structure", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestEmptyArray(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader(" [\n  ] "))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[]", EmptyArrayToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("[]"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[]", EmptyArrayToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestArray(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("[\n  \"foo\",\n  123\n]"))

  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", StringLiteralToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "123", IntegerLiteralToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "]", ArrayEndToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestNestedArray(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("[null,[-23.2,false,[\"hi\"]],\n  123\n, []]"))

  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "null", NullLiteralToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "-23.2", FloatLiteralToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "false", BooleanLiteralToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"hi\"", StringLiteralToken, 3, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "]", ArrayEndToken, 2, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "]", ArrayEndToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "123", IntegerLiteralToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[]", EmptyArrayToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "]", ArrayEndToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestBadArray(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("[null,,\n  123\n, []]"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "null", NullLiteralToken, 1, true, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, true, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character ',' in JSON", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)

  scan = NewScanner(strings.NewReader("[null  123\n, []]"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "[", ArrayStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "null", NullLiteralToken, 1, true, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, "Unexpected character '1' in data structure", scan.Error().Error())
  testAssertTrue(t, scan.Token() == nil)
}

func TestEmptyMap(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader(" {\n  } "))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{}", EmptyMapToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())

  scan = NewScanner(strings.NewReader("{}"))
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{}", EmptyMapToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestMap(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("{\"foo\": 12, \"bar\":false,\"21\":-2.4e7}"))

  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{", MapStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", MapKeyToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "12", IntegerLiteralToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"bar\"", MapKeyToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "false", BooleanLiteralToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"21\"", MapKeyToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "-2.4e7", ScientificLiteralToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "}", MapEndToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestNestedMap(t *testing.T) {
  var scan *Scanner

  scan = NewScanner(strings.NewReader("{\"foo\": {\"bar\":false, \"baz\":true } ,\"21\":{\"sci\":{\"num\":-2.4e7}}}"))

  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{", MapStartToken, 0, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"foo\"", MapKeyToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{", MapStartToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"bar\"", MapKeyToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "false", BooleanLiteralToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"baz\"", MapKeyToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "true", BooleanLiteralToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "}", MapEndToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ",", ValueSeparatorToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"21\"", MapKeyToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{", MapStartToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"sci\"", MapKeyToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "{", MapStartToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "\"num\"", MapKeyToken, 3, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, ":", MapColonToken, 3, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "-2.4e7", ScientificLiteralToken, 3, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "}", MapEndToken, 2, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "}", MapEndToken, 1, false, scan.Token())
  testAssertTrue(t, scan.Next())
  testAssertToken(t, "}", MapEndToken, 0, false, scan.Token())
  testAssertFalse(t, scan.Next())
  testAssertEqual(t, io.EOF, scan.Error())
}

func TestBadMap(t *testing.T) {

}

func TestMixedStream(t *testing.T) {

}
