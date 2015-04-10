package jsonfmt

import (
  "encoding/json"
  "io"
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
    "password": "iloveu😻",
    "address": map[string]interface{}{
      "street": "1 Infinite Loop",
      "city": "Cupertino",
      "zip": "91234",
      "pos": []string{"Apple", "HQ"},
    },
    "roles": []string{"eng", "employee"},
  }
}


func TestEncodeMapData(t *testing.T) {
  data := mockMapData()
  enc := NewOrderedEncoder(data)
  b := make([]byte, 1024)
  num, err := enc.Read(b)

  jstr := string(b[0:num])
  expected := "{\"address\":{\"city\":\"Cupertino\",\"pos\":[\"Apple\",\"HQ\"],\"street\":"+
    "\"1 Infinite Loop\",\"zip\":\"91234\"},\"age\":30,\"name\":\"Bob\",\"password\":\"iloveu😻\""+
    ",\"roles\":[\"eng\",\"employee\"]}"

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expected, jstr)

  var v interface{}
  err = json.Unmarshal([]byte(jstr), &v)
  testAssertNil(t, err)
}


func TestEncodeStructData(t *testing.T) {
  data := mockStructData()
  enc := NewOrderedEncoder(data)
  b := make([]byte, 1024)
  num, err := enc.Read(b)

  jstr := string(b[0:num])
  expected := `{"Address":{"City":"Cupertino","Street":"1 Infinite Loop","Zip":"91234"},"Age":30,"Name":"Bob"}`

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expected, jstr)

  var v interface{}
  err = json.Unmarshal([]byte(jstr), &v)
  testAssertNil(t, err)
}

func TestEncodeMultiData(t *testing.T) {
  enc := NewOrderedEncoder(mockStructData(), mockMapData())
  b := make([]byte, 1024)
  num, err := enc.Read(b)

  jstr := string(b[0:num])
  expected := `{"Address":{"City":"Cupertino","Street":"1 Infinite Loop","Zip":"91234"},"Age":30,"Name":"Bob"}`+
  "\n{\"address\":{\"city\":\"Cupertino\",\"pos\":[\"Apple\",\"HQ\"],\"street\":"+
    "\"1 Infinite Loop\",\"zip\":\"91234\"},\"age\":30,\"name\":\"Bob\",\"password\":\"iloveu😻\""+
    ",\"roles\":[\"eng\",\"employee\"]}"

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expected, jstr)
}

func TestEncodeWithTinyBuffer(t *testing.T) {
  enc := NewOrderedEncoder(mockStructData())
  b := make([]byte, 32)

  num, err := enc.Read(b)
  jstr := string(b[0:num])
  expected := `{"Address":{"City":"Cupertino","`

  testAssertNil(t, err)
  testAssertEqual(t, expected, jstr)

  num, err = enc.Read(b)
  jstr += string(b[0:num])
  expected += `Street":"1 Infinite Loop","Zip":`

  testAssertNil(t, err)
  testAssertEqual(t, expected, jstr)

  num, err = enc.Read(b)
  jstr += string(b[0:num])
  expected += `"91234"},"Age":30,"Name":"Bob"}`

  testAssertEqual(t, io.EOF, err)
  testAssertEqual(t, expected, jstr)
}
