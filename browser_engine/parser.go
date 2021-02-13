package main

import (
  "bufio"
  "log"
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

func (p *Parser) accept(check byte) bool {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  if check == next[0] {
    p.reader.ReadByte()
    return true
  }
  return false
}

func (p *Parser) expect(check byte) {
  
}

func (p *Parser) Parse() {
  r := p.reader
  b, err := r.Peek(8)
  if err != nil {
      fmt.Println(err)
  }
  fmt.Printf("%q\n", b)
}


