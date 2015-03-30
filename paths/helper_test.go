package paths

import (
  "testing"
)

func testAssertEqual(t *testing.T, v1, v2 interface{}) {
  if v1 != v2 { t.Fatalf("Expected %v but was %v", v1, v2) }
}

func testAssertTrue(t *testing.T, value bool) {
  if !value { t.Fatalf("Expected value to be true") }
}

func testAssertFalse(t *testing.T, value bool) {
  if value { t.Fatalf("Expected value to be false") }
}

func testAssertNil(t *testing.T, value interface{}) {
  if value != nil { t.Fatalf("Expected value to be nil") }
}

func testAssertNotNil(t *testing.T, value interface{}) {
  if value == nil { t.Fatalf("Expected value to be not nil") }
}

func testParsePath(t *testing.T, pStr string) *path {
  p, err := parsePath(pStr)
  testAssertNil(t, err)
  return p
}


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
      "pos": []string{"Apple", "HQ"},
    },
    "roles": []string{"eng", "employee"},
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