package jsonfmt

import (
  "encoding/json"
  "io"
  "../iterator"
)


type OrderedEncoder struct {
  datas []interface{}
  dataIndex int
  buffer []byte
  iterators []*iterator.DataIterator
}

func NewOrderedEncoder(datas ...interface{}) *OrderedEncoder {
  return &OrderedEncoder{datas, 0, []byte{}, []*iterator.DataIterator{}}
}

func (m *OrderedEncoder) Read(p []byte) (n int, err error) {
  bytesToRead := len(p)
  bytesRead := len(m.buffer)
  count := 0

  for bytesRead < bytesToRead && m.dataIndex < len(m.datas) && count < 100 {
    var bytes []byte
    var err error

    if len(m.iterators) == 0 {
      bytes, err, _ = m.encodeData(m.datas[m.dataIndex])
    } else {
      bytes, err, _ = m.encode(m.iterators[len(m.iterators)-1])
    }
    if err != nil { return bytesRead, err }

    m.buffer = append(m.buffer, bytes...)
    bytesRead = len(m.buffer)
    if len(m.iterators) == 0 {
      m.dataIndex++
      if m.dataIndex < len(m.datas) {
        m.buffer = append(m.buffer, byte('\n'))
        bytesRead++
      }
    }
    count++
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

func (m *OrderedEncoder) encodeData(data interface{}) ([]byte, error, bool) {
  nextIt, err := iterator.NewSortedDataIterator(data)
  if err != nil {
    return jsonEncode(data)
  } else {
    return m.encode(nextIt)
  }
}

func (m *OrderedEncoder) encode(it *iterator.DataIterator) ([]byte, error, bool) {
  b := []byte{}
  isMap := it.HasNamedKeys()

  if !it.Next() {
    m.iterators = m.iterators[0:len(m.iterators)-1]
    if isMap {
      b = append(b, byte('}'))
    } else {
      b = append(b, byte(']'))
    }
    if len(m.iterators) > 0 && !m.iterators[len(m.iterators)-1].IsLast() {
      b = append(b, byte(','))
    }
    return b, nil, false
  }

  entry := it.Value()

  if it.IsFirst() {
    if isMap {
      b = append(b, byte('{'))
    } else {
      b = append(b, byte('['))
    }
    m.iterators = append(m.iterators, it)
  }

  if isMap {
    key, err := json.Marshal(entry.Name)
    if err != nil { return b, err, false }
    b = append(b, key...)
    b = append(b, byte(':'))
  }

  bytes, err, eod := m.encodeData(entry.Value)
  if err != nil { return b, err, false }
  b = append(b, bytes...)

  if !it.IsLast() && eod {
    b = append(b, byte(','))
  }

  return b, err, false
}


func jsonEncode(data interface{}) ([]byte, error, bool) {
  fl, ok := data.(float64)
  var bytes []byte
  var err error
  // Hack to compensate for the fact that Go converts all undefined numbers into float64
  if ok && float64(int64(fl)) == fl {
    bytes, err = json.Marshal(int64(fl))
  } else {
    bytes, err = json.Marshal(data)
  }
  return bytes, err, true
}
