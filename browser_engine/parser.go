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

/**
 * Accept the next byte if it's any of the provided check bytes.
 * TODO: is this needed?
 */
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

/**
 * Consumes bytes until provided test function returns true.
 */
func (p *Parser) acceptByteGivenTest(test func(val byte) bool) (byte, bool) {
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

/**
 * Consumes all the bytes until provided test function returns true.
 */
func (p *Parser) acceptBytesUntilTest(test func(val byte) bool) (string, bool) {
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
    return "", false
  }
  return sb.String(), true
}

/**
 * Peeks ahead and consumes a string if it's found.
 */
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


/**
 * Looks ahead without consuming.
 */
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

// Consumes nothingness: carriage returns, newlines and spaces.
func (p *Parser) consumeWhitespace() {
  state := true
  for state == true {
    _, state = p.accept('\r', '\n', ' ')
  }
}

func isAlphanumericOrPunctuation(check byte) bool {
  return unicode.IsLetter(rune(check)) || unicode.IsNumber(rune(check)) || unicode.IsPunct(rune(check))
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
  p.accept('<')
  p.accept('!')
  p.acceptString("DOCTYPE html")
  p.accept('>')
  p.node()
  return true
}

/**
 * Parses a single node.
 */
func (p *Parser) node() bool {
  p.consumeWhitespace()

  var openTagName, openState = p.openTag()
  if !openState {
    return false
  }
  fmt.Println("Starting node! ", openTagName)

  var continueRecursion = true
  for continueRecursion {
    // Loop while we find either nodes or text.
    if p.node() {
      continue
    } else {
      // check for comments?

      // consume anything else.
      var val, valid = p.acceptBytesUntilTest(isAlphanumericOrPunctuation)
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
      "Tag name mismatch! Expected: ", openTagName, " Got: ", closeTagName,
    )
  }
  fmt.Println("Ending node! ", openTagName)
  return true
}

func (p *Parser) openTag() (string, bool) {
  if p.tagCloseSequence() {
    return "", false
  }

  var _, valid = p.accept('<')
  if !valid {
    return "", false
  }

  var tagName, tagCompleted = p.tagName()

  if tagCompleted {
    return tagName, true
  }

  // Tag isn't over yet, cover any attributes we can find.
  for p.attribute() == true {
    p.consumeWhitespace()

    // Quit the loop when we find the end of the tag.
    if p.tagEnd() {
      p.accept('>')
      return tagName, true
    }
  }

  return tagName, true
}

func (p *Parser) closeTag() (string, bool) {
  var isClose = p.tagCloseSequence()
  if !isClose {
    return "", false
  }

  p.acceptString("</")

  return p.tagName()
}

func (p *Parser) attribute() bool {
  var attributeName, nameState = p.acceptBytesUntilTest(func(val byte) bool {
    return val != '='
  })

  if !nameState {
    return false
  }

  p.accept('=')
  
  p.consumeWhitespace()
  
  if !p.quote() {
    return false
  }

  var attributeValue, _ = p.acceptBytesUntilTest(func(val byte) bool {
    return val != '"' && val != '\''
  })

  log.Println("Attribute: ", "(", attributeName, "=", attributeValue, ")")

  if !p.quote() {
    return false
  }

  return true
}

func (p *Parser) quote() bool {
  var _, quote = p.accept('"', '\'')
  // invalid attribute value
  if !quote {
    return false
  }
  return true
}

func (p *Parser) tagEnd() bool {
  var isEnd = p.assertNext('>')
  if isEnd {
    return true
  }
  return false
}

func (p *Parser) tagCloseSequence() bool {
  var isClose = p.assertNext('<', '/')
  if isClose {
    return true
  }
  return false
}

func (p *Parser) tagName() (string, bool) {
  var tagName, _ = p.acceptBytesUntilTest(func(val byte) bool {
    return val != '>' && val != ' '
  })

  // Now we've determined the tag name, consume any remaining whitespace.
  p.consumeWhitespace()

  // Exit early if we've reached the end of the tag.
  if p.tagEnd() {
    p.accept('>')
    return tagName, true
  }
  
  return tagName, false
}


func (p *Parser) commentStart() bool {
  var isComment = p.assertNext('<', '!', '-', '-')
  return isComment
}

func (p *Parser) commentEnd() bool {
  var isCommentEnd = p.assertNext('-', '-', '>')
  return isCommentEnd
}