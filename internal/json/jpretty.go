package json

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

const (
	jsonFormat uint = iota
)

func Pretty(r io.Reader, w io.Writer) error {

	buf := bufio.NewReaderSize(r, 4)

	var format uint
	for {
		ch, _, err := buf.ReadRune()
		if err != nil {
			return err
		}

		if f, ok := jsonFormat, true; ok {
			format = f
			if err := buf.UnreadRune(); err != nil {
				return err
			}
			break
		}

		if EndOfln(ch) {
			return errors.New("unable format")
		}
	}

	switch format {
	case jsonFormat:
		b, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}

		var out bytes.Buffer
		if err := json.Indent(&out, b, "", "\t"); err != nil {
			return err
		}

		if _, err := out.WriteTo(w); err != nil {
			return err
		}

	default:
		return errors.New("format error")
	}

	return nil
}

func EndOfln(ch rune) bool {
	return ch == '\n' || ch == '\r'
}
