package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var in io.Reader
	switch len(os.Args) {
	case 2:
		buf := &bytes.Buffer{}
		io.Copy(buf, os.Stdin)
		s := strings.TrimSpace(buf.String()) // trim trailing linebreak, see "linereaks test"
		in = strings.NewReader(s)
	case 3:
		in = must(os.Open(os.Args[2]))
	default:
		fmt.Fprintln(os.Stderr, "usage: xmlify xml|html [input.html]")
		os.Exit(1)
	}

	var f func(r io.Reader, w io.Writer)
	switch os.Args[1] {
	case "xml":
		f = decodeXML
	case "html":
		f = htmlToXML
	default:
		fmt.Fprintln(os.Stderr, "usage: xmlify xml|html [input.html]")
		os.Exit(1)
	}

	f(in, os.Stdout)
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
// XML.
func htmlToXML(r io.Reader, w io.Writer) {
	doc := must(html.Parse(r))

	bw := bufio.NewWriter(w)
	enc := xml.NewEncoder(bw)
	defer func() { _must(bw.Flush()) }()

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			enc.EncodeToken(startElem(n))
		case html.TextNode:
			enc.EncodeToken(charData(n))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			enc.EncodeToken(endElem(n))
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
