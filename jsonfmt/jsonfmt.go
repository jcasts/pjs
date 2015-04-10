package jsonfmt

import (
  "io"
  "strings"
)


type Processor interface {
  Handle(token Token) string
}


type Formatter struct {
  processors []Processor
}

func NewFormatter(processors ...Processor) *Formatter {
  return &Formatter{processors}
}

func (f *Formatter) Process(r io.Reader, w io.Writer) error {
  var err error
  scanner := NewScanner(r)
  for scanner.Next() {
    token := scanner.Token()
    err = scanner.Error()
    if err != nil { return err }

    _, err = w.Write([]byte(f.Handle(*token)))
    if err != nil { return err }
  }
  err = scanner.Error()
  return err
}

func (f *Formatter) Handle(token Token) string {
  for _, processor := range f.processors {
    token.Value = processor.Handle(token)
  }
  return token.Value
}


type consoleColorizer struct {
  stringColor string
  numberColor string
  boolColor string
  nullColor string
}

func NewConsoleColorProcessor() Processor {
  return &consoleColorizer{
    stringColor: "0;36",
    numberColor: "0;33",
    boolColor: "1;35",
    nullColor: "1;31",
  }
}

func (p *consoleColorizer) Handle(t Token) string {
  switch t.Type {
  case StringLiteralToken:
    return p.wrapColor(t.Value, p.stringColor)
  case IntegerLiteralToken, FloatLiteralToken, ScientificLiteralToken:
    return p.wrapColor(t.Value, p.numberColor)
  case BooleanLiteralToken:
    return p.wrapColor(t.Value, p.boolColor)
  case NullLiteralToken:
    return p.wrapColor(t.Value, p.nullColor)
  default:
    return t.Value
  }
}

func (p *consoleColorizer) wrapColor(value, color string) string {
  return "\033["+color+"m"+value+"\033[0m"
}


type indentProcessor struct {
  prefix string
  indent string
}

func NewIndentProcessor(prefix, indent string) Processor {
  return &indentProcessor{prefix, indent}
}

func (p *indentProcessor) Handle(t Token) string {
  switch t.Type {
  case MapKeyToken, StringLiteralToken, IntegerLiteralToken, FloatLiteralToken,
        BooleanLiteralToken, NullLiteralToken, EmptyMapToken, EmptyArrayToken:
    return t.Value
  case MapStartToken, ArrayStartToken:
    prefix := ""
    if t.Depth == 0 { prefix = p.prefix }
    return prefix + t.Value + "\n" + p.indentString("", t.Depth+1)
  case MapEndToken, ArrayEndToken:
    return "\n" + p.indentString(t.Value, t.Depth)
  case MapColonToken:
    return t.Value + " "
  case ValueSeparatorToken:
    return t.Value + "\n" + p.indentString("", t.Depth)
  case StartNewJsonToken:
    return "\n\n"
  default:
    return t.Value
  }
}

func (p *indentProcessor) indentString(s string, depth int) string {
  return p.prefix + strings.Repeat(p.indent, depth) + s
}
