package skiplist

import (
	"math/rand"
)

const (
	defaultMaxLevel = 32
	defaultP        = 0.25
	defaultRandSeed = 0
)

type SkipList struct {
	header *node
	tail   *node
	length int64
	level  int

	maxLevel int
	p        float64
	rnd      *rand.Rand
}

func New(maxlevel int, p float64, randseed int64) *SkipList {
	if maxlevel <= 1 || p <= 0 {
		panic("maxLevel must greater than 1, p must greater than 0")
	}

	s := &SkipList{
		header:   nil,
		tail:     nil,
		length:   0,
		level:    1,
		maxLevel: maxlevel,
		p:        p,
		rnd:      rand.New(rand.NewSource(randseed)),
	}

	s.header = newNode(s.maxLevel, 0, "")

	for j := 0; j < s.maxLevel; j++ {
		s.header.levels[j].next = nil
		s.header.levels[j].span = 0
	}

	s.header.pre = nil
	s.tail = nil
	return s
}

func NewDefault() *SkipList {
	return New(defaultMaxLevel, defaultP, defaultRandSeed)
}

func (s *SkipList) randomLevel() int {
	level := 1

	for s.rnd.Float64() < s.p {
		level += 1
	}

	if level > s.maxLevel {
		return s.maxLevel
	}
	return level
}

func (s *SkipList) Insert(score int64, val string) *node {
	update := make([]*node, s.maxLevel)
	rank := make([]uint, s.maxLevel)

	// Search the insert location, also calculates `update` and `rank`.
	// The search process is begin from the highest level's header.
	for i, n := s.level-1, s.header; i >= 0; i-- {
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for n.levels[i].next != nil &&
			(n.levels[i].next.score < score ||
				n.levels[i].next.score == score && strcmp(n.levels[i].next.val, val) < 0) {
			rank[i] += n.levels[i].span
			n = n.levels[i].next
		}
		update[i] = n
	}

	// Make a random level for the insert node.
	level := s.randomLevel()
	// If the insert process will create new levels, we need to
	// update the `rank` and `update`.
	if level > s.level {
		for i := s.level; i < level; i++ {
			// s.header is the only node in every levels,
			// since it doesn't has tail, so its pan is
			// the length of skiplist.
			update[i] = s.header
			update[i].levels[i].span = uint(s.length)
		}
		s.level = level
	}

	// Insert the new node into levels. Keep in mind here, we just
	// insert it to `node.levels`(it only includes next pointer).
	// But the level[0] is actually a doubled link list.
	n := newNode(level, score, val)
	for i := 0; i < level; i++ {
		n.levels[i].next = update[i].levels[i].next
		update[i].levels[i].next = n

		// (rank[0] - rank[i]) is actually the number of nodes between
		// update[i] and the new node in level i.
		n.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = 1 + (rank[0] - rank[i])
	}

	// Increment span for untouched levels, if the new node's level is
	// less than the skiplist's level.
	for i := level; i < s.level; i++ {
		update[i].levels[i].span += 1
	}

	// Update new node's pre.
	if update[0] != s.header {
		n.pre = update[0]
	}

	// Update new node's next's pre, Because the levels[0] is
	// doubled link list. But if new node's next is NIL, we
	// need to change s.tail to the new node.
	if n.levels[0].next != nil {
		n.levels[0].next.pre = n
	} else {
		s.tail = n
	}

	s.length += 1

	return n
}

func (s *SkipList) Delete(score int64, val string) bool {
	update := make([]*node, s.maxLevel)
	n := s.header
	for i := s.level - 1; i >= 0; i-- {
		for n.levels[i].next != nil &&
			(n.levels[i].next.score < score ||
				(n.levels[i].next.score == score && strcmp(n.levels[i].next.val, val) < 0)) {
			n = n.levels[i].next
		}
		update[i] = n
	}

	n = n.levels[0].next
	if n != nil && n.score == score && n.val == val {
		s.delete(n, update)
		return true
	}
	return false
}

func (s *SkipList) delete(n *node, update []*node) {
	// Delete node and update span for all levels.
	for i := 0; i < s.level; i++ {
		if update[i].levels[i].next == n {
			update[i].levels[i].next = n.levels[i].next
			update[i].levels[i].span += n.levels[i].span - 1
		} else {
			update[i].levels[i].span -= 1
		}
	}

	// Update n.next.pre if possible.
	if n.levels[0].next != nil {
		n.levels[0].next.pre = n.pre
	} else {
		s.tail = n.pre
	}

	// Update skiplist.level if some levels only includes header.
	for s.level > 1 && s.header.levels[s.level-1].next == nil {
		s.level -= 1
	}
	s.length -= 1
}

func (s *SkipList) Len() int64 {
	return s.length
}

func (s *SkipList) Head() *node {
	if s.length == 0 {
		return nil
	}
	return s.header.levels[0].next
}

func (s *SkipList) Tail() *node { return s.tail }

type node struct {
	val    string
	score  int64
	pre    *node
	levels []_nodeLevel
}

func newNode(level int, score int64, val string) *node {
	return &node{
		val:    val,
		score:  score,
		pre:    nil,
		levels: make([]_nodeLevel, level),
	}
}

func (n *node) Val() string { return n.val }

func (n *node) Score() int64 { return n.score }

type _nodeLevel struct {
	next *node
	span uint
}
