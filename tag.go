package urlquery

import (
	"strings"
	"sync"
)

// A tag is decoration of struct attribution
type tag struct {
	name    string
	options []string
}

// according to official standard
func newTag(s string) *tag {
	arr := strings.Split(s, ",")
	t := &tag{
		name:    arr[0],
		options: arr[1:],
	}
	return t
}

// get tag name
func (t *tag) getName() string {
	return t.name
}

// contains tag
func (t *tag) contains(option string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	mutex.Unlock()
	for _, o := range t.options {
		if o == option {
			return true
		}
	}
	return false
}
