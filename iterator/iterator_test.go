package iterator

import (
  "reflect"
  "testing"
)

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

  it, err := NewSortedDataIterator(data)
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.Index)
  testAssertEqual(t, "MyField0", val.Name)
  testAssertEqual(t, "MyField0", val.Key)
  testAssertEqual(t, "F0", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.Index)
  testAssertEqual(t, "MyField1", val.Name)
  testAssertEqual(t, "MyField1", val.Key)
  testAssertEqual(t, "F1", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.Index)
  testAssertEqual(t, "MyField2", val.Name)
  testAssertEqual(t, "MyField2", val.Key)
  testAssertEqual(t, "F2", val.Value)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorStruct(t *testing.T) {
  it, err := NewSortedDataIterator(mockStructData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.Index)
  testAssertEqual(t, "Address", val.Name)
  testAssertEqual(t, "Address", val.Key)
  testAssertEqual(t, reflect.TypeOf(&address{}), reflect.TypeOf(val.Value))

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.Index)
  testAssertEqual(t, "Age", val.Name)
  testAssertEqual(t, "Age", val.Key)
  testAssertEqual(t, 30, val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.Index)
  testAssertEqual(t, "Name", val.Name)
  testAssertEqual(t, "Name", val.Key)
  testAssertEqual(t, "Bob", val.Value)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorMap(t *testing.T) {
  it, err := NewSortedDataIterator(mockMapData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.Index)
  testAssertEqual(t, "address", val.Name)
  testAssertEqual(t, "address", val.Key)
  testAssertEqual(t, reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(val.Value))

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.Index)
  testAssertEqual(t, "age", val.Name)
  testAssertEqual(t, "age", val.Key)
  testAssertEqual(t, 30, val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.Index)
  testAssertEqual(t, "name", val.Name)
  testAssertEqual(t, "name", val.Key)
  testAssertEqual(t, "Bob", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.Index)
  testAssertEqual(t, "password", val.Name)
  testAssertEqual(t, "password", val.Key)
  testAssertEqual(t, "iloveu", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 4, val.Index)
  testAssertEqual(t, "roles", val.Name)
  testAssertEqual(t, "roles", val.Key)
  testAssertEqual(t, reflect.TypeOf([]string{}), reflect.TypeOf(val.Value))

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorSlice(t *testing.T) {
  it, err := NewSortedDataIterator(mockSliceData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.Index)
  testAssertEqual(t, "", val.Name)
  testAssertEqual(t, 0, val.Key)
  testAssertEqual(t, "Bob", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.Index)
  testAssertEqual(t, "", val.Name)
  testAssertEqual(t, 1, val.Key)
  testAssertEqual(t, 30, val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.Index)
  testAssertEqual(t, "", val.Name)
  testAssertEqual(t, 2, val.Key)
  testAssertEqual(t, "iloveu", val.Value)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.Index)
  testAssertEqual(t, "", val.Name)
  testAssertEqual(t, 3, val.Key)
  testAssertEqual(t, reflect.TypeOf([]interface{}{}), reflect.TypeOf(val.Value))
}
