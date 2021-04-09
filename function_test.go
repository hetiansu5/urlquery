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

func Test_repackArrayQueryKey(t *testing.T) {
	key := "[hts][0]"
	target := repackArrayQueryKey(key)
	if target != "[hts][]" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_repackArrayQueryKey1(t *testing.T) {
	key := "hts]"
	target := repackArrayQueryKey(key)
	if target != "hts]" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_repackArrayQueryKey2(t *testing.T) {
	key := "[hts"
	target := repackArrayQueryKey(key)
	if target != "[hts" {
		t.Error("failed to execute repackArrayQueryKey function")
	}
}

func Test_genNextParentNode(t *testing.T) {
	if genNextParentNode("", "test") != "test" {
		t.Error("failed to execute genNextParentNode")
	}

	if genNextParentNode("p", "test") != "p[test]" {
		t.Error("failed to execute genNextParentNode")
	}
}
