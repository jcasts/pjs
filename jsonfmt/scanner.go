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
  IntegerLiteralToken       // 1234
  FloatLiteralToken         // 0.234
  ScientificLiteralToken    // 0.234E-2
  BooleanLiteralToken       // true/false
  NullLiteralToken          // null
  MapKeyToken               // "something"
  MapStartToken             // {
  MapEndToken               // }
  MapColonToken             // :
  EmptyMapToken             // {}
  ArrayStartToken           // [
  ArrayEndToken             // ]
  EmptyArrayToken           // []
  ValueSeparatorToken       // ,
  StartNewJsonToken         // Used to delimit Json structures in a stream of many
)

type scannerDataType int
const (
  scannerMap scannerDataType = iota
  scannerArray
  scannerString
)


const bufSize = 1024
type stepFunc func(*Scanner, rune)error

type Scanner struct {
  buffer []byte
  bufferPos int
  bufferLen int
  reader io.Reader
  token *Token
  err error
  dataTypes []scannerDataType
  value string
  inStringEsc bool
  step stepFunc
}

func NewScanner(r io.Reader) *Scanner {
  return &Scanner{make([]byte, bufSize), 0, 0, r, nil, nil, []scannerDataType{}, "", false, parseAny}
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

  for s.token == nil {
    if s.bufferPos >= s.bufferLen - 1 {
      var err error
      s.bufferLen, err = s.reader.Read(s.buffer)
      if err != nil {
        if err == io.EOF {
          s.step(s, ' ') // Force trigger unfinished last values
        } else {
          s.err = err
        }
        return s.token != nil
      }
    }

    for i := s.bufferPos; i < s.bufferLen; i++ {
      s.bufferPos++
      b := s.buffer[i]
      err := s.step(s, rune(int(b)))
      if err != nil {
        s.err = err
        return false
      }
      if s.token != nil { break }
    }
  }

  return true
}

func (s *Scanner) inObjectType(t scannerDataType) bool {
  return s.dataTypes[len(s.dataTypes)-1] == t
}

func (s *Scanner) addObjectType(t scannerDataType) {
  s.dataTypes = append(s.dataTypes, t)
}

func (s *Scanner) popObjectType() {
  s.dataTypes = s.dataTypes[0:len(s.dataTypes)-1]
}

func (s *Scanner) finishWithTokenType(tokenType TokenType) {
  // TODO: Handle potential dangling EndOfValue chars by setting bufferPos back by one when necessary
  // Handle trimming of s.dataTypes
  s.token = &Token{s.value, tokenType, len(s.dataTypes),
        s.inObjectType(scannerArray)}
  s.step = parseNextInObject
}

func isBlank(char rune) bool {
  return char == '\n' || char == '\t' || char == '\r' || char == ' '
}

func isTermination(char rune) bool {
  return char == ',' || char == '}' || char == ']' // TODO
}

func isEndOfValue(char rune) bool {
  return isBlank(char) || isTermination(char) // TODO: Make sure we're in a map or array
}

func parseError() error {
  return errors.New("Unexpected character while reading JSON")
}

func parseAny(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  switch char {
    case '{':
      s.step = parseMapStart
    case '[':
      s.step = parseArrayStart
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

func parseMapStart(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == '}' {
    s.value += string(char)
    s.finishWithTokenType(EmptyMapToken)
  } else {
    s.bufferPos--
    s.addObjectType(scannerMap)
    s.finishWithTokenType(MapStartToken)
    s.step = parseMapKey
  }
  return nil
}

func parseArrayStart(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == ']' {
    s.value += string(char)
    s.finishWithTokenType(EmptyArrayToken)
  } else {
    s.bufferPos--
    s.addObjectType(scannerArray)
    s.finishWithTokenType(ArrayStartToken)
  }
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
    s.bufferPos--
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
    s.bufferPos--
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
    s.bufferPos--
    s.finishWithTokenType(ScientificLiteralToken)
    return nil
  }
  return parseScientific1(s, char)
}

func parseMapKey(s *Scanner, char rune) error {
  if char == '"' && s.value == "" {
    s.value += string(char)
    return nil
  } else if !s.inObjectType(scannerString) && s.value != "" {
    s.bufferPos--
    s.finishWithTokenType(MapKeyToken)
    s.step = parseMapKeyAssign
    return nil
  } else if s.value == "" {
    return parseError()
  }
  return parseString(s, char)
}

func parseMapKeyAssign(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == ':' {
    s.value += string(char)
    s.finishWithTokenType(MapColonToken)
    return nil
  }
  return parseError()
}

func parseNextInObject(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == ',' && len(s.dataTypes) > 0 {
    // Next Value
    s.value += string(char)
    s.finishWithTokenType(ValueSeparatorToken)
    if s.inObjectType(scannerMap) {
      s.step = parseMapKey
    } else {
      s.step = parseAny
    }
    return nil
  } else if char == ']' && s.inObjectType(scannerArray) {
    // End of array
    s.value += string(char)
    s.popObjectType()
    s.finishWithTokenType(ArrayEndToken)
    s.step = parseNextInObject
  } else if char == '}' && s.inObjectType(scannerMap) {
    // End of map
    s.value += string(char)
    s.popObjectType()
    s.finishWithTokenType(MapEndToken)
    s.step = parseNextInObject
  } else if len(s.dataTypes) == 0 && !isTermination(char) {
    // End of JSON
    s.bufferPos--
    s.finishWithTokenType(StartNewJsonToken)
    s.step = parseAny
  }
  return parseError()
}

func parseString(s *Scanner, char rune) error {
  inString := s.inObjectType(scannerString)
  if !inString && isEndOfValue(char) {
    s.bufferPos--
    s.finishWithTokenType(StringLiteralToken)
    return nil
  }

  s.value += string(char)

  if char == '"' && !s.inStringEsc && inString {
    s.popObjectType()
  } else {
    s.addObjectType(scannerString)
  }

  s.inStringEsc = char == '\\' && !s.inStringEsc
  // TODO: case '\u'?

  return nil
}

func parseFalse(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "false" {
    s.bufferPos--
    s.finishWithTokenType(BooleanLiteralToken)
    return nil
  }
  if char == 'a' && s.value == "f" {
    s.value += string(char)
  } else if char == 'l' && s.value == "fa" {
    s.value += string(char)
  } else if char == 's' && s.value == "fal" {
    s.value += string(char)
  } else if char == 'e' && s.value == "fals" {
    s.value += string(char)
  } else {
    return parseError()
  }
  return nil
}

func parseTrue(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "true" {
    s.bufferPos--
    s.finishWithTokenType(BooleanLiteralToken)
    return nil
  }
  if char == 'r' && s.value == "t" {
    s.value += string(char)
  } else if char == 'u' && s.value == "tr" {
    s.value += string(char)
  } else if char == 'e' && s.value == "tru" {
    s.value += string(char)
  } else {
    return parseError()
  }
  return nil
}

func parseNull(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "null" {
    s.bufferPos--
    s.finishWithTokenType(NullLiteralToken)
    return nil
  }
  if char == 'u' && s.value == "n" {
    s.value += string(char)
  } else if char == 'l' && s.value == "nu" {
    s.value += string(char)
  } else if char == 'l' && s.value == "nul" {
    s.value += string(char)
  } else {
    return parseError()
  }
  return nil
}
