package main

import (
  "bufio"
  "strings"
  "fmt"
)

type Parser struct {
  input string
  position int
  reader *bufio.Reader
}

func NewParser(input string, position int) *Parser {
  sr := strings.NewReader(input)
  var reader *bufio.Reader = bufio.NewReader(sr)
  return &Parser{
    input,
    position,
    reader,
  }
}

func (p *Parser) Parse() {
  r := p.reader
  b, err := r.Peek(3)
  if err != nil {
      fmt.Println(err)
  }
  fmt.Printf("%q\n", b)
}


