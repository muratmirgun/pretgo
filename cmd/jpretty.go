package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	unknownFormat uint = iota
	jsonFormat
)

func init() {
	rootCmd.AddCommand(xprettyCmd)
}

var xprettyCmd = &cobra.Command{
	Use:   "pjson",
	Short: "Pretty-print the serialised output",
	Run: func(cmd *cobra.Command, args []string) {
		if err := prettyjson(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	},
}

func prettyjson(r io.Reader, w io.Writer) error {

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

	default:
		return errors.New("known format error, please file a bug")
	}

	return nil
}

func EndOfLine(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func Whitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}
