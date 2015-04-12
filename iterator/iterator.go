package iterator

import (
  "fmt"
  "errors"
  "reflect"
  "sort"
)


type Iterator interface {
  Next() bool
  Value() Value
  HasNamedKeys() bool
  IsFirst() bool
  IsLast() bool
}

type Value interface {
  Index() int
  Name() string
  Key() interface{}
  Interface() interface{}
  HasIterator() bool
  Iterator() Iterator
}

type DataValue struct {
  index   int
  name    string
  key     interface{}
  value   interface{}
  iterator *DataIterator
}

func NewDataValue(data interface{}, sorted bool) *DataValue {
  var it *DataIterator
  if sorted {
    it, _ = NewSortedDataIterator(data)
  } else {
    it, _ = NewDataIterator(data)
  }
  return &DataValue{
    index: 0,
    name: "",
    key: nil,
    value: data,
    iterator: it,
  }
}
func (de *DataValue) Index() int { return de.index }
func (de *DataValue) Key() interface{} { return de.key }
func (de *DataValue) Name() string { return de.name }
func (de *DataValue) Interface() interface{} { return de.value }
func (de *DataValue) Iterator() Iterator { return de.iterator }
func (de *DataValue) HasIterator() bool { return de.iterator != nil }


type DataIterator struct {
  current   int
  keyCount  int
  keys      []reflect.Value
  data      reflect.Value
  sorted    bool
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
    d.sorted = true
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

func (d *DataIterator) HasNamedKeys() bool {
  return d.data.Kind() != reflect.Slice && d.data.Kind() != reflect.Array
}

func (d *DataIterator) IsFirst() bool {
  return d.current == 0
}

func (d *DataIterator) IsLast() bool {
  return d.current == d.keyCount - 1
}

func (d *DataIterator) Next() bool {
  d.current++
  return d.current < d.keyCount
}

func (d *DataIterator) Value() Value {
  if d.current >= d.keyCount { return nil }
  de := &DataValue{index: d.current}

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

  if d.sorted {
    de.iterator, _ = NewSortedDataIterator(de.value)
  } else {
    de.iterator, _ = NewDataIterator(de.value)
  }

  return de
}