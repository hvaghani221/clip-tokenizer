package main

type Element[K, V any] struct {
	next  *Element[K, V]
	prev  *Element[K, V]
	key   K
	value V
	list  *List[K, V]
}

type List[K, V any] struct {
	root Element[K, V]
	len  int
}

func NewList[K, V any]() *List[K, V] {
	list := &List[K, V]{
		root: Element[K, V]{},
		len:  0,
	}
	list.root.list = list
	list.root.next = &list.root
	list.root.prev = &list.root
	return list
}

func (l *List[K, V]) AddToFront(key K, value V) *Element[K, V] {
	elem := &Element[K, V]{
		key:   key,
		value: value,
		next:  l.root.next,
		prev:  &l.root,
		list:  l,
	}
	l.root.next.prev = elem
	l.root.next = elem
	l.len += 1
	return elem
}

func (l *List[K, V]) MoveToFront(e *Element[K, V]) {
	if e.list != l || l.root.next == e {
		return
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = &l.root
	e.next = l.root.next

	e.prev.next = e
	e.next.prev = e
}

func (l *List[K, V]) RemoveBack() (k K, v V) {
	if l.len == 0 {
		return
	}
	elem := l.root.prev
	elem.prev.next = &l.root
	l.root.prev = elem.prev
	elem.prev = nil
	elem.next = nil
	l.len -= 1
	return elem.key, elem.value
}

func (l *List[K, V]) Len() int {
	return l.len
}

func (l *List[K, V]) GetList() []K {
	keys := make([]K, 0, l.Len())
	vals := make([]V, 0, l.Len())
	head := l.root.next
	for i := 0; i < l.Len(); i++ {
		keys = append(keys, head.key)
		vals = append(vals, head.value)
		head = head.next
	}
	_ = vals
	return keys
}

type LRU[K comparable, V any] struct {
	queue *List[K, V]
	items map[K]*Element[K, V]
	cap   int
}

func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		queue: NewList[K, V](),
		items: make(map[K]*Element[K, V]),
		cap:   capacity,
	}
}

func (l LRU[K, V]) Get(key K) (value V, found bool) {
	e, found := l.items[key]
	if !found {
		return
	}
	l.queue.MoveToFront(e)
	return e.value, true
}

func (l LRU[K, V]) Put(key K, val V) {
	if element, found := l.items[key]; found {
		l.queue.MoveToFront(element)
		element.value = val
		return
	}
	if l.cap == l.queue.Len() {
		k, _ := l.queue.RemoveBack()
		delete(l.items, k)
	}

	e := l.queue.AddToFront(key, val)
	l.items[key] = e
}
