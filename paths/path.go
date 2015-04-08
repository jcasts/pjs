package paths

import (
  "encoding/json"
  "fmt"
  "regexp"
  "reflect"
  "strings"
)


type Path interface {
  String() string
  FindMatches(data interface{}) PathMatches
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

func (p *path) FindMatches(data interface{}) (matchSets PathMatches) {
  matchKeys := []string{""}
  uniqMatchSets := map[string]PathMatch{"": NewPathMatch(data)}
  for _, token := range p.tokens {
    newMatchKeys := []string{}
    newMatchSets := map[string]PathMatch{}
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

  matchSets = []PathMatch{}
  for _, key := range matchKeys {
    matchSets = append(matchSets, uniqMatchSets[key])
  }
  return
}

func matchPathToken(token *pathToken, dataSet PathMatch) (dataSets []PathMatch) {
  dataSets = []PathMatch{}
  if dataSet.Length() == 0 { return }

  if token.followParent() {
    dataSets = append(dataSets, dataSet.CopyAndTrim(1))
    return
  }

  it, err := newDataIterator(dataSet.Value())
  if err != nil { return }

  // Handle matching, recursion
  for it.Next() {
    entry := it.Value()
    if entry == nil { continue }

    if token.matches(entry.key, entry.value) {
      var lastDataSet PathMatch
      if (len(dataSets) > 0) { lastDataSet = dataSets[len(dataSets)-1] }
      newDataSet := dataSet.CopyAndAppend(entry.key, entry.value)
      if &lastDataSet == nil || !lastDataSet.Equals(newDataSet) {
        dataSets = append(dataSets, newDataSet)
      }

    } else if token.isRecursive() {
      newDataSet := dataSet.CopyAndAppend(entry.key, entry.value)
      dataSets = append(dataSets, matchPathToken(token, newDataSet)...)
    }
  }

  return
}


type PathMatches []PathMatch

func (pms PathMatches) MarshalJSON() (bytes []byte, err error) {
  return json.Marshal(pms.buildJsonStruct())
}

func (pms PathMatches) buildJsonStruct() interface{} {
  // TODO
  return nil
}


type PathMatch struct {
  nodes []*DataNode
  hashes []string
}

func (pm PathMatch) hashId() string {
  return strings.Join(pm.hashes, ":")
}

func (pm PathMatch) Length() int {
  return len(pm.nodes)
}

func (pm PathMatch) NodeAt(index int) *DataNode {
  return pm.nodes[index]
}

func (pm PathMatch) CopyAndTrim(rmSize int) PathMatch {
  length := pm.Length() - rmSize
  if (length <= 0) { length = 1}
  return PathMatch{
    nodes: pm.nodes[0:length],
    hashes: pm.hashes[0:length],
  }
}

func (pm PathMatch) CopyAndAppend(key, value interface{}) PathMatch {
  node := &DataNode{Key: key, Value: value}
  nodes := append([]*DataNode{}, pm.nodes...)
  nodes = append(nodes, node)
  hashes := append([]string{}, pm.hashes...)
  hashes = append(hashes, fmt.Sprintf("%v", node.Key))
  return PathMatch{
    nodes: nodes,
    hashes: hashes,
  }
}

func (pm PathMatch) Equals(other PathMatch) bool {
  return pm.hashId() == other.hashId()
}

func (pm PathMatch) Value() interface{} {
  return pm.nodes[pm.Length()-1].Value
}

func NewPathMatch(value interface{}) PathMatch {
  return PathMatch{
    nodes: []*DataNode{&DataNode{Value: value}},
    hashes: []string{""},
  }
}


type DataNode struct {
  Key interface{}
  Value interface{}
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
