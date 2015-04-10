package iterator

import (
  "fmt"
  "errors"
  "reflect"
  "sort"
)


type DataEntry struct {
  Index   int
  Name    string
  Key     interface{}
  Value   interface{}
}

type DataIterator struct {
  current   int
  keyCount  int
  keys      []reflect.Value
  data      reflect.Value
}


type itValueSorter []reflect.Value
func (v itValueSorter) Len() int { return len(v) }
func (v itValueSorter) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v itValueSorter) Less(i, j int) bool {
  // We're only concerned with JSON for now, so don't get too picky about
  // comparing types not in the JSON spec.
  if success, less := compareNumeric(v[i], v[j]); success {
    return less
  }

  v1 := fmt.Sprintf("%v", v[i].Interface())
  v2 := fmt.Sprintf("%v", v[j].Interface())
  return sort.StringsAreSorted([]string{v1, v2})
}

func compareNumeric(v1, v2 reflect.Value) (success bool, less bool) {
  intKinds := map[reflect.Kind]bool{
    reflect.Int: true,
    reflect.Int8: true,
    reflect.Int16: true,
    reflect.Int32: true,
    reflect.Int64: true,
  }
  uintKinds := map[reflect.Kind]bool{
    reflect.Uint: true,
    reflect.Uint8: true,
    reflect.Uint16: true,
    reflect.Uint32: true,
    reflect.Uint64: true,
  }
  floatKinds := map[reflect.Kind]bool{
    reflect.Float32: true,
    reflect.Float64: true,
  }

  if intKinds[v1.Kind()] && intKinds[v2.Kind()] {
    return true, v1.Int() < v2.Int()
  } else if intKinds[v1.Kind()] && uintKinds[v2.Kind()] {
    return true, v1.Int() < int64(v2.Uint())
  } else if intKinds[v1.Kind()] && floatKinds[v2.Kind()] {
    return true, v1.Int() < int64(v2.Float())
  } else if uintKinds[v1.Kind()] && intKinds[v2.Kind()] {
    return true, int64(v1.Uint()) < v2.Int()
  } else if uintKinds[v1.Kind()] && uintKinds[v2.Kind()] {
    return true, v1.Uint() < v2.Uint()
  } else if uintKinds[v1.Kind()] && floatKinds[v2.Kind()] {
    return true, v1.Uint() < uint64(v2.Float())
  } else if floatKinds[v1.Kind()] && intKinds[v2.Kind()] {
    return true, int64(v1.Float()) < v2.Int()
  } else if floatKinds[v1.Kind()] && uintKinds[v2.Kind()] {
    return true, uint64(v1.Float()) < v2.Uint()
  } else if floatKinds[v1.Kind()] && floatKinds[v2.Kind()] {
    return true, v1.Float() < v2.Float()
  } else {
    return false, false
  }
}


func NewSortedDataIterator(data interface{}) (d *DataIterator, err error) {
  d, err = NewDataIterator(data)
  if err == nil && len(d.keys) > 0 &&
      d.data.Kind() != reflect.Slice && d.data.Kind() != reflect.Array {
    sort.Sort(itValueSorter(d.keys))
  }
  return
}

func NewDataIterator(data interface{}) (d *DataIterator, err error) {
  val := reflect.ValueOf(data)
  if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
    val = val.Elem()
  }

  d = &DataIterator{current: -1, data: val}

  switch d.data.Kind() {
  case reflect.Struct:
    d.keys = deepGetStructFields(data)
    d.keyCount = len(d.keys)
  case reflect.Map:
    d.keys = d.data.MapKeys()
    d.keyCount = len(d.keys)
  case reflect.Slice, reflect.Array:
    d.keyCount = d.data.Len()
  default:
    // Not a traversable structure
    err = errors.New(fmt.Sprintf("Non-iteratable data structure %v", data))
    d = nil
  }
  return
}

func deepGetStructFields(data interface{}) (fieldValues []reflect.Value) {
  fieldValues = []reflect.Value{}
  val := reflect.ValueOf(data)
  if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
    val = val.Elem()
  }

  if val.Kind() != reflect.Struct { return }

  for i := 0; i < val.NumField(); i++ {
    if val.Field(i).CanInterface() {
      if val.Type().Field(i).Anonymous {
        fieldValues = append(fieldValues, deepGetStructFields(val.Field(i).Interface())...)
      } else {
        fieldValues = append(fieldValues, reflect.ValueOf(val.Type().Field(i).Name))
      }
    }
  }

  return
}

func (d *DataIterator) Reset() {
  d.current = -1
}

func (d *DataIterator) IsLast() bool {
  return d.current == d.keyCount - 1
}

func (d *DataIterator) Next() bool {
  d.current++
  return d.current < d.keyCount
}

func (d *DataIterator) Value() (de *DataEntry) {
  if d.current >= d.keyCount { return nil }
  de = &DataEntry{Index: d.current}

  switch d.data.Kind() {
  case reflect.Struct:
    // Private fields are skipped in the constructor function
    key := d.keys[d.current]
    de.Name = fmt.Sprintf("%v", key.Interface())
    de.Key = de.Name
    de.Value = d.data.FieldByName(de.Name).Interface()
  case reflect.Map:
    // This is build specifically for JSON which can't have
    // anything other than strings as a map key
    key := d.keys[d.current]
    de.Name = fmt.Sprintf("%v", key.Interface())
    de.Key = de.Name
    de.Value = d.data.MapIndex(key).Interface()
  case reflect.Slice, reflect.Array:
    de.Value = d.data.Index(d.current).Interface()
    de.Key = d.current
  default:
    // Not a traversable structure
    return nil
  }
  return
}