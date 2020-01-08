package query

type tag struct {
	name  string
	flags map[string]bool
}

func newTag(origin string) *tag {
	t := &tag{}
	t.flags = map[string]bool{}
	t.init(origin)
	return t
}

func (t *tag) getName() string {
	return t.name
}

func (t *tag) hasFlag(flag string) bool {
	_, ok := t.flags[flag]
	return ok
}

func (t *tag) init(origin string) {
	offset := 0
	size := 0
	for _, c := range []byte(origin + " ") {
		if c <= ' ' && isSpace(c) {
			if size > 0 {
				if offset - size == 0 {
					t.name = origin[offset - size : offset]
				} else {
					t.flags[origin[offset - size : offset]] = true
				}
			}
			size = 0
		} else {
			size++
		}
		offset++
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
