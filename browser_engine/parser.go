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
  var reader = bufio.NewReader(sr)
  return &Parser{
    input,
    position,
    reader,
  }
}

/**
 * Accept the next byte if it's any of the provided check bytes.
 */
func (p *Parser) accept(check ...byte) bool {
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
      return true
    }
  }
  return false
}

/**
 * Asserts that the next len(check) bytes match the check string.
 */
func (p *Parser) assertString(check string) bool {
  var next, err = p.reader.Peek(len(check))
  if err != nil {
    log.Fatal(err)
  }
  return string(next) == check
}

/**
 * Peeks ahead and consumes a string if it's found.
 */
func (p *Parser) acceptString(check string) bool {
  if p.assertString(check) {
    var readUntil = check[len(check) - 1]
    p.reader.ReadBytes(readUntil)
    return true
  }
  return false
}

/**
 * Consumes bytes until provided test function returns true.
 */
func (p *Parser) acceptByteGivenTest(test func(val byte) bool) (byte, bool) {
  var next, err = p.reader.Peek(1)
  if err != nil {
    log.Fatal(err)
  }
  valid := test(next[0])
  if valid {
    _, err = p.reader.ReadByte()
    if err != nil {
      log.Fatal(err)
    }
    return next[0], true
  }
  return 0, false
}

/**
 * Consumes all the bytes until provided test function returns true.
 */
func (p *Parser) acceptBytesUntilTest(test func(val byte) bool) string {
  var sb strings.Builder
  var state = true
  var val byte
  for state {
    val, state = p.acceptByteGivenTest(test)
    if state {
      sb.WriteByte(val)
    }
  }
  if len(sb.String()) == 0 {
    return ""
  }
  return sb.String()
}

// Consumes nothingness: carriage returns, newlines and spaces.
func (p *Parser) consumeWhitespace() {
  state := true
  for state == true {
    state = p.accept('\r', '\n', ' ')
  }
}

func isAlphanumericOrPunctuation(check byte) bool {
  return unicode.IsLetter(rune(check)) ||
    unicode.IsNumber(rune(check)) ||
    unicode.IsPunct(rune(check)) ||
    check == ' ' ||
    check == '\r' ||
    check == '\n'
}

func isEquals(check byte) bool {
  return check != '='
}

func isEqualsOrCloseBracket(check byte) bool {
  return check != '=' && check != '>'
}

func isQuote(check byte) bool {
  return check != '"' && check != '\''
}

/* Actual parsing starts here */

func (p *Parser) Parse() bool {
  if p.document() {
    return true
  }
  log.Fatal("Problem parsing input")
  return false
}

func (p *Parser) document() bool {
  p.acceptString("<!DOCTYPE html>")
  p.node()
  return true
}

/**
 * Parses a single node.
 */
func (p *Parser) node() bool {
  p.consumeWhitespace()

  if !p.openTag() {
    return false
  }

  for p.node() || p.text() {}

  if !p.closeTag() {
    return false
  }

  return true
}

func (p *Parser) text() bool {
   var val = p.acceptBytesUntilTest(isAlphanumericOrPunctuation)
   if len(val) > 0 {
     fmt.Println("Consumed string: ", val)
     return true
   }
   return false
}

func (p *Parser) openTag() bool {
  // if it's a close tag, bail out.
  if p.assertString("</") {
    return false
  }

  if !p.accept('<') {
    return false
  }

  var tagName = p.tagName()
  fmt.Println("Open tag: ", tagName)

  // Exit early if we've reached the end of the tag.
  if p.accept('>') {
    return true
  }

  // Tag isn't over yet, cover any attributes we can find.
  for p.attribute() == true {
    p.consumeWhitespace()
    // Quit the loop when we find the end of the tag.
    if p.accept('>') {
      return true
    }
  }

  return false
}

func (p *Parser) closeTag() bool {
  if !p.acceptString("</") {
    return false
  }
  
  var tagName = p.tagName()

  fmt.Println("Close tag: ", tagName)

  return p.accept('>')
}

func (p *Parser) tagName() string {
  var tagName = p.acceptBytesUntilTest(func(val byte) bool {
    return !(val == '>' || val == ' ')
  })

  // Now we've determined the tag name, consume any remaining whitespace.
  p.consumeWhitespace()
  return tagName
}

func (p *Parser) attribute() bool {
  var attributeName = p.acceptBytesUntilTest(isEqualsOrCloseBracket)


  // We have an attribute without a value
  if !p.accept('=') {
    log.Println("Attribute: ", "(", attributeName, ")")
    return true
  }
  
  p.consumeWhitespace()
  
  if !p.accept('"', '\'') {
    return false
  }

  var attributeValue = p.acceptBytesUntilTest(isQuote)

  if !p.accept('"', '\'') {
    return false
  }

  log.Println("Attribute: ", "(", attributeName, "=", attributeValue, ")")
  return true
}

//
//func (p *Parser) commentStart() bool {
// var start = p.acceptString("<!--")
// return isComment
//}
//
//func (p *Parser) commentEnd() bool {
// var _, isCommentEnd = p.acceptString("-->")
// return isCommentEnd
//}