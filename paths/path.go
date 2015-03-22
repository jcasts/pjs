package paths

import (
  //"errors"
  "fmt"
  "regexp"
  "reflect"
  //"strconv"
)


type Path interface {
  String() string
  Matches(data interface{}) (bool, [][]interface{})
}


type path struct {
  raw string
  tokens []*pathToken
}

func (p *path) String() string {
  return p.raw
}

func (p *path) Matches(data interface{}) (bool, [][]interface{}) {
  success := false
  //TODO: Implement
  return success, nil
}


type pathToken struct {
  keyMatcher *tokenMatcher
  valueMatcher *tokenMatcher
}

func (t *pathToken) matches(key interface{}, value interface{}) bool {
  return (t.keyMatcher == nil || t.keyMatcher.matches(key)) &&
    (t.valueMatcher == nil || t.valueMatcher.matches(value))
}

func (t *pathToken) followParent() bool {
  return t.keyMatcher.matcherType == parentMatcher
}

func (t *pathToken) isRecursive() bool {
  return t.keyMatcher.recursive
}

func (t *pathToken) isAny() bool {
  return t.keyMatcher.matcherType == anyMatcher &&
    (t.valueMatcher == nil || t.valueMatcher.matcherType == anyMatcher)
}


type tokenMatcherType int
const (
  stringMatcher  tokenMatcherType = iota
  intMatcher
  rangeMatcher 
  parentMatcher
  anyMatcher
)

type tokenMatcher struct {
  regexpMatcher *regexp.Regexp
  rangeMatcher []int64
  intMatcher int64
  matcherType tokenMatcherType
  recursive bool
}

func (tm *tokenMatcher) matches(value interface{}) bool {
  if tm.matcherType == anyMatcher {
    return true
  }

  if tm.matcherType == stringMatcher {
    return tm.regexpMatcher != nil && tm.regexpMatcher.MatchString(fmt.Sprintf("%v", value))
  }

  switch value.(type) {
  case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
    var num int64
    num = reflect.ValueOf(value).Convert(reflect.TypeOf(num)).Int()
    return tm.matchesInt(num)
  default:
    return tm.regexpMatcher != nil && tm.regexpMatcher.MatchString(fmt.Sprintf("%v", value))
  }
}

func (tm *tokenMatcher) matchesInt(num int64) bool {
  if tm.matcherType == intMatcher {
    return tm.intMatcher == num
  } else if tm.matcherType == rangeMatcher {
    return num >= tm.rangeMatcher[0] && num <= tm.rangeMatcher[1]
  } else {
    return false
  }
}


type dataNode struct {
  parent *dataNode
  key interface{}
  value interface{}
}


func NewPath(pathStr string) (Path, error) {
  return parsePath(pathStr)
}
