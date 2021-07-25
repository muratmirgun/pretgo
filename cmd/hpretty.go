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
		wrp, _ := strconv.Atoi(args[0])
		node, err := html.Parse(os.Stdin)
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed parsing HTML: ", err)
			os.Exit(1)
		}
		if err := Print(os.Stdout, node, indent, wrp); err != nil {
			fmt.Fprint(os.Stderr, "Failed printing HTML: ", err)
			os.Exit(1)
		}
	},
}

func Print(w io.Writer, root *html.Node, indent string, wrp int) error {
	p := printer{
		w:      w,
		iStr:   indent,
		wWidth: wrp,
		lStart: true,
	}
	if err := p.doc(root); err != nil {
		return err
	}
	return p.werr
}

type tagSet map[string]struct{}

func nwtagSet(tags []string) tagSet {
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

var tagvoid = nwtagSet(strings.Fields("area base br col embed hr img input link meta param source track wbr"))

var inlnTags = nwtagSet(strings.Fields("a amp-img b code em i img picture span s source strong"))

var clstags = nwtagSet(strings.Fields("li"))

var ltrltags = nwtagSet(strings.Fields("noscript script style"))

var spacekeepTags = nwtagSet(strings.Fields("pre"))

type printer struct {
	w      io.Writer
	werr   error
	iStr   string
	wWidth int

	lvl       int
	ltrlkeep  int
	spckeeper int
	lStart    bool
	lWidth    int
}

func (p *printer) inLiteral() bool {
	return p.ltrlkeep > 0
}
func (p *printer) inspcKeep() bool {
	return p.spckeeper > 0
}

func (p *printer) doc(n *html.Node) error {
	if n.Type != html.DocumentNode {
		return fmt.Errorf("root node has non-document type %v", n.Type)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.DoctypeNode:
			p.write("<!DOCTYPE " + c.Data + ">")
			p.endline()
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

	inln := inlnTags.has(n)
	if forceinln := p.tagOpener(n); forceinln {
		inln = true
	}

	literal := ltrltags.has(n)
	if literal {
		p.ltrlkeep++
	}
	spcKeep := spacekeepTags.has(n)
	if spcKeep {
		p.spckeeper++
	}

	omitClose := clstags.has(n)
	if !inln && !omitClose {
		p.endline()
	}

	if tagvoid.has(n) {
		if literal || spcKeep {
			panic(fmt.Sprintf("<%s> is both literal/keep-space and void", n.Data))
		}
		return nil
	}

	// Indent if needed and print the children.
	if !inln {
		p.lvl++
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
	if !inln {
		p.lvl--
		p.endline()
	}

	if !omitClose {
		p.Indnt()
		p.write(closeTag(n))
	}
	if literal {
		p.ltrlkeep--
	}
	if spcKeep {
		p.spckeeper--
	}
	if !inln {
		p.endline()
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
	s = escapeTXT(s)

	if p.inspcKeep() {
		p.write(s)
		return nil
	}

	s = collapseText(s, n)
	if s == "" {
		return nil
	}

	p.Indnt()

	startSpace := s[0] == ' '
	endSpace := s[len(s)-1] == ' '

	wrpStart := 0
	if (inlnTags.has(n.PrevSibling) || inlnTags.has(n.Parent)) && !startSpace {
		wrpStart = 1
	}

	words := strings.Fields(strings.TrimSpace(s))
	for i, w := range words {
		if (i == 0 && startSpace) || i != 0 {
			w = " " + w
		}
		if i == len(words)-1 && endSpace && w != " " {
			w = w + " "
		}

		if i < wrpStart {
			p.write(w)
		} else {
			p.wrp(w, "")
		}
	}
	return nil
}

func (p *printer) Indnt() {
	if p.inLiteral() || p.inspcKeep() || !p.lStart {
		return
	}
	s := strings.Repeat(p.iStr, p.lvl)
	p.write(s) // updates lineStart and lineWidth
}

func (p *printer) wrp(s, extra string) {
	if !p.inLiteral() && !p.inspcKeep() &&
		p.wWidth > 0 && p.lWidth+len(s) > p.wWidth {
		p.endline()
		p.Indnt()
		s = extra + strings.TrimLeft(s, " ")
	}
	p.write(s)
}

func (p *printer) endline() {
	if p.inLiteral() || p.inspcKeep() {
		return
	}
	if p.lStart {
		return
	}
	p.write("\n")
	p.lStart = true
	p.lWidth = 0
}

func (p *printer) write(s string) {
	if p.werr != nil {
		return
	}
	_, p.werr = io.WriteString(p.w, s)
	p.lStart = false
	p.lWidth += len(s)
}

func (p *printer) tagOpener(n *html.Node) (forceinln bool) {
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

	inln := inlnTags.has(n)
	wouldwrp := p.wWidth > 0 && p.lWidth+tagLen > p.wWidth
	prev := n.PrevSibling
	prevTextNotSpace := prev != nil && prev.Type == html.TextNode &&
		(prev.Data == "" || !whtspace.MatchString(prev.Data[len(prev.Data)-1:]))
	startSpaceMatters := inlnTags.has(prev) || inlnTags.has(n.Parent) || prevTextNotSpace
	if !inln || (wouldwrp && !startSpaceMatters) {
		p.endline()
	}

	startedLine := p.lStart
	p.Indnt()

	if !ltrltags.has(n) && !p.inLiteral() &&
		!spacekeepTags.has(n) && !p.inspcKeep() {
		childLen := -1
		if n.FirstChild == nil {
			childLen = 0
		} else if hasSingleChild(n) && n.FirstChild.Type == html.TextNode {
			childLen = len(collapseText(escapeTXT(n.FirstChild.Data), n.FirstChild))
		}
		if childLen >= 0 && (p.lWidth+tagLen+childLen+len(closeTag(n)) < p.wWidth || p.wWidth <= 0) {
			forceinln = true
		}
	}

	var unwrpTokens int
	var wrpIndent string
	if startedLine {
		wrpIndent = strings.Repeat(p.iStr, 2)
		unwrpTokens = 1
		if len(tokens[0]) < len(wrpIndent) {
			unwrpTokens = 2
		}
	} else if (inln || forceinln) && startSpaceMatters {
		unwrpTokens = 1
	}
	for i, t := range tokens {
		if i < unwrpTokens {
			p.write(t)
		} else {
			p.wrp(t, wrpIndent)
		}
	}

	return forceinln
}

func hasSingleChild(n *html.Node) bool {
	return n.FirstChild != nil && n.FirstChild == n.LastChild
}

func closeTag(n *html.Node) string {
	if n.Type != html.ElementNode || tagvoid.has(n) || clstags.has(n) {
		return ""
	}
	return "</" + n.Data + ">"
}

func escapeTXT(s string) string {
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	return s
}

var whtspace *regexp.Regexp = regexp.MustCompile(`[\t\n\f\r ]+`)

func collapseText(s string, n *html.Node) string {
	s = whtspace.ReplaceAllString(s, " ")

	if !inlnTags.has(n.Parent) {
		if !inlnTags.has(n.PrevSibling) {
			s = strings.TrimLeft(s, " ")
		}
		if !inlnTags.has(n.NextSibling) {
			s = strings.TrimRight(s, " ")
		}
	}

	return s
}
