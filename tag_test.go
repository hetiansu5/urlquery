package query

import (
	"testing"
)

func Test_newTag(t *testing.T) {
	tg := newTag("id ignore vip")

	if tg.name != "id" {
		t.Error("id is wrong")
	}

	if len(tg.flags) != 2 {
		t.Error("flag's length is wrong")
	}

	if _, ok := tg.flags["ignore"]; !ok {
		t.Error("flag ignore is not found")
	}

	if _, ok := tg.flags["vip"]; !ok {
		t.Error("flag vip is not found")
	}
}

func Test_newTag_NotName(t *testing.T)  {
	tg := newTag(" ignore vip")

	if tg.name != "" {
		t.Error("id is wrong")
	}

	if len(tg.flags) != 2 {
		t.Error("flag's length is wrong")
	}

	if _, ok := tg.flags["ignore"]; !ok {
		t.Error("flag ignore is not found")
	}

	if _, ok := tg.flags["vip"]; !ok {
		t.Error("flag vip is not found")
	}
}


func Test_newTag_SpecialSpace(t *testing.T)  {
	tg := newTag(" ignore \t vip ")

	if tg.name != "" {
		t.Error("id is wrong")
	}

	if len(tg.flags) != 2 {
		t.Error("flag's length is wrong")
	}

	if _, ok := tg.flags["ignore"]; !ok {
		t.Error("flag ignore is not found")
	}

	if _, ok := tg.flags["vip"]; !ok {
		t.Error("flag vip is not found")
	}
}