package urlquery

import (
	"testing"
)

func Test_newTag(t *testing.T) {
	tg := newTag("id,ignore,vip")

	if tg.name != "id" {
		t.Error("name is wrong")
	}

	if len(tg.options) != 2 {
		t.Error("the length of options is wrong")
	}

	if !tg.contains("ignore") {
		t.Error("options ignore is not found")
	}

	if !tg.contains("vip") {
		t.Error("options vip is not found")
	}

	if tg.contains("vip1") {
		t.Error("options vip1 is found")
	}
}

func Test_newTag_NotName(t *testing.T) {
	tg := newTag(",ignore,vip")

	if tg.name != "" {
		t.Error("name is wrong")
	}

	if !tg.contains("ignore") {
		t.Error("options ignore is wrong")
	}

	if !tg.contains("vip") {
		t.Error("options vip is not found")
	}
}

func Test_newTag_SpecialSpace(t *testing.T) {
	tg := newTag(",ignore,vip")

	if tg.getName() != "" {
		t.Error("name is wrong")
	}

	if len(tg.options) != 2 {
		t.Error("options's length is wrong")
	}
}
