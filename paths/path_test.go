package paths

import(
  "testing"
  "reflect"
)


func TestPathMatchStruct(t *testing.T) {
  var p *path
  var err error
  var matches []*PathMatch

  data := mockStructData()

  p, err = parsePath("foo")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  testAssertEqual(t, 0, len(matches))

  p, err = parsePath("A*")
  testAssertNil(t, err)
  matches = p.FindMatches(mockStructData())
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, "Address", matches[0].Key)
  testAssertEqual(t, reflect.TypeOf(data.Address), reflect.TypeOf(matches[0].Value))
  testAssertEqual(t, 0, len(matches[0].ChildMatches))
  testAssertEqual(t, "Age", matches[1].Key)
  testAssertEqual(t, 30, matches[1].Value)

  p, err = parsePath("*/Zip|City")
  testAssertNil(t, err)
  matches = p.FindMatches(mockStructData())
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, "Address", matches[0].Key)
  matches = matches[0].ChildMatches
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, "City", matches[0].Key)
  testAssertEqual(t, "Cupertino", matches[0].Value)
  testAssertEqual(t, "Zip", matches[1].Key)
  testAssertEqual(t, "91234", matches[1].Value)
}


func TestPathMatchParent(t *testing.T) {
  var p *path
  var err error
  var matches []*PathMatch

  data := mockMapData()

  p, err = parsePath("*/zip|city/../a*")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, "address", matches[0].Key)
  testAssertEqual(t, "age", matches[1].Key)

  p, err = parsePath("*/zip|city/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, "address", matches[0].Key)

  p, err = parsePath("*/*/..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  testAssertEqual(t, 2, len(matches))
  testAssertEqual(t, "address", matches[0].Key)
  testAssertEqual(t, "roles", matches[1].Key)

  p, err = parsePath("*/pos/*=Apple/../..")
  testAssertNil(t, err)
  matches = p.FindMatches(data)
  testAssertEqual(t, 1, len(matches))
  testAssertEqual(t, "address", matches[0].Key)
}


func TestPathMatchRecursive(t *testing.T) {

}
