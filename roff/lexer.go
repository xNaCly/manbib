package roff

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"
)

type Lexer struct {
	scan    *bufio.Scanner
	isAtEnd bool
	curLine []rune
	curChar rune
	linePos int
	line    int
}

func New() *Lexer {
	return &Lexer{
		scan:    nil,
		isAtEnd: false,
		curLine: nil,
		curChar: 0,
		linePos: 0,
		line:    0,
	}
}

func (l *Lexer) NewInput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	l.scan = bufio.NewScanner(gz)
	l.scan.Scan()

	l.curLine = []rune(l.scan.Text())
	if len(l.curLine) == 0 {
		return errors.New("given input is empty")
	}
	l.curChar = l.curLine[l.linePos]
	return nil
}

func (l *Lexer) advance(n int) {
	if l.linePos+n >= len(l.curLine) {
		l.curChar = '\n'
		l.linePos++
		return
	}

	l.linePos += n
	l.curChar = l.curLine[l.linePos]
}

func (l *Lexer) prev() rune {
	return l.curLine[l.linePos-1]
}

func (l *Lexer) advanceLine() {
	ok := l.scan.Scan()

	if l.scan.Err() != nil || !ok {
		l.isAtEnd = true
		return
	}

	l.curLine = []rune(l.scan.Text())
	l.line++
	l.linePos = 0
	for len(l.curLine) == 0 && ok {
		ok = l.scan.Scan()
		l.curLine = []rune(l.scan.Text())
		l.line++
	}
	if !ok {
		l.isAtEnd = true
		return
	}
	l.curChar = l.curLine[l.linePos]
}

func (l *Lexer) peekEquals(c rune) bool {
	if l.linePos >= len(l.curLine) {
		return l.curLine[l.linePos+1] == c
	}
	return false
}

func (l *Lexer) Start() []Token {
	res := make([]Token, 0)
	for !l.isAtEnd {
		kind := 0
		switch l.curChar {
		case '"':
			if l.prevChar == '\\' {
				l.advanceLine()
				continue
			}
		case '.':
			if l.prevChar == '\\' {
				res = append(res, Token{
					Kind:    TEXT,
					Content: "\\.",
					Pos:     l.linePos,
					Line:    l.line,
				})
				l.advance(1)
				continue
			}
			l.advance(1)
		case '\n':
			res = append(res, Token{
				Kind: NEWLINE,
				Pos:  l.linePos,
				Line: l.line,
			})
			l.advanceLine()
			continue
		case 'B':
			if l.peekEquals('I') {
				kind = BOLDITALIC
				l.advance(3)
			} else if l.peekEquals('R') {
				kind = BOLDROMAN
				l.advance(3)
			} else {
				kind = BOLD
			}
		case 'I':
			if l.peekEquals('B') {
				kind = ITALICBOLD
				l.advance(3)
			} else if l.peekEquals('R') {
				kind = ITALICROMAN
				l.advance(3)
			} else if l.peekEquals('P') {
				kind = INDENTPARAGRAPH
				l.advance(3)
			} else {
				kind = ITALIC
			}
		case 'E':
			if l.peekEquals('X') {
				kind = EXAMPLESTART
				l.advance(3)
			} else if l.peekEquals('E') {
				kind = EXAMPLEEND
				l.advance(3)
			}
		case 'T':
			if l.peekEquals('H') {
				kind = TITLEHEADING
				l.advance(3)
			} else if l.peekEquals('P') {
				kind = TAGGEDPARAGRAPH
				l.advance(3)
			}
		default:
			b := strings.Builder{}
			for {
				if (l.curChar == '.' && l.prevChar != '\\') || l.curChar == '\n' {
					break
				}
				b.WriteRune(l.curChar)
				l.advance(1)
			}

			res = append(res, Token{
				Pos:     l.linePos - b.Len(),
				Kind:    TEXT,
				Line:    l.line,
				Content: b.String(),
			})

			continue
		}

		l.advance(1)
		if kind == 0 {
			continue
		}
		res = append(res, Token{
			Kind: kind,
			Pos:  l.linePos,
			Line: l.line,
		})
	}
	return res
}
