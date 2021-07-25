package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

func init() {
	rootCmd.AddCommand(prettyCmd)
}

var prettyCmd = &cobra.Command{
	Use:   "phtml",
	Short: "Pretty-print the serialised output",
	Run: func(cmd *cobra.Command, args []string) {
		indent := "  "
		wrap, _ := strconv.Atoi(args[0])
		node, err := html.Parse(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed parsing HTML: ", err)
			os.Exit(1)
		}
		if err := Print(os.Stdout, node, indent, wrap); err != nil {
			fmt.Fprint(os.Stderr, "Failed printing HTML: ", err)
			os.Exit(1)
		}
	},
}

func Print(w io.Writer, root *html.Node, indent string, wrap int) error {
	p := printer{
		w:         w,
		indentStr: indent,
		wrapWidth: wrap,
		lineStart: true,
	}
	if err := p.doc(root); err != nil {
		return err
	}
	return p.werr
}

// tagSet holds a set of HTML tag names.
type tagSet map[string]struct{}

func newTagSet(tags []string) tagSet {
	ts := make(tagSet)
	for _, t := range tags {
		ts[t] = struct{}{}
	}
	return ts
}

func (ts tagSet) has(n *html.Node) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}
	_, ok := ts[n.Data]
	return ok
}

var voidTags = newTagSet(strings.Fields("area base br col embed hr img input link meta param source track wbr"))

var inlineTags = newTagSet(strings.Fields("a amp-img b code em i img picture span s source strong"))

var omitCloseTags = newTagSet(strings.Fields("li"))

var literalTags = newTagSet(strings.Fields("noscript script style"))

var keepSpaceTags = newTagSet(strings.Fields("pre"))

type printer struct {
	w         io.Writer
	werr      error
	indentStr string
	wrapWidth int

	level          int
	literalDepth   int
	keepSpaceDepth int
	lineStart      bool
	lineWidth      int
}

func (p *printer) inLiteral() bool {
	return p.literalDepth > 0
}
func (p *printer) inKeepSpace() bool {
	return p.keepSpaceDepth > 0
}

func (p *printer) doc(n *html.Node) error {
	if n.Type != html.DocumentNode {
		return fmt.Errorf("root node has non-document type %v", n.Type)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.DoctypeNode:
			p.write("<!DOCTYPE " + c.Data + ">")
			p.endl()
		case html.ElementNode:
			if err := p.element(c); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled doc child %q with type %v", c.Data, c.Type)
		}
	}
	return nil
}

func (p *printer) element(n *html.Node) error {
	tag := n.Data
	if n.Type != html.ElementNode {
		return fmt.Errorf("got non-element node %q of type %v", tag, n.Type)
	}

	inline := inlineTags.has(n)
	if forceInline := p.openTag(n); forceInline {
		inline = true
	}

	literal := literalTags.has(n)
	if literal {
		p.literalDepth++
	}
	keepSpace := keepSpaceTags.has(n)
	if keepSpace {
		p.keepSpaceDepth++
	}

	omitClose := omitCloseTags.has(n)
	if !inline && !omitClose {
		p.endl()
	}

	if voidTags.has(n) {
		if literal || keepSpace {
			panic(fmt.Sprintf("<%s> is both literal/keep-space and void", n.Data))
		}
		return nil
	}

	// Indent if needed and print the children.
	if !inline {
		p.level++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.ElementNode:
			if err := p.element(c); err != nil {
				return err
			}
		case html.TextNode:
			if err := p.text(c); err != nil {
				return err
			}
		case html.CommentNode:
			continue
		default:
			return fmt.Errorf("unexpected node %q of type %d", c.Data, c.Type)
		}
	}
	if !inline {
		p.level--
		p.endl()
	}

	if !omitClose {
		p.maybeIndent()
		p.write(closeTag(n))
	}
	if literal {
		p.literalDepth--
	}
	if keepSpace {
		p.keepSpaceDepth--
	}
	if !inline {
		p.endl()
	}
	return nil
}

func (p *printer) text(n *html.Node) error {
	if n.Type != html.TextNode {
		panic(fmt.Sprintf("Got non-text node %q (type %v)", n.Data, n.Type))
	}
	if len(n.Data) == 0 {
		return nil
	}

	if p.inLiteral() {
		p.write(n.Data)
		return nil
	}

	s := n.Data
	s = escapeText(s)

	if p.inKeepSpace() {
		p.write(s)
		return nil
	}

	s = collapseText(s, n)
	if s == "" {
		return nil
	}

	p.maybeIndent()

	startSpace := s[0] == ' '
	endSpace := s[len(s)-1] == ' '

	wrapStart := 0
	if (inlineTags.has(n.PrevSibling) || inlineTags.has(n.Parent)) && !startSpace {
		wrapStart = 1
	}

	words := strings.Fields(strings.TrimSpace(s))
	for i, w := range words {
		if (i == 0 && startSpace) || i != 0 {
			w = " " + w
		}
		if i == len(words)-1 && endSpace && w != " " {
			w = w + " "
		}

		if i < wrapStart {
			p.write(w)
		} else {
			p.wrap(w, "")
		}
	}
	return nil
}

func (p *printer) maybeIndent() {
	if p.inLiteral() || p.inKeepSpace() || !p.lineStart {
		return
	}
	s := strings.Repeat(p.indentStr, p.level)
	p.write(s) // updates lineStart and lineWidth
}

func (p *printer) wrap(s, extra string) {
	if !p.inLiteral() && !p.inKeepSpace() &&
		p.wrapWidth > 0 && p.lineWidth+len(s) > p.wrapWidth {
		p.endl()
		p.maybeIndent()
		s = extra + strings.TrimLeft(s, " ")
	}
	p.write(s)
}

func (p *printer) endl() {
	if p.inLiteral() || p.inKeepSpace() {
		return
	}
	if p.lineStart {
		return
	}
	p.write("\n")
	p.lineStart = true
	p.lineWidth = 0
}

func (p *printer) write(s string) {
	if p.werr != nil {
		return
	}
	_, p.werr = io.WriteString(p.w, s)
	p.lineStart = false
	p.lineWidth += len(s)
}

func (p *printer) openTag(n *html.Node) (forceInline bool) {
	tokens := append([]string{}, "<"+n.Data)
	for _, a := range n.Attr {
		as := " " + a.Key
		if len(a.Val) > 0 {
			escaped := strings.Replace(a.Val, `"`, `&quot;`, -1)
			as += `="` + escaped + `"`
		}
		tokens = append(tokens, as)
	}
	tokens[len(tokens)-1] += ">"
	tagLen := len(strings.Join(tokens, ""))

	inline := inlineTags.has(n)
	wouldWrap := p.wrapWidth > 0 && p.lineWidth+tagLen > p.wrapWidth
	prev := n.PrevSibling
	prevTextNotSpace := prev != nil && prev.Type == html.TextNode &&
		(prev.Data == "" || !whitespace.MatchString(prev.Data[len(prev.Data)-1:]))
	startSpaceMatters := inlineTags.has(prev) || inlineTags.has(n.Parent) || prevTextNotSpace
	if !inline || (wouldWrap && !startSpaceMatters) {
		p.endl()
	}

	startedLine := p.lineStart
	p.maybeIndent()

	if !literalTags.has(n) && !p.inLiteral() &&
		!keepSpaceTags.has(n) && !p.inKeepSpace() {
		childLen := -1
		if n.FirstChild == nil {
			childLen = 0
		} else if hasSingleChild(n) && n.FirstChild.Type == html.TextNode {
			childLen = len(collapseText(escapeText(n.FirstChild.Data), n.FirstChild))
		}
		if childLen >= 0 && (p.lineWidth+tagLen+childLen+len(closeTag(n)) < p.wrapWidth || p.wrapWidth <= 0) {
			forceInline = true
		}
	}

	var unwrapTokens int
	var wrapIndent string
	if startedLine {
		wrapIndent = strings.Repeat(p.indentStr, 2)
		unwrapTokens = 1
		if len(tokens[0]) < len(wrapIndent) {
			unwrapTokens = 2
		}
	} else if (inline || forceInline) && startSpaceMatters {
		unwrapTokens = 1
	}
	for i, t := range tokens {
		if i < unwrapTokens {
			p.write(t)
		} else {
			p.wrap(t, wrapIndent)
		}
	}

	return forceInline
}

func hasSingleChild(n *html.Node) bool {
	return n.FirstChild != nil && n.FirstChild == n.LastChild
}

func closeTag(n *html.Node) string {
	if n.Type != html.ElementNode || voidTags.has(n) || omitCloseTags.has(n) {
		return ""
	}
	return "</" + n.Data + ">"
}

func escapeText(s string) string {
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	return s
}

var whitespace *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]+`)

func collapseText(s string, n *html.Node) string {
	s = whitespace.ReplaceAllString(s, " ")

	if !inlineTags.has(n.Parent) {
		if !inlineTags.has(n.PrevSibling) {
			s = strings.TrimLeft(s, " ")
		}
		if !inlineTags.has(n.NextSibling) {
			s = strings.TrimRight(s, " ")
		}
	}

	return s
}
