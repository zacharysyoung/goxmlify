package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecodeBadXML(t *testing.T) {
	for _, tc := range []struct {
		xml, want string
	}{
		{
			xml:  "<a><b></a>",
			want: "<a><b></b></a>",
		},
		{
			xml:  `<a><b nonce></a>`,
			want: `<a><b nonce="nonce"></b></a>`,
		},
		{
			xml:  `<a><b nonce/></a>`,
			want: `<a><b nonce="nonce"></b></a>`,
		},
		{
			xml:  `<a><b nonce=foo></a>`,
			want: `<a><b nonce="foo"></b></a>`,
		},
	} {
		r := strings.NewReader(tc.xml)
		b := &bytes.Buffer{}
		decodeXML(r, b)

		if got := b.String(); got != tc.want {
			t.Errorf("decodeXML(%q)\ngot  %q\nwant %q", tc.xml, got, tc.want)
		}
	}
}

// Weird command-line behavior cannot replicate in a test, yet:
//
//	echo '<a><b nonce></a>' | goxmlify html
//
// yields:
//
//	<html><head></head><body><a><b nonce=""></b></a><b nonce="">
//	</b></body></html>
//
// with extra <b nonce="">\n</b>
// ??
func TestHTMLToXML(t *testing.T) {
	const (
		htmlPre  = "<html><head>"
		htmlMid  = "</head><body>"
		htmlPost = "</body></html>"
	)

	for _, tc := range []struct {
		html, xml string
	}{
		{
			html: `<link nonce>`,
			xml:  htmlPre + `<link nonce=""></link>` + htmlMid + htmlPost,
		},
		{
			html: `<a><b></a>`,
			xml:  htmlPre + htmlMid + `<a><b></b></a>` + htmlPost,
		},
		{
			html: `<a><b nonce></a>`,
			xml:  htmlPre + htmlMid + `<a><b nonce=""></b></a>` + htmlPost,
		},
		// linebreaks drastically alter effect of HTML parser/renderer
		{
			html: "<a><b nonce></a>\n",
			xml:  htmlPre + htmlMid + "<a><b nonce=\"\"></b></a><b nonce=\"\">\n</b>" + htmlPost,
		},
	} {
		r := strings.NewReader(tc.html)
		b := &bytes.Buffer{}
		htmlToXML(r, b)

		got := b.String()

		if got != tc.xml {
			t.Errorf("htmlToXML(`%s`)\ngot  %q\nwant %q", tc.html, got, tc.xml)
		}
	}
}
