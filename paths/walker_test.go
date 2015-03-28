package paths

import (
  //"fmt"
  "reflect"
  "testing"
)

type address struct {
  Street string
  Zip string
  City string
}

type person struct {
  Name string
  Age int
  password string
  Address *address
}

func mockStructData() *person {
  addr := &address{Street: "1 Infinite Loop", Zip: "91234", City: "Cupertino"}
  return &person{Name: "Bob", Age: 30, Address: addr, password: "iloveu"}
}

func mockMapData() map[string]interface{} {
  return map[string]interface{}{
    "name": "Bob",
    "age": 30,
    "password": "iloveu",
    "address": map[string]interface{}{
      "street": "1 Infinite Loop",
      "city": "Cupertino",
      "zip": "91234",
    },
  }
}

func mockSliceData() []interface{} {
  return []interface{}{
    "Bob",
    30,
    "iloveu",
    []interface{}{
      "1 Infinite Loop",
      "cupertino",
      "91234",
    },
  }
}

func TestDataIteratorStruct(t *testing.T) {
  it, err := newDataIterator(mockStructData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "Address", val.name)
  testAssertEqual(t, "Address", val.key)
  testAssertEqual(t, reflect.TypeOf(&address{}), reflect.TypeOf(val.value))
  testAssertTrue(t, val.iterator != nil)
  testAssertEqual(t, *val.value.(*address), val.iterator.data.Interface())

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "Age", val.name)
  testAssertEqual(t, "Age", val.key)
  testAssertEqual(t, 30, val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "Name", val.name)
  testAssertEqual(t, "Name", val.key)
  testAssertEqual(t, "Bob", val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorMap(t *testing.T) {
  it, err := newDataIterator(mockMapData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "address", val.name)
  testAssertEqual(t, "address", val.key)
  testAssertEqual(t, reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(val.value))
  testAssertTrue(t, val.iterator != nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "age", val.name)
  testAssertEqual(t, "age", val.key)
  testAssertEqual(t, 30, val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "name", val.name)
  testAssertEqual(t, "name", val.key)
  testAssertEqual(t, "Bob", val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.index)
  testAssertEqual(t, "password", val.name)
  testAssertEqual(t, "password", val.key)
  testAssertEqual(t, "iloveu", val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertFalse(t, it.Next())
  testAssertTrue(t, it.Value() == nil)
}

func TestDataIteratorSlice(t *testing.T) {
  it, err := newDataIterator(mockSliceData())
  testAssertNil(t, err)

  testAssertTrue(t, it.Next())
  val := it.Value()
  testAssertEqual(t, 0, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 0, val.key)
  testAssertEqual(t, "Bob", val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 1, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 1, val.key)
  testAssertEqual(t, 30, val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 2, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 2, val.key)
  testAssertEqual(t, "iloveu", val.value)
  testAssertTrue(t, val.iterator == nil)

  testAssertTrue(t, it.Next())
  val = it.Value()
  testAssertEqual(t, 3, val.index)
  testAssertEqual(t, "", val.name)
  testAssertEqual(t, 3, val.key)
  testAssertEqual(t, reflect.TypeOf([]interface{}{}), reflect.TypeOf(val.value))
  testAssertTrue(t, val.iterator != nil)
}
