package filetype

import (
	"gopkg.in/h2non/filetype.v0/matchers"
	"gopkg.in/h2non/filetype.v0/types"
	"testing"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		buf []byte
		ext string
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "jpg"},
		{[]byte{0xFF, 0xD8, 0x00}, "unknown"},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, "png"},
	}

	for _, test := range cases {
		match, err := Match(test.buf)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if match.Extension != test.ext {
			t.Fatalf("Invalid image type: %s", match.Extension)
		}
	}
}

func TestMatches(t *testing.T) {
	cases := []struct {
		buf   []byte
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, true},
		{[]byte{0xFF, 0x0, 0x0}, false},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, true},
	}

	for _, test := range cases {
		if Matches(test.buf) != test.match {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}

func TestAddMatcher(t *testing.T) {
	fileType := AddType("foo", "foo/foo")

	AddMatcher(fileType, func(buf []byte) bool {
		return len(buf) == 2 && buf[0] == 0x00 && buf[1] == 0x00
	})

	if !Is([]byte{0x00, 0x00}, "foo") {
		t.Fatalf("Type cannot match")
	}

	if !IsSupported("foo") {
		t.Fatalf("Not supported extension")
	}

	if !IsMIMESupported("foo/foo") {
		t.Fatalf("Not supported MIME type")
	}
}

func TestMatchMap(t *testing.T) {
	cases := []struct {
		buf  []byte
		kind types.Type
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, types.Get("jpg")},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, types.Get("png")},
		{[]byte{0xFF, 0x0, 0x0}, Unknown},
	}

	for _, test := range cases {
		if kind := MatchMap(test.buf, matchers.Image); kind != test.kind {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}

func TestMatchesMap(t *testing.T) {
	cases := []struct {
		buf   []byte
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, true},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, true},
		{[]byte{0xFF, 0x0, 0x0}, false},
	}

	for _, test := range cases {
		if match := MatchesMap(test.buf, matchers.Image); match != test.match {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}