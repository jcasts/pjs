package jsonfmt

import (
  "fmt"
  "errors"
  "io"
  "unicode/utf8"
  "strconv"
)

type Token struct {
  Value string
  Type TokenType
  Depth int
  InMap bool
}

type TokenType int
const (
  AnyToken TokenType = iota
  StringLiteralToken        // "Foo"
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

func tokenTypeName(t TokenType) string {
  switch t {
    case StringLiteralToken:
      return "string"
    case IntegerLiteralToken:
      return "integer"
    case FloatLiteralToken:
      return "float"
    case ScientificLiteralToken:
      return "scientific"
    case BooleanLiteralToken:
      return "boolean"
    case NullLiteralToken:
      return "null"
    case MapKeyToken:
      return "map key"
    case MapStartToken, MapEndToken, MapColonToken, EmptyMapToken:
      return "map"
    case ArrayStartToken, ArrayEndToken, EmptyArrayToken:
      return "array"
    case ValueSeparatorToken:
      return "data structure"
    default:
      return "JSON" 
  }
}


const bufSize = 128 // Benchmarks show this is a pretty optimal buffer size on large JSON structs
type stepFunc func(*Scanner, rune)error

type Scanner struct {
  buffer []byte
  bufferPos int
  bufferLen int
  runeLen int
  streamPos int
  reader io.Reader
  token *Token
  err error
  dataTypes []scannerDataType
  value string
  step stepFunc
  ignoreConsoleChars bool
}

func NewScanner(r io.Reader) *Scanner {
  s := &Scanner{make([]byte, bufSize), 0, 0, 0, 0, r, nil, nil, []scannerDataType{}, "", nil, true}
  s.setStep(parseAny)
  return s
}

func (s *Scanner) Token() *Token {
  return s.token
}

func (s *Scanner) Error() error {
  return s.err
}

func (s *Scanner) Next() bool {
  s.err = nil
  s.token = nil
  s.value = ""
  runeBuffer := []byte{}
  runeError := false

  for s.token == nil {
    if s.bufferPos > s.bufferLen - 1 || runeError {
      var err error
      s.bufferLen, err = s.reader.Read(s.buffer)
      s.bufferPos = 0
      if err != nil && s.bufferLen == 0 {
        s.err = err
        if err == io.EOF {
          s.step(s, ' ') // Force trigger unfinished last values
        }
        return s.token != nil
      }
    }

    lastRuneLen := 0
    posDiff := len(runeBuffer)
    index := s.bufferPos
    if index < 0 { index = 0 }
    runeBuffer = append(runeBuffer, s.buffer[index:s.bufferLen]...)
    s.bufferPos -= posDiff

    for i := s.bufferPos; i < s.bufferLen; i+=lastRuneLen {
      var char rune
      runeError = false
      char, lastRuneLen = utf8.DecodeRune(runeBuffer)

      if char == utf8.RuneError {
        runeError = true
        if len(runeBuffer) >= 8 {
          s.err = errors.New(fmt.Sprintf("Unparsable UTF-8: %v", runeBuffer))
          return false
        } else {
          break
        }
      }
      s.runeLen = lastRuneLen
      s.bufferPos += lastRuneLen
      runeBuffer = runeBuffer[lastRuneLen:]

      err := s.step(s, char)

      if err != nil {
        s.err = err
        return false
      }
      s.streamPos += 1
      if s.token != nil { break }
    }
  }

  return true
}

func (s *Scanner) emptyObjectType() bool {
  return len(s.dataTypes) == 0
}

func (s *Scanner) inObjectType(t scannerDataType) bool {
  if s.emptyObjectType() { return false }
  return s.dataTypes[len(s.dataTypes)-1] == t
}

func (s *Scanner) addObjectType(t scannerDataType) {
  s.dataTypes = append(s.dataTypes, t)
}

func (s *Scanner) popObjectType() {
  if len(s.dataTypes) == 0 { return }
  s.dataTypes = s.dataTypes[0:len(s.dataTypes)-1]
}

func (s *Scanner) setStep(fn stepFunc) {
  s.step = func(s *Scanner, char rune) error {
    if checkParseConsoleEsc(s, char, fn) { return nil }
    return fn(s, char)
  }
}

func (s *Scanner) stepBack() {
  if s.streamPos > 0 {
    s.bufferPos -= s.runeLen
    s.streamPos--
  }
}

func (s *Scanner) finishWithTokenType(tokenType TokenType) {
  s.token = &Token{s.value, tokenType, len(s.dataTypes),
        s.inObjectType(scannerMap)}
  s.setStep(parseNextInObject)
}

func checkParseConsoleEsc(s *Scanner, char rune, next stepFunc) bool {
  if char != '\033' || !s.ignoreConsoleChars { return false }

  fn1 := func(s *Scanner, char rune) error {
    switch char {
      case 'm':
        s.setStep(next)
      case ';', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      default:
        return parseError(s.streamPos, char, AnyToken)
    }
    return nil
  }

  fn0 := func(s *Scanner, char rune) error {
    if char != '[' { return parseError(s.streamPos, char, AnyToken) }
    s.step = fn1
    return nil
  }

  s.step = fn0

  return true
}

func isBlank(char rune) bool {
  return char == '\n' || char == '\t' || char == '\r' || char == ' '
}

func isTermination(char rune) bool {
  return char == ':' || char == ',' || char == '}' || char == ']'
}

func isEndOfValue(char rune) bool {
  return isBlank(char) || isTermination(char)
}

func parseError(pos int, char rune, t TokenType) error {
  msg := "Unexpected character '"+string(char)+"' in "+tokenTypeName(t)+" at position "+strconv.Itoa(pos)
  return errors.New(msg)
}

func parseAny(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  switch char {
    case '{':
      s.setStep(parseMapStart)
    case '[':
      s.setStep(parseArrayStart)
    case 'f':
      s.setStep(parseFalse)
    case 't':
      s.setStep(parseTrue)
    case 'n':
      s.setStep(parseNull)
    case '"':
      s.setStep(parseString)
    case '-':
      s.setStep(parseNegNumber)
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.setStep(parseNumber)
    default:
      return parseError(s.streamPos, char, AnyToken)
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
    s.stepBack()
    s.finishWithTokenType(MapStartToken)
    s.addObjectType(scannerMap)
    s.setStep(parseMapKey)
  }
  return nil
}

func parseArrayStart(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == ']' {
    s.value += string(char)
    s.finishWithTokenType(EmptyArrayToken)
  } else {
    s.stepBack()
    s.finishWithTokenType(ArrayStartToken)
    s.addObjectType(scannerArray)
    s.setStep(parseAny)
  }
  return nil
}

func parseNegNumber(s *Scanner, char rune) error {
  switch char {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.setStep(parseNumber)
    default:
      return parseError(s.streamPos, char, IntegerLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseNumber(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.stepBack()
    s.finishWithTokenType(IntegerLiteralToken)
    return nil
  }

  switch char {
    case '.':
      s.setStep(parseFloat1)
    case 'e', 'E':
      s.setStep(parseScientific0)
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      if s.value == "0" || s.value == "-0" { return parseFloat0(s, char) }
    default:
      return parseError(s.streamPos, char, IntegerLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseFloat0(s *Scanner, char rune) error {
  if char == '.' {
    s.setStep(parseFloat1)
  } else {
    return parseError(s.streamPos, char, FloatLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseFloat1(s *Scanner, char rune) error {
  switch char {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.setStep(parseFloat2)
    default:
      return parseError(s.streamPos, char, FloatLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseFloat2(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.stepBack()
    s.finishWithTokenType(FloatLiteralToken)
    return nil
  } else if char == 'e' || char == 'E' {
    s.value += string(char)
    s.setStep(parseScientific0)
    return nil
  }
  return parseFloat1(s, char)
}

func parseScientific0(s *Scanner, char rune) error {
  switch char {
    case '-', '+':
      s.setStep(parseScientific1)
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.setStep(parseScientific2)
    default:
      return parseError(s.streamPos, char, ScientificLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseScientific1(s *Scanner, char rune) error {
  switch char {
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      s.setStep(parseScientific2)
    default:
      return parseError(s.streamPos, char, ScientificLiteralToken)
  }
  s.value += string(char)
  return nil
}

func parseScientific2(s *Scanner, char rune) error {
  if isEndOfValue(char) {
    s.stepBack()
    s.finishWithTokenType(ScientificLiteralToken)
    return nil
  }
  return parseScientific1(s, char)
}

func parseMapKey(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char != '"' {
    return parseError(s.streamPos, char, MapKeyToken)
  }
  s.value += string(char)
  s.setStep(stringLikeStep(MapKeyToken))
  return nil
}

func parseNextInObject(s *Scanner, char rune) error {
  if isBlank(char) { return nil }
  if char == ',' && !s.emptyObjectType() {
    // Next Value
    s.value += string(char)
    s.finishWithTokenType(ValueSeparatorToken)
    if s.inObjectType(scannerMap) {
      s.setStep(parseMapKey)
    } else {
      s.setStep(parseAny)
    }
    return nil
  } else if char == ':' && s.inObjectType(scannerMap) {
    // End of map key
    s.value += string(char)
    s.finishWithTokenType(MapColonToken)
    s.setStep(parseAny)
    return nil
  } else if char == ']' && s.inObjectType(scannerArray) {
    // End of array
    s.value += string(char)
    s.popObjectType()
    s.finishWithTokenType(ArrayEndToken)
    s.setStep(parseNextInObject)
    return nil
  } else if char == '}' && s.inObjectType(scannerMap) {
    // End of map
    s.value += string(char)
    s.popObjectType()
    s.finishWithTokenType(MapEndToken)
    s.setStep(parseNextInObject)
    return nil
  } else if s.emptyObjectType() && !isTermination(char) {
    // End of JSON
    s.stepBack()
    s.finishWithTokenType(StartNewJsonToken)
    s.setStep(parseAny)
    return nil
  }
  return parseError(s.streamPos, char, ValueSeparatorToken)
}

func parseString(s *Scanner, char rune) error {
  return stringLikeStep(StringLiteralToken)(s, char)
}

func stringLikeStep(tokenType TokenType) stepFunc {
  return func(s *Scanner, char rune) error {
    return parseStringLike(s, char, tokenType)
  }
}

func stringLikeEscStep(tokenType TokenType) stepFunc {
  return func(s *Scanner, char rune) error {
    return parseStringLikeEsc(s, char, tokenType)
  }
}

func parseStringLike(s *Scanner, char rune, tokenType TokenType) error {
  if !s.inObjectType(scannerString) && s.value == "\"" {
    s.addObjectType(scannerString)
  }

  if s.inObjectType(scannerString) {
    s.value += string(char)

    if char == '"' && len(s.value) > 1 {
      s.popObjectType()
    } else if char == '\\' {
      s.setStep(stringLikeEscStep(tokenType))
    } else {
      s.setStep(stringLikeStep(tokenType))
    }
    // TODO: case '\u'?
  } else {
    if isEndOfValue(char) {
      s.stepBack()
      s.finishWithTokenType(tokenType)
    } else {
      return parseError(s.streamPos, char, tokenType)
    }
  }

  return nil
}

func parseStringLikeEsc(s *Scanner, char rune, tokenType TokenType) error {
  s.value += string(char)
  s.setStep(stringLikeStep(tokenType))
  return nil
}

func parseFalse(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "false" {
    s.stepBack()
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
    return parseError(s.streamPos, char, BooleanLiteralToken)
  }
  return nil
}

func parseTrue(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "true" {
    s.stepBack()
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
    return parseError(s.streamPos, char, BooleanLiteralToken)
  }
  return nil
}

func parseNull(s *Scanner, char rune) error {
  if isEndOfValue(char) && s.value == "null" {
    s.stepBack()
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
    return parseError(s.streamPos, char, NullLiteralToken)
  }
  return nil
}
