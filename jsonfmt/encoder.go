package jsonfmt

import (
  "encoding/json"
  "io"
  "../iterator"
)


type Encoder struct {
  datas []iterator.Value
  dataIndex int
  buffer []byte
  iterators []iterator.Iterator
}

func NewEncoder() *Encoder {
  return &Encoder{[]iterator.Value{}, 0, []byte{}, []iterator.Iterator{}}
}

func (m *Encoder) Queue(data iterator.Value) {
  m.datas = append(m.datas, data)
}

func (m *Encoder) Read(p []byte) (n int, err error) {
  bytesToRead := len(p)
  bytesRead := len(m.buffer)

  for bytesRead < bytesToRead && len(m.datas) > 0 {
    var bytes []byte
    var err error

    if len(m.iterators) == 0 {
      // Start new data encoding
      bytes, err = m.encodeData(m.datas[0])
    } else {
      // Continue previous encoding
      bytes, err = m.encode(m.iterators[len(m.iterators)-1])
    }
    if err != nil { return bytesRead, err }

    m.buffer = append(m.buffer, bytes...)
    bytesRead = len(m.buffer)
    if len(m.iterators) == 0 {
      m.datas = m.datas[1:len(m.datas)]
      if len(m.datas) > 0 {
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

func (m *Encoder) encodeData(data iterator.Value) ([]byte, error) {
  if data.HasIterator() {
    m.iterators = append(m.iterators, data.Iterator())
    return []byte{}, nil
  } else {
    return m.jsonEncode(data.Interface())
  }
}

func (m *Encoder) encode(it iterator.Iterator) ([]byte, error) {
  b := []byte{}
  isMap := it.HasNamedKeys()

  if !it.Next() {
    if it.IsFirst() {
      // Empty iterator
      if isMap {
        b = append(b, byte('{'))
      } else {
        b = append(b, byte('['))
      }
    }
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
    key, err := json.Marshal(entry.Name())
    if err != nil { return b, err }
    b = append(b, key...)
    b = append(b, byte(':'))
  }

  bytes, err := m.encodeData(entry)
  if err != nil { return b, err }
  b = append(b, bytes...)

  return b, nil
}

func (m *Encoder) tryAppendComma(b []byte) []byte {
  if len(m.iterators) > 0 && !m.iterators[len(m.iterators)-1].IsLast() {
    b = append(b, byte(','))
  }
  return b
}


func (m *Encoder) jsonEncode(data interface{}) ([]byte, error) {
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
