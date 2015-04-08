package jsonfmt


type OrderedEncoder struct {
  datas []interface{}
}

func NewOrderedEncoder(datas ...interface{}) *OrderedEncoder {
  return &OrderedEncoder{datas}
}

func (m* OrderedEncoder) Read(p []byte) (n int, err error) {
  return 0, nil
}