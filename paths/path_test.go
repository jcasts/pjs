package paths

import(
  "fmt"
  "reflect"
  "sort"
  "testing"
)


func (pms PathMatches) Len() int { return len(pms) }
func (pms PathMatches) Swap(i, j int) { pms[i], pms[j] = pms[j], pms[i] }
func (pms PathMatches) Less(i, j int) bool {
  v1 := pms[i].hashId()
  v2 := pms[j].hashId()
  return sort.StringsAreSorted([]string{v1, v2})
}


func TestPathMatchStruct(t *testing.T) {
  var p *path
  var err error
  var matches PathMatches

  data := mockStructData()

  p, err = parsePath("foo")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 0, len(matches))

  p, err = parsePath("A*")
  testAssertNil(t, err)
  matches = p.FindMatches(mockStructData())
  sort.Sort(matches)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, 2, matches[0].Length())
  testAssertEqual(t, "Address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, reflect.TypeOf(data.Address), reflect.TypeOf(matches[0].NodeAt(1).Value))
  testAssertEqual(t, reflect.TypeOf(data.Address), reflect.TypeOf(matches[0].Value()))
  testAssertEqual(t, 2, matches[1].Length())
  testAssertEqual(t, "Age", matches[1].NodeAt(1).Key)
  testAssertEqual(t, 30, matches[1].NodeAt(1).Value)
  testAssertEqual(t, 30, matches[1].Value())

  p, err = parsePath("*/Zip|City")
  testAssertNil(t, err)
  matches = p.FindMatches(mockStructData())
  sort.Sort(matches)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, 3, matches[0].Length())
  testAssertEqual(t, "Address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "City", matches[0].NodeAt(2).Key)
  testAssertEqual(t, "Cupertino", matches[0].NodeAt(2).Value)
  testAssertEqual(t, "Cupertino", matches[0].Value())
  testAssertEqual(t, "Zip", matches[1].NodeAt(2).Key)
  testAssertEqual(t, "91234", matches[1].NodeAt(2).Value)
  testAssertEqual(t, "91234", matches[1].Value())
}


func TestPathMatchParent(t *testing.T) {
  var p *path
  var err error
  var matches PathMatches

  data := mockMapData()

  fmt.Println("===== P1")
  p, err = parsePath("*/zip|city/../../a*")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  fmt.Printf("%v\n", matches[1].NodeAt(1).Key)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, 2, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, 2, matches[1].Length())
  testAssertEqual(t, "age", matches[1].NodeAt(1).Key)

  fmt.Println("===== P2")
  p, err = parsePath("*/zip|city/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)

  fmt.Println("===== P3")
  p, err = parsePath("*/*/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "roles", matches[1].NodeAt(1).Key)

  fmt.Println("===== P4")
  p, err = parsePath("*/pos/*=Apple/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "pos", matches[0].NodeAt(2).Key)
}


func TestPathMatchRecursive(t *testing.T) {
  var p *path
  var err error
  var matches PathMatches

  data := mockMapData()

  fmt.Println("===== P0")
  p, err = parsePath("address/**=Apple")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 4, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "pos", matches[0].NodeAt(2).Key)
  testAssertEqual(t, 0, matches[0].NodeAt(3).Key)

  fmt.Println("===== P1")
  p, err = parsePath("**=Apple/../..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 2, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)

  fmt.Println("===== P2")
  p, err = parsePath("**/pos/*=Apple/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 3, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "pos", matches[0].NodeAt(2).Key)

  fmt.Println("===== P3")
  p, err = parsePath("**/*=Apple/../..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 2, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)

  fmt.Println("===== P4")
  p, err = parsePath("address/**=Apple/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 3, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "pos", matches[0].NodeAt(2).Key)

  fmt.Println("===== P5")
  p, err = parsePath("address/**/*=Apple/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  sort.Sort(matches)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, 3, matches[0].Length())
  testAssertEqual(t, "address", matches[0].NodeAt(1).Key)
  testAssertEqual(t, "pos", matches[0].NodeAt(2).Key)
}
