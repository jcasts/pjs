package paths

import (
  "fmt"
  "errors"
  "reflect"
)

func (p *path) matchData(data interface{}) (pathMatches []PathMatch) {
  pathMatches = []PathMatch{}

  it, err := newDataIterator(data)
  for err != nil && it.Next() {

  }

  return
}


type dataEntry struct {
  index int
  name  string
  value reflect.Value
}

type dataIterator struct {
  current   int
  keyCount  int
  data      reflect.Value
}

func newDataIterator(data interface{}) (d *dataIterator, err error) {
  val := reflect.ValueOf(data)
  if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
    val = val.Elem()
  }

  d = &dataIterator{current: -1, data: val}

  switch d.data.Kind() {
  case reflect.Struct:
    d.keyCount = d.data.Len()
  case reflect.Map:
    d.keyCount = len(d.data.MapKeys())
  case reflect.Slice, reflect.Array:
    d.keyCount = d.data.NumField()
  default:
    // Not a traversable structure
    err = errors.New(fmt.Sprintf("Non-traversable data structure %v", data))
  }
  return
}

func (d *dataIterator) Next() bool {
  d.current++
  if d.data.Kind() == reflect.Struct {
    for d.data.Type().Field(d.current).Anonymous && d.current < d.keyCount {
      d.current++
    }
  }
  return d.current >= d.keyCount
}

func (d *dataIterator) Value() (de *dataEntry) {
  de = &dataEntry{index: d.current}

  switch d.data.Kind() {
  case reflect.Struct:
    // Anonymous fields are skipped in the Next() method
    de.name = d.data.Type().Field(d.current).Name
    de.value = d.data.Field(d.current)
  case reflect.Map:
    key := d.data.MapKeys()[d.current]
    // This is build specifically for JSON which can't have
    // anything other than strings as a map key
    de.name = fmt.Sprintf("%v", key.Interface())
    de.value = d.data.MapIndex(key)
  case reflect.Slice, reflect.Array:
    de.value = d.data.Index(d.current)
  default:
    // Not a traversable structure
    return nil
  }
  return
}