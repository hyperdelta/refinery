package trie

import (
	"github.com/hyperdelta/refinery/log"
	"fmt"
)

var (
	logger = log.Get()
)

type Trie struct {
	root *TrieElement
}

type TrieElement struct {
	prefix string
	data interface{}
	children []*TrieElement
}

func NewTrie() *Trie {
	t := new(Trie)
	t.root = new(TrieElement)

	return t
}

func NewTrieElement() *TrieElement {
	elem := new(TrieElement)
	return elem
}

func (t* Trie) Print() {
	fmt.Printf("\nprint = %s\n", t)
}

func (t* Trie) Clear() {
	t.root = new(TrieElement)
}

func (t* Trie) Add(data interface{}, prefix ...string) {
	elem := t.root._retrieveElement(prefix...)

	if elem == nil {
		t.root.Add(data, prefix...)
	} else {
		elem.data = data
	}
}

func (t* Trie) Retrieve(prefix ...string) interface{} {
	return t.root.Retrieve(prefix...)
}

func (e* TrieElement) Retrieve(prefix ...string) interface{} {
	elem := e._retrieveElement(prefix...)

	if elem != nil {
		return elem.data
	}

	return nil
}

func (e *TrieElement) _retrieveElement(prefix ...string) *TrieElement {
	if prefix == nil || len(prefix) == 0 {
		return e
	}

	var p string = prefix[0]

	if e.children != nil {
		for _, elem := range e.children {
			if elem.prefix == p {
				if len(prefix) > 1 {
					return elem._retrieveElement(prefix[1:]...)
				} else {
					return elem._retrieveElement(nil...)
				}
			}
		}
	}

	return nil
}

func (e* TrieElement) Add(data interface{}, prefix ...string) {

	if prefix == nil || len(prefix) == 0 {
		e.data = data
		return
	}

	var p string = prefix[0]
	var findChild bool = false

	if e.children != nil {
		for _, elem := range e.children {
			if elem.prefix == p {
				findChild = true

				if len(prefix) > 1 {
					elem.Add(data, prefix[1:]...)
				} else {
					elem.Add(data, nil...)
				}
			}
		}
	}

	if !findChild {
		newElem := NewTrieElement()
		newElem.prefix = p
		e.children = append(e.children, newElem)

		if len(prefix) > 1 {
			newElem.Add(data, prefix[1:]...)
		} else {
			newElem.Add(data, nil...)
		}
	}

}