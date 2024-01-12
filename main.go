package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

func main() {
	var in io.Reader
	switch len(os.Args) {
	case 1:
		in = os.Stdin
	case 2:
		in = must(os.Open(os.Args[1]))
	default:
		fmt.Fprintln(os.Stderr, "usage: xmlify [input.hmtl]")
		os.Exit(1)
	}

	doc := must(html.Parse(in))
	w := bufio.NewWriter(os.Stdout)
	htmlToXML(doc, w)
	_must(w.Flush())
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

func htmlToXML(doc *html.Node, w io.Writer) {
	w = bufio.NewWriter(w)
	e := xml.NewEncoder(w)

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			e.EncodeToken(startElem(n))
		case html.TextNode:
			e.EncodeToken(charData(n))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		if n.Type == html.ElementNode {
			e.EncodeToken(endElem(n))
		}
	}
	f(doc)

	_must(w.(*bufio.Writer).Flush())
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
