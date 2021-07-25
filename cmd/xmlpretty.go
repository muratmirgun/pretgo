package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

const (
	xmlFormat uint = iota
)

func init() {
	rootCmd.AddCommand(jprettyCmd)
}

var jprettyCmd = &cobra.Command{
	Use:   "pxml",
	Short: "Pretty-print the serialised output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := prettyxml(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func prettyxml(r io.Reader, w io.Writer) error {

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
			return errors.New("unable to recognize this format")
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

			e.EncodeToken(t)
		}

		return e.Flush()

	default:
		return errors.New("known format error, please file a bug")
	}

	return nil
}

func endOfLine(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func whiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}
