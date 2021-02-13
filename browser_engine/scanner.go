package main

import (
	"bufio"
	"strings"
)

type Token byte

const (
	TAG_NAME Token = iota //
	OPEN_TAG
	CLOSE_TAG
	EOF
	WHITESPACE
)

type Scanner struct {
	reader *bufio.Reader
}

func NewScanner(input string) *Scanner {
	sr := strings.NewReader(input)
	return &Scanner{reader: bufio.NewReader(sr)}
}

func (s* Scanner) peek() (byte, error) {
	var ch, err = s.reader.Peek(1)
	if err != nil {
		return 0, err
	}
	return ch[0], nil
}

func (s* Scanner) read() (byte, error) {
	var ch, err = s.reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return ch, nil
}