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
  // TODO: Have a way to request new partial data instead of instantiating
  // with a bunch of data existing items. E.g:
  // encoder.startMap(); encoder.addToMap("key", value); encoder.endMap()
  // encoder.startArray(); encoder.addToArray(value); encoder.endArray()
  // This would mean no ordering is really possible in maps.

  bytesToRead := len(p)
  bytesRead := len(m.buffer)

  for bytesRead < bytesToRead && m.dataIndex < len(m.datas) {
    var bytes []byte
    var err error

    if len(m.iterators) == 0 {
      // Start new data encoding
      bytes, err = m.encodeData(m.datas[m.dataIndex])
    } else {
      // Continue previous encoding
      bytes, err = m.encode(m.iterators[len(m.iterators)-1])
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

func (m *OrderedEncoder) encodeData(data interface{}) ([]byte, error) {
  nextIt, err := iterator.NewSortedDataIterator(data)
  if err != nil {
    return m.jsonEncode(data)
  } else {
    m.iterators = append(m.iterators, nextIt)
    return []byte{}, nil
  }
}

func (m *OrderedEncoder) encode(it *iterator.DataIterator) ([]byte, error) {
  b := []byte{}
  isMap := it.HasNamedKeys()

  if !it.Next() {
    m.iterators = m.iterators[0:len(m.iterators)-1]
    if isMap {
      b = append(b, byte('}'))
    } else {
      b = append(b, byte(']'))
    }
    b = m.tryAppendComma(b)
    return b, nil
  }

  entry := it.Value()

  if it.IsFirst() {
    if isMap {
      b = append(b, byte('{'))
    } else {
      b = append(b, byte('['))
    }
  }

  if isMap {
    key, err := json.Marshal(entry.Name)
    if err != nil { return b, err }
    b = append(b, key...)
    b = append(b, byte(':'))
  }

  bytes, err := m.encodeData(entry.Value)
  if err != nil { return b, err }
  b = append(b, bytes...)

  return b, nil
}

func (m *OrderedEncoder) tryAppendComma(b []byte) []byte {
  if len(m.iterators) > 0 && !m.iterators[len(m.iterators)-1].IsLast() {
    b = append(b, byte(','))
  }
  return b
}


func (m *OrderedEncoder) jsonEncode(data interface{}) ([]byte, error) {
  fl, ok := data.(float64)
  var bytes []byte
  var err error
  // Hack to compensate for the fact that Go converts all undefined numbers into float64
  if ok && float64(int64(fl)) == fl {
    bytes, err = json.Marshal(int64(fl))
  } else {
    bytes, err = json.Marshal(data)
  }
  bytes = m.tryAppendComma(bytes)
  return bytes, err
}
