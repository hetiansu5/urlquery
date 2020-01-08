package query

import (
	"testing"
)

func Test_unpackKey_NotEnd(t *testing.T) {
	key := "hts"
	pre, suf := unpackKey(key)
	if pre != "hts" || suf != "" {
		t.Error("unpack error")
	}
}

func Test_unpackKey_RightSquareBracketEnd(t *testing.T) {
	key := "hts%5B0%5D"
	pre, suf := unpackKey(key)
	if pre != "hts" || suf != "%5B0%5D" {
		t.Error("unpack error")
	}
}

func Test_unpackKey_LeftSquareBracketEnd(t *testing.T) {
	key := "%5Bhts%5D%5B0%5D"
	pre, suf := unpackKey(key)
	if pre != "hts" || suf != "%5B0%5D" {
		t.Error("unpack error")
	}
}
