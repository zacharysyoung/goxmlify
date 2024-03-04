package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	method = flag.String("method", "", "process as xml, xhtml, or xhtml-clean; html-clean strips HTML <html>, <head>, and <body> tags.")
)

func main() {
	usage := func() {
		fmt.Fprintln(os.Stderr, "usage: xmlify -method=xml|xhtml|xhtml-clean [input.html]")
		os.Exit(1)
	}

	flag.Parse()
	if *method == "" {
		usage()
	}

	var in io.Reader
	switch args := flag.Args(); len(args) {
	case 0:
		buf := &bytes.Buffer{}
		io.Copy(buf, os.Stdin)
		s := strings.TrimSpace(buf.String()) // trim trailing linebreak, see "linereaks test"
		in = strings.NewReader(s)
	case 1:
		in = must(os.Open(args[0]))
	default:
		usage()
	}

	switch *method {
	case "xml":
		decodeXML(in, os.Stdout)
	case "xhtml":
		htmlToXML(in, os.Stdout, false)
	case "xhtml-clean":
		htmlToXML(in, os.Stdout, true)
	default:
		usage()
	}
}

// decodeXML attemtps to decode XML from r and re-encode to w
// fixing some bad XML along the way.
func decodeXML(r io.Reader, w io.Writer) {
	dec := xml.NewDecoder(r)
	dec.Strict = false

	bw := bufio.NewWriter(w)
	enc := xml.NewEncoder(bw)
	defer func() { _must(bw.Flush()) }()

	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			_must(err)
		}
		enc.EncodeToken(t)
	}
}

// htmlToXML uses the html packages tokenizer to parse malformed
// HTML (as XML), and then render back to HTML (XML)... except
// it enforces some HTML properties that have nothing to do with
// XML. clean strips ancillary HTML <html> <head> and <body> tags.
func htmlToXML(r io.Reader, w io.Writer, clean bool) {
	skip := func(n *html.Node) bool {
		return clean && (n.Data == "html" || n.Data == "head" || n.Data == "body")
	}

	doc := must(html.Parse(r))

	bw := bufio.NewWriter(w)
	enc := xml.NewEncoder(bw)
	defer func() { _must(bw.Flush()) }()

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			if !skip(n) {
				_must(enc.EncodeToken(startElem(n)))
			}
		case html.TextNode:
			_must(enc.EncodeToken(charData(n)))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			if !skip(n) {
				_must(enc.EncodeToken(endElem(n)))
			}
		}
	}
	f(doc)
}

func startElem(n *html.Node) xml.StartElement {
	e := xml.StartElement{
		Name: xml.Name{Space: "", Local: n.Data},
	}
	for _, x := range n.Attr {
		e.Attr = append(e.Attr, xml.Attr{Name: xml.Name{Space: "", Local: x.Key}, Value: x.Val})
	}
	return e
}

func endElem(n *html.Node) xml.EndElement {
	return xml.EndElement{
		Name: xml.Name{Space: "", Local: n.Data},
	}
}

func charData(n *html.Node) xml.CharData {
	return []byte(n.Data)
}

func _must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func must[T any](obj T, err error) T {
	_must(err)
	return obj
}
