package urlquery

type tag struct {
	name  string
	flags map[string]struct{}
}

func newTag(origin string) *tag {
	t := &tag{}
	t.flags = map[string]struct{}{}
	t.init(origin)
	return t
}

func (t *tag) getName() string {
	return t.name
}

//has flag
func (t *tag) hasFlag(flags ...string) bool {
	for _, flag := range flags {
		_, ok := t.flags[flag]
		if ok {
			return true
		}
	}
	return false
}

func (t *tag) init(origin string) {
	offset := 0
	size := 0
	for _, c := range []byte(origin + " ") {
		if c <= ' ' && isSpace(c) {
			if size > 0 {
				if offset-size == 0 {
					t.name = origin[offset-size : offset]
				} else {
					t.flags[origin[offset-size:offset]] = struct{}{}
				}
			}
			size = 0
		} else {
			size++
		}
		offset++
	}
}

//is space character?
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
