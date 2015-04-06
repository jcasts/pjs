package paths

import (
  "fmt"
  "errors"
  "reflect"
  "sort"
)


type dataEntry struct {
  index   int
  name    string
  key     interface{}
  value   interface{}
}

type dataIterator struct {
  current   int
  keyCount  int
  keys      []reflect.Value
  data      reflect.Value
}

type valueSorter []reflect.Value
func (v valueSorter) Len() int { return len(v) }
func (v valueSorter) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v valueSorter) Less(i, j int) bool {
  v1 := fmt.Sprintf("%v", v[i].Interface())
  v2 := fmt.Sprintf("%v", v[j].Interface())
  return sort.StringsAreSorted([]string{v1, v2})
}

func newDataIterator(data interface{}) (d *dataIterator, err error) {
  val := reflect.ValueOf(data)
  if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
    val = val.Elem()
  }

  d = &dataIterator{current: -1, data: val}

  switch d.data.Kind() {
  case reflect.Struct:
    d.keys = []reflect.Value{}
    for i := 0; i < d.data.NumField(); i++ {
      if d.data.Field(i).CanInterface() {
        d.keys = append(d.keys, reflect.ValueOf(d.data.Type().Field(i).Name))
      }
    }
    d.keyCount = len(d.keys)
    sort.Sort(valueSorter(d.keys))
  case reflect.Map:
    d.keys = d.data.MapKeys()
    d.keyCount = len(d.keys)
    sort.Sort(valueSorter(d.keys))
  case reflect.Slice, reflect.Array:
    d.keyCount = d.data.Len()
  default:
    // Not a traversable structure
    err = errors.New(fmt.Sprintf("Non-iteratable data structure %v", data))
    d = nil
  }
  return
}

func (d *dataIterator) Reset() {
  d.current = -1
}

func (d *dataIterator) Next() bool {
  d.current++
  return d.current < d.keyCount
}

func (d *dataIterator) Value() (de *dataEntry) {
  if d.current >= d.keyCount { return nil }
  de = &dataEntry{index: d.current}

  switch d.data.Kind() {
  case reflect.Struct:
    // Private fields are skipped in the constructor function
    key := d.keys[d.current]
    de.name = fmt.Sprintf("%v", key.Interface())
    de.key = de.name
    de.value = d.data.FieldByName(de.name).Interface()
  case reflect.Map:
    // This is build specifically for JSON which can't have
    // anything other than strings as a map key
    key := d.keys[d.current]
    de.name = fmt.Sprintf("%v", key.Interface())
    de.key = de.name
    de.value = d.data.MapIndex(key).Interface()
  case reflect.Slice, reflect.Array:
    de.value = d.data.Index(d.current).Interface()
    de.key = d.current
  default:
    // Not a traversable structure
    return nil
  }
  return
}