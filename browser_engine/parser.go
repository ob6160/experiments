package main

import (
  "bufio"
  "fmt"
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

func (p *Parser) acceptTest(test func(val byte) bool) (byte, bool) {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  value := next[0]
  valid := test(value)
  if valid {
    _, err = p.reader.ReadByte()
    if err != nil {
      log.Fatal(err)
    }
    return value, true
  }
  return 0, false
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

func (p *Parser) assertNext(chars ...byte) bool {
  var next, err = p.reader.Peek(len(chars))
  if err != nil {
    return false
  }
  for i, v := range chars {
    if v != next[i] {
      return false
    }
  }
  return true
}

func (p *Parser) consumeWhitespace() {
  // Consume nothingness
  state := true
  for state == true {
    _, state = p.accept('\r', '\n', ' ')
  }
}

func isAlphanumericOrPunctuation(check byte) bool {
  return unicode.IsLetter(rune(check)) || unicode.IsNumber(rune(check)) || unicode.IsPunct(rune(check))
}

func (p *Parser) consumeGivenTest(test func(val byte) bool) (string, bool) {
  var sb strings.Builder
  var state = true
  var val byte
  for state {
    val, state = p.acceptTest(test)
    if state {
      sb.WriteByte(val)
    }
  }
  if len(sb.String()) == 0 {
    return "", false
  }
  return sb.String(), true
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
  p.consumeWhitespace()
  var openTagName, openState = p.openTag()
  if !openState {
    return false
  }
  fmt.Println("Starting node! ", openTagName[:len(openTagName)-1])

  var continueRecursion = true
  for continueRecursion {
    // Loop while we find either nodes or text.
    if p.node() {
      continue
    } else {
      var val, valid = p.consumeGivenTest(isAlphanumericOrPunctuation)
      if valid {
        fmt.Println("Consumed string: ", val)
        continue
      }
    }
    // Reached a point with no text or nodes, exit.
    continueRecursion = false
  }

  var closeTagName, closeState = p.closeTag()
  if !closeState {
    return false
  }

  if openTagName != closeTagName {
    log.Fatal(
      "Tag name mismatch! Expected: ", openTagName[:len(openTagName)-1], " Got: ", closeTagName[:len(closeTagName)-1],
    )
  }
  fmt.Println("Ending node! ", openTagName[:len(closeTagName)-1])
  return true
}

func (p *Parser) openTag() (string, bool) {
  var isClose = p.assertNext('<', '/')

  if isClose {
    return "", false
  }

  var _, valid = p.accept('<')
  if !valid {
    return "", false
  }

  var tagName, state = p.acceptUntil('>')
  return tagName, state
}

func (p *Parser) closeTag() (string, bool) {
  var _, valid = p.acceptString("</")

  if !valid {
    return "", false
  }

  var tagName, state = p.acceptUntil('>')
  return tagName, state
}





