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

func (p *Parser) acceptString(check string) bool {
  var next, err = p.reader.Peek(len(check) - 1)
  if err != nil {
    log.Fatal(err)
  }
  if string(next) == check {
    var readUntil byte = next[len(next) - 1]
    p.reader.ReadBytes(readUntil)
    return true
  }
  return false
}

func (p *Parser) expect(check byte) bool {
  if p.accept(check) {
    return true
  }
  log.Fatal("Syntax error")
  return false
}

func (p *Parser) expectString(check string) bool {
  if p.accept(check[0]) {
    return true
  }
  log.Fatal("Syntax error")
  return false
}


func (p *Parser) Parse() bool {
  if p.document() {
    return true
  }
  log.Fatal("Problem parsing input")
  return false
}

//func (p *Parser) isTagname(tag string) bool {
//  return true
//}
//
//func (p *Parser) isWhitespace(char byte) bool {
//  return true
//}

func (p *Parser) document() bool {
  p.expect('<')
  p.expect('!')
  p.expectString("DOCTYPE html", '>')

}

func (p *Parser) node() bool {
  
}

func (p *Parser) openTag() bool {
  
}

func (p *Parser) closeTag() bool {
  
}





