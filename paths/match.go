package paths

import (
  "fmt"
  "reflect"
  "sort"
  "strings"
  "../iterator"
)

type matchesIteratorValue struct {
  index int
  name string
  key interface{}
  value interface{}
  hasNamedKeys bool
  childKeys []reflect.Value
  children map[string]*matchesIteratorValue
  iterator *matchesIterator
}

func newMatchesSortedIteratorValue(matches Matches) *matchesIteratorValue {
  if len(matches) == 0 { return nil }

  itValue := &matchesIteratorValue{
    childKeys: []reflect.Value{},
    children: map[string]*matchesIteratorValue{},
  }

  for _, m := range matches {
    v := itValue

    for _, node := range m.nodes {
      key := fmt.Sprintf("%v", node.Key)
      if _, ok := v.children[key]; !ok {
        v.children[key] = &matchesIteratorValue{
          childKeys: []reflect.Value{},
          children: map[string]*matchesIteratorValue{},
        }
        v.childKeys = append(v.childKeys, reflect.ValueOf(node.Key))
      }
      child := v.children[key]
      val := reflect.ValueOf(node.Value)
      if val.Kind() == reflect.Struct || val.Kind() == reflect.Map {
        child.hasNamedKeys = true
      }
      child.index, _ = node.Key.(int)
      child.name = key
      child.value = node.Value
      child.key = node.Key

      v = child
    }
  }
  return itValue.childForKey(itValue.childKeys[0])
}

func (v *matchesIteratorValue) Index() int {
  return v.index
}

func (v *matchesIteratorValue) Name() string {
  return v.name
}

func (v *matchesIteratorValue) Key() interface{} {
  return v.key
}

func (v *matchesIteratorValue) Interface() interface{} {
  return v.value
}

func (v *matchesIteratorValue) HasIterator() bool {
  return len(v.children) > 0
}

func (v *matchesIteratorValue) Iterator() iterator.Iterator {
  return newSortedMatchesIterator(v)
}

func (v *matchesIteratorValue) childForKey(key reflect.Value) *matchesIteratorValue {
  var k interface{}
  if key.IsValid() { k = key.Interface() }
  return v.children[fmt.Sprintf("%v", k)]
}

type matchesIterator struct {
  value *matchesIteratorValue
  keys []reflect.Value
  sort bool
  current int
}

func newSortedMatchesIterator(v *matchesIteratorValue) *matchesIterator {
  keys := v.childKeys
  sort.Sort(iterator.ValueSorter(keys))
  return &matchesIterator{current: -1, value: v, sort: true, keys: keys}
}

func (i *matchesIterator) Next() bool {
  i.current++
  return i.current < len(i.value.children)
}

func (i *matchesIterator) Value() iterator.Value {
  if i.current >= len(i.keys) { return nil }
  key := i.keys[i.current]
  return i.value.childForKey(key)
}

func (i *matchesIterator) HasNamedKeys() bool {
  return i.value.hasNamedKeys
}

func (i *matchesIterator) IsFirst() bool {
  return i.current == 0
}

func (i *matchesIterator) IsLast() bool {
  return i.current == len(i.value.children) - 1
}


type Matches []Match

func (pms Matches) IteratorValue() iterator.Value {
  itValue := newMatchesSortedIteratorValue(pms)
  if itValue == nil { return nil } // Interface to pointer to nil not recognized as nil
  return itValue
}


type DataNode struct {
  Key interface{}
  Value interface{}
}


type Match struct {
  nodes []*DataNode
  hashes []string
  hashIdStr string
}

func (pm Match) hashId() string {
  return pm.hashIdStr
}

func (pm Match) Length() int {
  return len(pm.nodes)
}

func (pm Match) NodeAt(index int) *DataNode {
  return pm.nodes[index]
}

func (pm Match) CopyAndTrim(rmSize int) Match {
  length := pm.Length() - rmSize
  if (length <= 0) { length = 1}
  return Match{
    nodes: pm.nodes[0:length],
    hashes: pm.hashes[0:length],
    hashIdStr: strings.Join(pm.hashes[0:length], ":"),
  }
}

func (pm Match) CopyAndAppend(key, value interface{}) Match {
  node := &DataNode{Key: key, Value: value}
  nodes := append([]*DataNode{}, pm.nodes...)
  nodes = append(nodes, node)
  hashes := append([]string{}, pm.hashes...)
  hashes = append(hashes, fmt.Sprintf("%v", node.Key))
  return Match{
    nodes: nodes,
    hashes: hashes,
    hashIdStr: strings.Join(hashes, ":"),
  }
}

func (pm Match) Equals(other Match) bool {
  return pm.hashId() == other.hashId()
}

func (pm Match) Value() interface{} {
  return pm.nodes[pm.Length()-1].Value
}

func NewMatch(value interface{}) Match {
  return Match{
    nodes: []*DataNode{&DataNode{Value: value}},
    hashes: []string{""},
    hashIdStr: "",
  }
}
