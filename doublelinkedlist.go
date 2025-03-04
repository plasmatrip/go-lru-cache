package lrucache

import "fmt"

type DoubleLinkedList[T any] interface {
	Head() *Node[T]
	Tail() *Node[T]
	Push(value T)
	Remove(node *Node[T])
	RemoveTail()
	MoveToHead(node *Node[T])
	String() string
}

type Node[T any] struct {
	data T
	next *Node[T]
	prev *Node[T]
}

type DoubleLinkedListImpl[T any] struct {
	head *Node[T]
	tail *Node[T]
	cap  int
	len  int
}

func NewDoubleLinkedList[T any](cap int) *DoubleLinkedListImpl[T] {
	head := &Node[T]{}
	tail := &Node[T]{}
	head.next, tail.prev = tail, head

	return &DoubleLinkedListImpl[T]{
		head: head,
		tail: tail,
		cap:  cap,
		len:  0,
	}
}

func (d *DoubleLinkedListImpl[T]) MoveToHead(node *Node[T]) {
	d.Remove(node)
	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
}

func (d *DoubleLinkedListImpl[T]) Push(value T) {
	if d.len == d.cap {
		d.RemoveTail()
	}

	node := &Node[T]{data: value}

	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
	d.len++
}

func (d *DoubleLinkedListImpl[T]) Remove(node *Node[T]) {
	if d.len == 0 {
		return
	}

	node.prev.next, node.next.prev = node.next, node.prev
	d.len--
}

func (d *DoubleLinkedListImpl[T]) RemoveTail() {
	if d.len == 0 {
		return
	}

	d.tail.prev, d.tail.prev.prev.next = d.tail.prev.prev, d.tail
	d.len--
}

func (d *DoubleLinkedListImpl[T]) Head() *Node[T] {
	return d.head.next
}

func (d *DoubleLinkedListImpl[T]) Tail() *Node[T] {
	return d.tail.prev
}

func (d *DoubleLinkedListImpl[T]) String() string {
	next := d.head.next
	str := ""
	for next.next != nil {
		str += fmt.Sprintf("%v ", next.data)
		next = next.next
	}
	return "List: [ " + str + "]"
}
