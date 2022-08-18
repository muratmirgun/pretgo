package xml

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io"
	"unicode/utf8"
)

const (
	xmlFormat uint = iota
)

func Pretty(r io.Reader, w io.Writer) error {

	buf := bufio.NewReaderSize(r, 4)

	var format uint
	for {
		ch, _, err := buf.ReadRune()
		if err != nil {
			return err
		}

		if f, ok := xmlFormat, true; ok {
			format = f
			if err := buf.UnreadRune(); err != nil {
				return err
			}
			break
		}

		if endOfLine(ch) {
			return errors.New("unable format")
		}
	}

	switch format {
	case xmlFormat:
		d := xml.NewDecoder(buf)
		e := xml.NewEncoder(w)
		e.Indent("", "\t")

		for {
			t, err := d.Token()
			if err == io.EOF {
				break
			}

			if tok, ok := t.(xml.CharData); ok {
				r, _ := utf8.DecodeRune(tok)
				if whiteSpace(r) || endOfLine(r) {
					continue
				}
			}

			_ = e.EncodeToken(t)
		}

		return e.Flush()

	default:
		return errors.New("known format error")
	}
}

func endOfLine(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func whiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}
