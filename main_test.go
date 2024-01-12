package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestHTMLToXML(t *testing.T) {
	for _, tc := range []struct {
		html, xml string
	}{
		{
			html: `<html><head><link nonce></head><body></body></html>`,
			xml:  `<html><head><link nonce=""></link></head><body></body></html>`,
		},
	} {
		doc := must(html.Parse(strings.NewReader(tc.html)))
		b := &bytes.Buffer{}

		htmlToXML(doc, b)

		if got := b.String(); got != tc.xml {
			t.Errorf("htmlToXML(`%s`)\ngot  %s\nwant %s", tc.html, got, tc.xml)
		}
	}
}
