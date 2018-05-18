package main

import (
	"fmt"
)

// Represent a node
type Node struct {
	data interface{}
	next *Node
	prev *Node
}

// Represent a linked list
type List struct {
	head *Node
	tail *Node
}

//inserting data into linked list
func (L *List) insert(data interface{}) {
	newNode := &Node{
		data: data,
		next: nil,
		prev: nil,
	}

	if L.head == nil && L.tail == nil {
		L.head = newNode
		L.tail = newNode
	} else {
		L.tail.next = newNode
		temp := L.tail
		L.tail = newNode
		newNode.prev = temp
	}
}

//Printing the linked list
func (L *List) Printing() {
	for temp := L.head; temp != nil; temp = temp.next {
		fmt.Printf("%v\n", temp.data)
	}
}

func main() {
	L := &List{}
	L.insert(4)
	L.insert(5)
	L.insert(6)
	L.insert(7)
	L.Printing()
}
