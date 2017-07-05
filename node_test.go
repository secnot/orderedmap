package orderedmap

import "testing"

func TestNewNode(t *testing.T) {
	node1 := newNode(1, 10, nil, nil)
	node2 := newNode(2, 55, nil, nil)

	if node1.Next != nil || node1.Prev != nil {
		t.Error("Pointer assignment error")
	}

	node := newNode(3, 33, node1, node2)
	if node.Next != node1 {
		t.Error("Next pointer assignment error")
	}

	if node.Prev != node2 {
		t.Error("Prev pointer assignment error")
	}

	if node.Value != 33 {
		t.Error("Value assignment error")
	}

	if node.Key != 3 {
		t.Error("Key assignment error")
	}
}
