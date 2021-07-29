package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/muratmirgun/pretgo/internal/json"
	"github.com/muratmirgun/pretgo/internal/xml"
	"github.com/yosssi/gohtml"
	"log"
	"os"
)

func main() {

	wordPtr := flag.String("format", "nil", "format style")
	flag.Parse()

	switch *wordPtr {
	case "json":
		if err := json.Pretty(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	case "xml":
		if err := xml.Pretty(os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	case "html":
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(gohtml.Format(scanner.Text()))
		}
	}
}
