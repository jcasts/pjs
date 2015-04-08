package jsonfmt

import (
  "errors"
  "io"
)

type Token struct {
  Value string
  Type TokenType
  Depth int
  InArray bool
}

type TokenType int
const (
  StringLiteralToken TokenType = iota // "Foo"
  IntegerLiteralToken   // 1234
  FloatLiteralToken     // 0.234
  ScientificLiteralToken     // 0.234E-2
  BooleanLiteralToken   // true/false
  NullLiteralToken      // null
  MapKeyToken           // "something"
  MapStartToken         // {
  MapEndToken           // }
  MapColonToken         // :
  EmptyMapToken         // {}
  ArrayStartToken       // [
  ArrayEndToken         // ]
  EmptyArrayToken       // []
  ValueSeparatorToken   // ,
)

type scannerDataType int
const (
  scannerMap scannerDataType = iota
  scannerArray
)



type stepFunc func(*Scanner, rune)error

type Scanner struct {
  reader io.Reader
  token *Token
  err error
  dataTypes []scannerDataType
  value string
  inStringEsc bool
  inNumber bool
  step stepFunc
}

func NewScanner(r io.Reader) *Scanner {
  return &Scanner{r, nil, nil, []scannerDataType{}, "", false, false, parseAny}
}

func (s *Scanner) Token() *Token {
  return s.token
}

func (s *Scanner) Error() error {
  return s.err
}

func (s *Scanner) Next() bool {
  s.token = nil
  s.value = ""
  buffer := []byte{}

  for s.token == nil {
    byteCount, err := s.reader.Read(buffer)
    if byteCount == 0 {
      if err != io.EOF { s.err = err }
      return false
    }

    for _, b := range buffer {
      s.step(s, rune(int(b)))
      if s.token != nil { break }
      // TODO: Handle left over buffer
    }

    if s.token == nil { s.step(s, ' ') } // Force trigger unfinished values
  }

  return true
}

func (s *Scanner) finishWithTokenType(tokenType TokenType) {
  s.token = &Token{s.value, tokenType, len(s.dataTypes),
        s.dataTypes[len(s.dataTypes)-1] == scannerArray}
  s.step = parseAny
  // TODO: Handle potential dangling chars
}

func isBlank(char rune) bool {
  return char == '\n' || char == '\t' || char == '\r' || char == ' '
}

func isEndOfValue(char rune) bool {
  return isBlank(char) || char == ','
}

func parseError() error {
  return errors.New("Unexpected character while reading JSON")
}

func parseAny(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  switch char {
    case '{':
    case '}':
    case '[':
    case ']':
    case ':':
      s.value += string(char)
      s.finishWithTokenType(MapColonToken)
      return nil
    case ',':
      s.value += string(char)
      s.finishWithTokenType(ValueSeparatorToken)
      return nil
    case 'f':
      s.step = parseFalse
    case 't':
      s.step = parseTrue
    case 'n':
      s.step = parseNull
    case '"':
      s.step = parseString
    case '-':
      s.step = parseNegNumber
    case '0':
      s.step = parseFloat0
    case '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.step = parseNumber
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseNegNumber(s *Scanner, char rune) error {
  switch char {
    case '0':
      s.step = parseFloat0
    case '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.step = parseNumber
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseNumber(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.finishWithTokenType(IntegerLiteralToken)
    return nil
  }

  switch char {
    case '.':
      s.step = parseFloat1
    case 'e', 'E':
      s.step = parseScientific0
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseFloat0(s *Scanner, char rune) error {
  if char == '.' {
    s.step = parseFloat1
  } else {
    return parseError()
  }
  s.value += string(char)
  return nil
}

func parseFloat1(s *Scanner, char rune) error {
  switch char {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.step = parseFloat2
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseFloat2(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.finishWithTokenType(FloatLiteralToken)
    return nil
  } else if char == 'e' || char == 'E' {
    s.step = parseScientific0
    return nil
  }
  return parseFloat1(s, char)
}

func parseScientific0(s *Scanner, char rune) error {
  switch char {
    case '-', '+':
      s.step = parseScientific1
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.step = parseScientific2
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseScientific1(s *Scanner, char rune) error {
  switch char {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.step = parseScientific2
    default:
      return parseError()
  }
  s.value += string(char)
  return nil
}

func parseScientific2(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.finishWithTokenType(ScientificLiteralToken)
    return nil
  }
  return parseScientific1(s, char)
}

func parseString(s *Scanner, char rune) error {
  if char == '"' && !s.inStringEsc {
    s.finishWithTokenType(StringLiteralToken)
    return nil
  }

  s.inStringEsc = char == '\\' && !s.inStringEsc
      //case '\u':

  s.value += string(char)
  return nil
}

func parseFalse(s *Scanner, char rune) error {
  if char == 'a' && s.value == "f" {
    s.value += string(char)
  } else if char == 'l' && s.value == "fa" {
    s.value += string(char)
  } else if char == 's' && s.value == "fal" {
    s.value += string(char)
  } else if char == 'e' && s.value == "fals" {
    s.value += string(char)
    s.finishWithTokenType(BooleanLiteralToken)
  } else {
    return parseError()
  }
  return nil
}

func parseTrue(s *Scanner, char rune) error {
  if char == 'r' && s.value == "t" {
    s.value += string(char)
  } else if char == 'u' && s.value == "tr" {
    s.value += string(char)
  } else if char == 'e' && s.value == "tru" {
    s.value += string(char)
    s.finishWithTokenType(BooleanLiteralToken)
  } else {
    return parseError()
  }
  return nil
}

func parseNull(s *Scanner, char rune) error {
  if char == 'u' && s.value == "n" {
    s.value += string(char)
  } else if char == 'l' && s.value == "nu" {
    s.value += string(char)
  } else if char == 'l' && s.value == "nul" {
    s.value += string(char)
    s.finishWithTokenType(NullLiteralToken)
  } else {
    return parseError()
  }
  return nil
}
