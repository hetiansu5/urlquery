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
	key := "hts[0]"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "[0]" {
		t.Error("unpack error")
	}
}

func Test_unpackQueryKey_LeftSquareBracketEnd(t *testing.T) {
	key := "[hts][0]"
	pre, suf := unpackQueryKey(key)
	if pre != "hts" || suf != "[0]" {
		t.Error("unpack error")
	}
}
