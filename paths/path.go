package paths

import (
  "fmt"
  "regexp"
  "reflect"
  "github.com/jcasts/pjs/iterator"
)


type Path interface {
  String() string
  FindMatches(data interface{}) Matches
}

type path struct {
  raw string
  tokens []*pathToken
}

func NewPath(pathStr string) (Path, error) {
  return parsePath(pathStr)
}

func (p *path) String() string {
  return p.raw
}

func (p *path) FindMatches(data interface{}) (matchSets Matches) {
  matchKeys := []string{""}
  uniqMatchSets := map[string]Match{"": NewMatch(data)}
  for _, token := range p.tokens {
    newMatchKeys := []string{}
    newMatchSets := map[string]Match{}
    for _, key := range matchKeys {
      matchSet := uniqMatchSets[key]
      results := matchPathToken(token, matchSet)
      for _, pathMatch := range results {
        newKey := pathMatch.hashId()
        if _, ok := newMatchSets[newKey]; !ok {
          newMatchKeys = append(newMatchKeys, newKey)
        }
        newMatchSets[newKey] = pathMatch
      }
    }
    matchKeys = newMatchKeys
    uniqMatchSets = newMatchSets
    if len(uniqMatchSets) == 0 { return }
  }

  matchSets = []Match{}
  for _, key := range matchKeys {
    matchSets = append(matchSets, uniqMatchSets[key])
  }
  return
}

func matchPathToken(token *pathToken, dataSet Match) (dataSets []Match) {
  dataSets = []Match{}
  if dataSet.Length() == 0 { return }

  if token.followParent() {
    dataSets = append(dataSets, dataSet.CopyAndTrim(1))
    return
  }

  it, err := iterator.NewDataIterator(dataSet.Value())
  if err != nil { return }
  counter := 0
  matchCounter := 0

  // Handle matching, recursion
  for it.Next() {
    entry := it.Value()
    if entry == nil { continue }

    counter++
    if token.matches(entry.Key(), entry.Interface()) {
      matchCounter++
      var lastDataSet Match
      if (len(dataSets) > 0) { lastDataSet = dataSets[len(dataSets)-1] }
      newDataSet := dataSet.CopyAndAppend(entry.Key(), entry.Interface())
      if token.isRecursive() && token.inverseMatcher() && entry.HasIterator() {
        dataSets = append(dataSets, matchPathToken(token, newDataSet)...)
      } else if &lastDataSet == nil || !lastDataSet.Equals(newDataSet) {
        dataSets = append(dataSets, newDataSet)
      }

    } else if token.isRecursive() && !token.inverseMatcher() {
      matchCounter++
      newDataSet := dataSet.CopyAndAppend(entry.Key(), entry.Interface())
      dataSets = append(dataSets, matchPathToken(token, newDataSet)...)
    }
  }

  if token.isInverseChildMatch() && matchCounter < counter {
    dataSets = []Match{}
  }

  return
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

func (t *pathToken) inverseMatcher() bool {
  return t.keyMatcher.inverse
}

func (t *pathToken) isInverseChildMatch() bool {
  return t.keyMatcher.inverse && t.keyMatcher.exclusive
}


type tokenMatcherType int
const (
  stringMatcher  tokenMatcherType = iota
  intMatcher
  boolMatcher
  rangeMatcher 
  parentMatcher
  anyMatcher
  nilMatcher
)

type tokenMatcher struct {
  regexpMatcher *regexp.Regexp
  rangeMatcher []int64
  intMatcher int64
  boolMatcher bool
  matcherType tokenMatcherType
  recursive bool
  inverse bool
  exclusive bool
}

func (tm *tokenMatcher) matches(value interface{}) bool {
  if tm.inverse {
    return !tm.matchesInterface(value)
  } else {
    return tm.matchesInterface(value)
  }
}

func (tm *tokenMatcher) matchesInterface(value interface{}) bool {
  if tm.matcherType == anyMatcher {
    return true
  }

  if tm.matcherType == stringMatcher {
    return tm.regexpMatcher != nil && tm.regexpMatcher.MatchString(fmt.Sprintf("%v", value))
  }

  if value == nil && tm.matcherType == nilMatcher {
    return true
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
