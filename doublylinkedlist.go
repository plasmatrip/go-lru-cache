package golrucache

import "fmt"

type node[T any] struct {
	data T
	next *node[T]
	prev *node[T]
}

type doublyLinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	cap  int
	len  int
}

func NewDoublyLinkedList[T any](cap int) *doublyLinkedList[T] {
	head := &node[T]{}
	tail := &node[T]{}
	head.next, tail.prev = tail, head

	return &doublyLinkedList[T]{
		head: head,
		tail: tail,
		cap:  cap,
		len:  0,
	}
}

func (d *doublyLinkedList[T]) moveToHead(node *node[T]) {
	d.Remove(node)
	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
}

func (d *doublyLinkedList[T]) Push(value T) {
	if d.len == d.cap {
		d.removeTail()
	}

	node := &node[T]{data: value}

	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
	d.len++
}

func (d *doublyLinkedList[T]) Remove(node *node[T]) {
	if d.len == 0 {
		return
	}

	node.prev.next, node.next.prev = node.next, node.prev
	d.len--
}

func (d *doublyLinkedList[T]) removeTail() {
	if d.len == 0 {
		return
	}

	d.tail.prev, d.tail.prev.prev.next = d.tail.prev.prev, d.tail
	d.len--
}

func (d *doublyLinkedList[T]) Head() *node[T] {
	return d.head.next
}

func (d *doublyLinkedList[T]) Tail() *node[T] {
	return d.tail.prev
}

func (d *doublyLinkedList[T]) String() string {
	next := d.head.next
	str := ""
	for next.next != nil {
		str += fmt.Sprintf("%v ", next.data)
		next = next.next
	}
	return "List: [ " + str + "]"
}
