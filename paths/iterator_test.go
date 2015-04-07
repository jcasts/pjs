package paths

import (
  "fmt"
  "reflect"
  "sort"
  "testing"
)

type itValueSorter []reflect.Value
func (v itValueSorter) Len() int { return len(v) }
func (v itValueSorter) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v itValueSorter) Less(i, j int) bool {
  v1 := fmt.Sprintf("%v", v[i].Interface())
  v2 := fmt.Sprintf("%v", v[j].Interface())
  return sort.StringsAreSorted([]string{v1, v2})
}

func newSortedDataIterator(data interface{}) (d *dataIterator, err error) {
  d, err = newDataIterator(data)
  if err == nil && len(d.keys) > 0 { sort.Sort(itValueSorter(d.keys)) }
  return
}

type StructTest1 struct {
  MyField1 string
  MyField2 string
}

type StructTest2 struct {
  StructTest1
  MyField0 string
}

func TestDataIteratorEmbeddedStruct(t *testing.T) {
  data := StructTest2{StructTest1{"F1", "F2"}, "F0"}

  it, err := newSortedDataIterator(data)
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "MyField0", val.name)
  testAssertEqual(t, "MyField0", val.key)
  testAssertEqual(t, "F0", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "MyField1", val.name)
  testAssertEqual(t, "MyField1", val.key)
  testAssertEqual(t, "F1", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "MyField2", val.name)
  testAssertEqual(t, "MyField2", val.key)
  testAssertEqual(t, "F2", val.value)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorStruct(t *testing.T) {
  it, err := newSortedDataIterator(mockStructData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "Address", val.name)
  testAssertEqual(t, "Address", val.key)
  testAssertEqual(t, reflect.TypeOf(&address{}), reflect.TypeOf(val.value))

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "Age", val.name)
  testAssertEqual(t, "Age", val.key)
  testAssertEqual(t, 30, val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "Name", val.name)
  testAssertEqual(t, "Name", val.key)
  testAssertEqual(t, "Bob", val.value)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorMap(t *testing.T) {
  it, err := newSortedDataIterator(mockMapData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "address", val.name)
  testAssertEqual(t, "address", val.key)
  testAssertEqual(t, reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(val.value))

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "age", val.name)
  testAssertEqual(t, "age", val.key)
  testAssertEqual(t, 30, val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "name", val.name)
  testAssertEqual(t, "name", val.key)
  testAssertEqual(t, "Bob", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.index)
  testAssertEqual(t, "password", val.name)
  testAssertEqual(t, "password", val.key)
  testAssertEqual(t, "iloveu", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 4, val.index)
  testAssertEqual(t, "roles", val.name)
  testAssertEqual(t, "roles", val.key)
  testAssertEqual(t, reflect.TypeOf([]string{}), reflect.TypeOf(val.value))

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorSlice(t *testing.T) {
  it, err := newSortedDataIterator(mockSliceData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 0, val.key)
  testAssertEqual(t, "Bob", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 1, val.key)
  testAssertEqual(t, 30, val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 2, val.key)
  testAssertEqual(t, "iloveu", val.value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 3, val.key)
  testAssertEqual(t, reflect.TypeOf([]interface{}{}), reflect.TypeOf(val.value))
}
