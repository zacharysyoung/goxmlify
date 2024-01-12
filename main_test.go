package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHTMLToXML(t *testing.T) {
	const (
		htmlPre  = "<html><head>"
		htmlPost = "</body></html>"
	)
	for _, tc := range []struct {
		html, xml string
	}{
		{
			html: `<link nonce>`,
			xml:  `<link nonce=""></link>`,
		},
	} {
		r := strings.NewReader(tc.html)
		b := &bytes.Buffer{}
		htmlToXML(r, b)

		got := b.String()

		if strings.Index(got, htmlPre) != 0 {
			t.Errorf("htmlToXML(`%s`)\nstarts with %s\nwant       %s", tc.html, got[:10], htmlPre)
		}

		if !strings.Contains(got, tc.xml) {
			t.Errorf("htmlToXML(`%s`)→%q does not contain %q", tc.html, got, tc.xml)
		}
	}
}

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
