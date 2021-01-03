package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie_Search(t *testing.T) {
	/**
	 * Your Trie object will be instantiated and called as such:
	 * obj := Constructor();
	 * obj.Insert(word);
	 * param_2 := obj.Search(word);
	 * param_3 := obj.StartsWith(prefix);
	 */
	assert.Equal(t, NewByInt([]int{1, 2, 3}).Search("1"), true)
	assert.Equal(t, NewByInt([]int{1, 2, 3}).Search(""), false)
	assert.Equal(t, NewByString([]string{"123", "123", "3"}).StartsWith("3"), true)
	assert.Equal(t, NewByString([]string{"123", "123", "3"}).StartsWith("1"), true)
	assert.Equal(t, NewByInt(nil).Search(""), false)
	assert.Equal(t, NewByInt(nil).Search("1"), false)
}
