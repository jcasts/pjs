package jsonfmt

import (
  "encoding/json"
  "io"
  "reflect"
  "../iterator"
)


type OrderedEncoder struct {
  datas []interface{}
  dataIndex int
  buffer []byte
}

func NewOrderedEncoder(datas ...interface{}) *OrderedEncoder {
  return &OrderedEncoder{datas, 0, []byte{}}
}

func (m* OrderedEncoder) Read(p []byte) (n int, err error) {
  bytesToRead := len(p)
  bytesRead := len(m.buffer)

  // TODO: Check if there's a way we can encode less data at once to make it more stream-friendly
  for bytesRead < bytesToRead && m.dataIndex < len(m.datas) {
    encoded, err := encode(m.datas[m.dataIndex])
    if err != nil { return bytesRead, err }

    m.buffer = append(m.buffer, encoded...)
    bytesRead = len(m.buffer)
    m.dataIndex++

    if m.dataIndex < len(m.datas) {
      m.buffer = append(m.buffer, byte('\n'))
      bytesRead++
    }
  }

  if bytesRead > bytesToRead {
    copy(p, m.buffer[0:bytesToRead])
    m.buffer = m.buffer[bytesToRead:]
    return bytesToRead, nil
  } else {
    copy(p, m.buffer)
    m.buffer = []byte{}
    return bytesRead, io.EOF
  }
}


func encode(data interface{}) ([]byte, error) {
  it, err := iterator.NewSortedDataIterator(data)
  if err != nil { return json.Marshal(data) }

  b := []byte{}

  value := reflect.ValueOf(data)
  isMap := value.Kind() != reflect.Slice && value.Kind() != reflect.Array

  if isMap {
    b = append(b, byte('{'))
  } else {
    b = append(b, byte('['))
  }

  for it.Next() {
    entry := it.Value()
    if isMap {
      key, err := json.Marshal(entry.Name)
      if err != nil { return b, err }
      b = append(b, key...)
      b = append(b, byte(':'))
    }
    bytes, err := encode(entry.Value)
    if err != nil { return b, err }
    b = append(b, bytes...)
    if !it.IsLast() {
      b = append(b, byte(','))
    }
  }

  if isMap {
    b = append(b, byte('}'))
  } else {
    b = append(b, byte(']'))
  }

  return b, nil
}
