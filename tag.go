package urlquery

type tag struct {
	name    string
	options []string
}

func newTag(origin string) *tag {
	t := &tag{}
	t.options = []string{}
	t.init(origin)
	return t
}

func (t *tag) getName() string {
	return t.name
}

//contains
func (t *tag) contains(options ...string) bool {
	for _, option := range options {
		for _, o := range t.options {
			if o == option {
				return true
			}
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
					t.options = append(t.options, origin[offset-size:offset])
				}
			}
			size = 0
		} else {
			size++
		}
		offset++
	}
}
