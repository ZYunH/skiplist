package skiplist

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	l := New(10,0.3,2)
	l.Insert(12, "12")
	l.Insert(13, "13")
	l.Insert(14, "14")
	l.Insert(15, "15")
	l.Insert(11, "11")
	l.Insert(16, "16")
	l.Insert(17, "17")
	printSkipList(l)
	l.Delete(17, "17")
	printSkipList(l)
}

func BenchmarkSkipList_Insert(b *testing.B) {
	l := NewDefault()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Insert(int64(i), strconv.Itoa(i))
	}
}
