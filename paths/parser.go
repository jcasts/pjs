package paths

import (
  "errors"
  //"fmt"
  "regexp"
  "strconv"
)

func parsePath(pathStr string) (*path, error) {
  newPath := &path{pathStr, []*pathToken{}}
  var err error
  var matcher *tokenMatcher

  if pathStr == "" {
    err = errors.New("Paths can't be empty")
    return nil, err
  }

  item := ""
  escMode := false
  newTokenKey := false
  newTokenValue := false
  processingKey := true
  currToken := &pathToken{}

  if pathStr[len(pathStr)-1:] != "/" {
    pathStr = pathStr + "/"
  }

  for _, ch := range pathStr {
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
        if !processingKey || item != "" {
          newTokenKey = true
        }
      case '=':
        if !processingKey {
          err = errors.New("Multiple '=' invalid. Use \\= to match character.")
          break
        }
        newTokenValue = true
      default:
        item += regexp.QuoteMeta(runeToAscii(ch))
    }

    if newTokenKey || newTokenValue {
      matcher, err = matcherForTokenString(item, processingKey)
      item = ""
      if processingKey {
        currToken.keyMatcher = matcher
      } else {
        currToken.valueMatcher = matcher
      }
    }

    if err != nil {
      break
    }

    if newTokenKey {
      newPath.tokens = append(newPath.tokens, currToken)
      currToken = &pathToken{}
      processingKey = true
    }

    if newTokenValue {
      processingKey = false
    }

    newTokenKey = false
    newTokenValue = false
  }

  if err != nil {
    newPath = nil
  }

  return newPath, err
}

func matcherForTokenString(tokenStr string, isKey bool) (*tokenMatcher, error) {
  var err error
  rangeRegexp := regexp.MustCompile("^(\\d+)\\\\\\.\\\\\\.(\\d+)$")
  digitMatches := rangeRegexp.FindAllStringSubmatch(tokenStr, -1)
  digits := []string{}
  if len(digitMatches) > 0 { digits = digitMatches[0] }

  matcher := &tokenMatcher{}

  if len(digits) >= 3 {
    rStart, _ := strconv.ParseInt(digits[1], 10, 64)
    rEnd, _ := strconv.ParseInt(digits[2], 10, 64)
    matcher.rangeMatcher = []int64{int64(rStart), int64(rEnd)}
    matcher.matcherType = rangeMatcher

  } else if tokenStr == "\\.\\." {
    if !isKey { return nil, errors.New("Invalid path value '..'. Use '\\.\\.' to match characters.") }
    matcher.matcherType = parentMatcher

  } else if tokenStr == ".*.*" {
    if !isKey { return nil, errors.New("Invalid path value '**'. Use '\\*\\*' to match characters.") }
    matcher.matcherType = recursiveMatcher

  } else if tokenStr == ".*" {
    matcher.matcherType = anyMatcher

  } else {
    matcher.regexpMatcher, err = regexp.Compile("^" + tokenStr + "$")
    matcher.matcherType = stringMatcher
  }

  return matcher, err
}

func runeToAscii(r rune) string {
    if r < 128 {
        return string(r)
    } else {
        return "\\u" + strconv.FormatInt(int64(r), 16)
    }
}
