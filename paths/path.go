package paths

import (
  "fmt"
  "regexp"
  "reflect"
)


type Path interface {
  String() string
  FindMatches(data interface{}) []*PathMatch
}


type path struct {
  raw string
  tokens []*pathToken
}

func (p *path) String() string {
  return p.raw
}

func (p *path) FindMatches(data interface{}) []*PathMatch {
  if it, err := newDataIterator(data); err == nil {
    return p.walkData(it, nil, 0)
  }
  return []*PathMatch{}
}

func (p *path) walkData(it *dataIterator, parent *PathMatch, pathDepth int) (pathMatches []*PathMatch) {
  if pathDepth >= len(p.tokens) { return nil }
  token := p.tokens[pathDepth]
  pathMatches = []*PathMatch{}

  for it.Next() {
    entry := it.Value()
    if entry != nil && token.matches(entry.key, entry.value) {
      match := &PathMatch{Key: entry.key, Value: entry.value, ParentMatch: parent}
      if entry.iterator != nil {
        match.ChildMatches = p.walkData(entry.iterator, match, pathDepth + 1)
      }
      if len(match.ChildMatches) > 0 || pathDepth == len(p.tokens) - 1 {
        pathMatches = append(pathMatches, match)
      }
    }
  }

  return
}

type PathMatch struct {
  Key interface{}
  Value interface{}
  ChildMatches []*PathMatch
  ParentMatch *PathMatch
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
  boolMatcher
  rangeMatcher 
  parentMatcher
  anyMatcher
)

type tokenMatcher struct {
  regexpMatcher *regexp.Regexp
  rangeMatcher []int64
  intMatcher int64
  boolMatcher bool
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
