package urlquery

import (
	"testing"
)

func Test_unpackQueryKey_NotEnd(t *testing.T) {
	key := "hts"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "" {
		t.Error("unpack error")
	}
}

func Test_unpackQueryKey_RightSquareBracketEnd(t *testing.T) {
	key := "hts%5B0%5D"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "%5B0%5D" {
		t.Error("unpack error")
	}
}

func Test_unpackQueryKey_LeftSquareBracketEnd(t *testing.T) {
	key := "%5Bhts%5D%5B0%5D"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "%5B0%5D" {
		t.Error("unpack error")
	}
}
