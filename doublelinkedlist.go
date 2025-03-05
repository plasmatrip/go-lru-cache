package lrucache

import "fmt"

type DoubleLinkedList[T any] interface {
	Head() *node[T]
	Tail() *node[T]
	Push(value T)
	Remove(node *node[T])
	RemoveTail()
	MoveToHead(node *node[T])
	String() string
}

type node[T any] struct {
	data T
	next *node[T]
	prev *node[T]
}

type doubleLinkedList[T any] struct {
	head *node[T]
	tail *node[T]
	cap  int
	len  int
}

func NewDoubleLinkedList[T any](cap int) *doubleLinkedList[T] {
	head := &node[T]{}
	tail := &node[T]{}
	head.next, tail.prev = tail, head

	return &doubleLinkedList[T]{
		head: head,
		tail: tail,
		cap:  cap,
		len:  0,
	}
}

func (d *doubleLinkedList[T]) MoveToHead(node *node[T]) {
	d.Remove(node)
	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
}

func (d *doubleLinkedList[T]) Push(value T) {
	if d.len == d.cap {
		d.RemoveTail()
	}

	node := &node[T]{data: value}

	d.head.next, node.next, node.prev, d.head.next.prev = node, d.head.next, d.head, node
	d.len++
}

func (d *doubleLinkedList[T]) Remove(node *node[T]) {
	if d.len == 0 {
		return
	}

	node.prev.next, node.next.prev = node.next, node.prev
	d.len--
}

func (d *doubleLinkedList[T]) RemoveTail() {
	if d.len == 0 {
		return
	}

	d.tail.prev, d.tail.prev.prev.next = d.tail.prev.prev, d.tail
	d.len--
}

func (d *doubleLinkedList[T]) Head() *node[T] {
	return d.head.next
}

func (d *doubleLinkedList[T]) Tail() *node[T] {
	return d.tail.prev
}

func (d *doubleLinkedList[T]) String() string {
	next := d.head.next
	str := ""
	for next.next != nil {
		str += fmt.Sprintf("%v ", next.data)
		next = next.next
	}
	return "List: [ " + str + "]"
}
