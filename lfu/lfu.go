package lfu

type LFUCache struct {
	Cache               map[string]*Node
	freq                map[int]*DoubleList
	ncap, size, minFreq int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		Cache: make(map[string]*Node),
		freq:  make(map[int]*DoubleList),
		ncap:  capacity,
	}
}

func (l *LFUCache) Get(key string) int {
	if node, ok := l.Cache[key]; ok {
		l.IncFreq(node)
		return node.val
	}
	return -1
}

func (l *LFUCache) Put(key string, value int) {
	if l.ncap == 0 {
		return
	}
	if node, ok := l.Cache[key]; ok {
		node.val = value
		l.IncFreq(node)
	} else {
		if l.size >= l.ncap {
			node := l.freq[l.minFreq].RemoveLast()
			delete(l.Cache, node.key)
			l.size--
		}
		x := &Node{key: key, val: value, Freq: 1}
		l.Cache[key] = x
		if l.freq[1] == nil {
			l.freq[1] = CreateDL()
		}
		l.freq[1].AddFirst(x)
		l.minFreq = 1
		l.size++
	}
}

func (l *LFUCache) IncFreq(node *Node) {
	_freq := node.Freq
	l.freq[_freq].Remove(node)
	if l.minFreq == _freq && l.freq[_freq].IsEmpty() {
		l.minFreq++
		delete(l.freq, _freq)
	}

	node.Freq++
	if l.freq[node.Freq] == nil {
		l.freq[node.Freq] = CreateDL()
	}
	l.freq[node.Freq].AddFirst(node)
}

type DoubleList struct {
	head, tail *Node
}

type Node struct {
	prev, next *Node
	key        string
	val, Freq  int
}

func CreateDL() *DoubleList {
	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head
	return &DoubleList{
		head: head,
		tail: tail,
	}
}

func (l *DoubleList) AddFirst(node *Node) {
	node.next = l.head.next
	node.prev = l.head

	l.head.next.prev = node
	l.head.next = node
}

func (l *DoubleList) Remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev

	node.next = nil
	node.prev = nil
}

func (l *DoubleList) RemoveLast() *Node {
	if l.IsEmpty() {
		return nil
	}

	last := l.tail.prev
	l.Remove(last)

	return last
}

func (l *DoubleList) IsEmpty() bool {
	return l.head.next == l.tail
}
