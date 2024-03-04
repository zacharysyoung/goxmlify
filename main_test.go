package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	testCases := []struct {
		input           string
		xml, xhtmlClean string
	}{
		{
			input:      "<a><b></a>",
			xml:        "<a><b></b></a>",
			xhtmlClean: "<a><b></b></a>",
		},
		{
			input:      `<a><b nonce=foo></a>`,
			xml:        `<a><b nonce="foo"></b></a>`,
			xhtmlClean: `<a><b nonce="foo"></b></a>`,
		},
		{
			input:      `<a><b nonce></a>`,
			xml:        `<a><b nonce="nonce"></b></a>`, // xml fills in the attribs value w/attribs name
			xhtmlClean: `<a><b nonce=""></b></a>`,      // xhtml leaves it blanks
		},
		{
			input:      "<a><b>foo</a>",
			xml:        "<a><b>foo</b></a>",
			xhtmlClean: "<a><b>foo</b></a>",
		},
		// -- Line break oddities
		//    linebreak in the middle is OK
		{
			input:      "<a><b>\n</a>",
			xml:        "<a><b>\n</b></a>",
			xhtmlClean: "<a><b>\n</b></a>",
		},
		//    trailing linebreaks alters output of XHTML
		{
			input:      "<a><b></a>\n",
			xml:        "<a><b></b></a>\n",
			xhtmlClean: "<a><b></b></a><b>\n</b>",
		},
	}

	for _, tc := range testCases {
		t.Run("xml", func(t *testing.T) {
			r := strings.NewReader(tc.input)
			b := &bytes.Buffer{}
			decodeXML(r, b)
			got := b.String()
			if got != tc.xml {
				t.Errorf("decodeXML(%s)\n got %s\nwant %s", tc.input, got, tc.xml)
			}
		})
		t.Run("xhtmlClean", func(t *testing.T) {
			r := strings.NewReader(tc.input)
			b := &bytes.Buffer{}
			htmlToXML(r, b, true)
			got := b.String()
			if got != tc.xhtmlClean {
				t.Errorf("htmlToXML(%s, true)\n got %s\nwant %s", tc.input, got, tc.xhtmlClean)
			}
		})
	}
}

func TestXHTML(t *testing.T) {
	testCases := []struct {
		input             string
		xhtmlClean, xhtml string
	}{
		// link must appear inside <head> for HTML
		{
			input:      `<link nonce>`,
			xhtmlClean: `<link nonce=""></link>`,
			xhtml:      `<html><head><link nonce=""></link></head><body></body></html>`,
		},
	}

	for _, tc := range testCases {
		t.Run("xhtml", func(t *testing.T) {
			r := strings.NewReader(tc.input)
			b := &bytes.Buffer{}
			htmlToXML(r, b, false)
			got := b.String()
			if got != tc.xhtml {
				t.Errorf("htmlToXML(%s, false)\n got %s\nwant %s", tc.input, got, tc.xhtml)
			}
		})
		t.Run("xhtmlClean", func(t *testing.T) {
			r := strings.NewReader(tc.input)
			b := &bytes.Buffer{}
			htmlToXML(r, b, true)
			got := b.String()
			if got != tc.xhtmlClean {
				t.Errorf("htmlToXML(%s, false)\n got %s\nwant %s", tc.input, got, tc.xhtmlClean)
			}
		})
	}
}
