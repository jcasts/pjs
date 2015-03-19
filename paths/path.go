package paths

import (
  "errors"
  "fmt"
  "regexp"
  "strconv"
)


type Path interface {
  String() string
  Matches(data interface{}) (bool, [][]interface{})
}


type path struct {
  raw string
  tokens []*token
}

func (p *path) String() string {
  return p.raw
}

func (p *path) Matches(data interface{}) (bool, [][]interface{}) {
  success := false
  return success, nil
}


type tokenMatcher struct {
  regexpMatcher *regexp.Regexp
  rangeMatcher []int64
}

func (tm *tokenMatcher) matches(value interface{}) bool {
  if len(tm.rangeMatcher) < 2 {
    value = fmt.Sprintf("%s", value)
  }

  switch value.(type) {
  case int8, int16, int32, int64, uint8, uint16, uint32:
    num, _ := value.(int64)
    return num >= tm.rangeMatcher[0] && num <= tm.rangeMatcher[1]
  default:
    return tm.regexpMatcher.MatchString(fmt.Sprintf("%s", value))
  }
}


type token struct {
  keyMatcher *tokenMatcher
  valueMatcher *tokenMatcher
}

func (t *token) matches(key interface{}, value interface{}) bool {
  return (t.keyMatcher == nil || t.keyMatcher.matches(key)) &&
    (t.valueMatcher == nil || t.valueMatcher.matches(value))
}


func NewPath(pathStr string) (Path, error) {
  var newPath *path;
  var err error

  if pathStr == "" {
    err = errors.New("Paths can't be empty")
  }

  if err != nil {
    newPath, err = parsePath(pathStr)
  }

  return newPath, err
}


func parsePath(pathStr string) (*path, error) {
  newPath := &path{pathStr, []*token{}}
  var err error
  var matcher *tokenMatcher

  rangeMatcher := regexp.MustCompile("^\\^(\\d+)\\.\\.(\\d+)$")

  item := "^"
  escMode := false
  newTokenKey := false
  newTokenValue := false
  processingKey := true
  currToken := token{}

  for _, ch := range pathStr {
    if newTokenKey || newTokenValue {
      matcher := &tokenMatcher{}
      digits := rangeMatcher.FindAllStringSubmatch(item, -1)[0]
      if len(digits) >= 3 {
        rStart, _ := strconv.ParseInt(digits[1], 10, 64)
        rEnd, _ := strconv.ParseInt(digits[2], 10, 64)
        matcher.rangeMatcher = []int64{rStart, rEnd}
      } else {
        matcher.regexpMatcher, err = regexp.Compile(item + "$")
      }
      item = "^"
    }

    if err != nil {
      newPath = nil
      break
    }

    if newTokenKey {
      if processingKey {
        currToken.keyMatcher = matcher
      } else {
        currToken.valueMatcher = matcher
      }
      newPath.tokens = append(newPath.tokens, &currToken)
      currToken = token{}
      newTokenKey = false
      newTokenValue = false
      processingKey = true
    }

    if newTokenValue {
      currToken.keyMatcher = matcher
      newTokenValue = false
      processingKey = false
    }

    if escMode {
      item += regexp.QuoteMeta(runeToAscii(ch))
      escMode = false
      continue
    }

    switch ch {
      case '*', '?':
        item += "." + runeToAscii(ch)
      case '\\':
        escMode = true
      case '(', ')', '|':
        item += runeToAscii(ch)
      case '/':
        newTokenKey = true
      case '=':
        newTokenValue = true
      default:
        item += regexp.QuoteMeta(runeToAscii(ch))
    }
  }

  return newPath, err
}

func runeToAscii(r rune) string {
    if r < 128 {
        return string(r)
    } else {
        return "\\u" + strconv.FormatInt(int64(r), 16)
    }
}

