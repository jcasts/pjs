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