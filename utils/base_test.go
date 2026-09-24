package utils

import (
	"testing"
)

// Node 双向链表节点
type Node struct {
	Value int
	Prev  *Node
	Next  *Node
}

// Link2Point 双向链表
type Link2Point struct {
	List []*Node
}

func (l *Link2Point) push(value int) *Node {
	newNode := &Node{
		Value: value,
	}

	if len(l.List) != 0 {
		var firstN *Node
		for _, n := range l.List {
			if n.Prev == nil {
				firstN = n
				break
			}
		}

		currentN := firstN
		var lastN *Node
		for {
			if currentN.Next != nil {
				currentN = currentN.Next
			} else {
				lastN = currentN
				break
			}
		}

		lastN.Next = newNode
	} else {
		l.List = append(l.List, newNode)
	}

	return newNode
}

func TestLink2Point(t *testing.T) {
	l := &Link2Point{}
	if got := l.push(1); got == nil || got.Value != 1 {
		t.Fatalf("push(1) = %+v", got)
	}

}
