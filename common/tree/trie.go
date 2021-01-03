package tree

import (
	"fmt"
	"strconv"
)

type Trie struct {
	IsWord   bool
	Children map[rune]*Trie
}

/** Initialize your data structure here. */
func Constructor() *Trie {
	return &Trie{IsWord: false, Children: make(map[rune]*Trie)}
}

func NewByInt(s []int) *Trie {
	t := Constructor()
	if s == nil {
		return t
	}
	for _, i := range s {
		t.Insert(strconv.Itoa(i))
	}
	return t
}

func NewByInt32(s []int32) *Trie {
	t := Constructor()
	if s == nil {
		return t
	}
	for _, i := range s {
		t.Insert(fmt.Sprintf("%d", i))
	}
	return t
}

func NewByInt64(s []int64) *Trie {
	t := Constructor()
	if s == nil {
		return t
	}
	for _, i := range s {
		t.Insert(fmt.Sprintf("%d", i))
	}
	return t
}

func NewByString(s []string) *Trie {
	t := Constructor()
	if s == nil {
		return t
	}
	for _, i := range s {
		t.Insert(i)
	}
	return t
}

/** Inserts a word into the trie. */
func (t *Trie) Insert(word string) {
	root := t
	for _, item := range word {
		p := item - 'a'
		if root.Children[p] == nil {
			root.Children[p] = &Trie{IsWord: false, Children: make(map[rune]*Trie)}
		}
		root = root.Children[p]
	}
	root.IsWord = true
}

/** Returns if the word is in the trie. */
func (t *Trie) Search(word string) bool {
	root := t

	for _, item := range word {
		p := item - 'a'

		if root.Children[p] == nil {
			return false
		}
		root = root.Children[p]
	}
	return root.IsWord && root != nil
}

/** Returns if there is any word in the trie that starts with the given prefix. */
func (t *Trie) StartsWith(prefix string) bool {

	root := t

	for _, item := range prefix {
		p := item - 'a'

		if root.Children[p] == nil {
			return false
		}
		root = root.Children[p]
	}
	return root != nil
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */
