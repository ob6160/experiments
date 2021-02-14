package main

import (
  "bufio"
  "log"
  "strings"
  "unicode"
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

func (p *Parser) accept(check ...byte) (byte, bool) {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  for _, checkByte := range check {
    if checkByte == next[0] {
      _, err = p.reader.ReadByte()
      if err != nil {
        log.Fatal(err)
      }
      return checkByte, true
    }
  }
  return next[0], false
}

func (p *Parser) acceptString(check string) (string, bool) {
  var next, err = p.reader.Peek(len(check))
  if err != nil {
    log.Fatal(err)
  }
  if string(next) == check {
    var readUntil = next[len(next) - 1]
    p.reader.ReadBytes(readUntil)
    return check, true
  }
  return string(next), false
}

func (p *Parser) acceptUntil(delim byte) (string, bool) {
  var accepted, err = p.reader.ReadBytes(delim)
  if err != nil {
    // delimiter never reached :(
    log.Fatal(err)
    return string(accepted), false
  }
  return string(accepted), true
}

func (p *Parser) acceptWhitespace() {
  // Consume nothingness
  state := true
  for state == true {
    _, state = p.accept('\r', '\n', ' ')
  }
}

func (p *Parser) acceptAlphanumeric() (string, bool) {

}

func (p *Parser) expect(check byte) bool {
  var val, state = p.accept(check)
  if state {
    return true
  }
  log.Fatal("Syntax error, Expected: ", string(check), " Got: ", string(val))
  return false
}

func (p *Parser) expectString(check string) bool {
  var val, state = p.acceptString(check)
  if state {
    return true
  }
  log.Fatal("Syntax error! Expected: ", check, " Got: ", val)
  return false
}

//func (p *Parser) expectAlphanumeric() (string, bool) {
// var val, state = p.acc
//}


func (p *Parser) Parse() bool {
  if p.document() {
    return true
  }
  log.Fatal("Problem parsing input")
  return false
}

func (p *Parser) document() bool {
  p.expect('<')
  p.expect('!')
  p.expectString("DOCTYPE html")
  p.expect('>')
  p.node()
  return true
}

func (p *Parser) node() bool {
  p.acceptWhitespace()
  var openTagName, state = p.openTag()
  if !state {
    return false
  }
  for p.node() {
    
  }
  var closeTagName, _ = p.closeTag()
  
  if openTagName != closeTagName {
    log.Fatal(
      "Tag name mismatch! Expected: ", openTagName[:len(openTagName)-1], " Got: ", closeTagName[:len(closeTagName)-1],
    )
  }
  return false
}

func (p *Parser) openTag() (string, bool) {
  if !p.expect('<') {
    return "", false
  }
  var tagName, state = p.acceptUntil('>')
  return tagName, state
}

func (p *Parser) closeTag() (string, bool) {
  p.expect('<')
  p.expect('/')
  var tagName, state = p.acceptUntil('>')
  return tagName, state
}





