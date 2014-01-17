package utils

import "fmt"

type Node struct {
	Value string
}

func (node *Node) String() string {
	return fmt.Sprintf("Node value:%v\n", node.Value)
}

//Queue is a basic FIFO queue that resizes as needed.
type Queue struct {
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}

func (queue Queue) String() string {
	return fmt.Sprintf("Nodes:\n%v\nSize: %d, Head: %d,Tail: %d, Count: %d", queue.nodes, queue.size, queue.head, queue.tail, queue.count)
}

//Get number of the remaining nodes in a queue
func (queue Queue) Length() int {
	return queue.count
}

//NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}
