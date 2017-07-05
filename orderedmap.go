// Package orderedmap is a Go implementation of Python's OrderedDict class, a map
// that preserves the order of insertion, so key:value pairs can be iterated in
// the order they where added.
// It can also be used as a stack (LIFO) or queue (FIFO).
package orderedmap

import "fmt"

// OrderedMap class
type OrderedMap struct {
	table map[interface{}]*node
	root  *node
}

// NewOrderedMap creates an empty OrderedMap
func NewOrderedMap() *OrderedMap {
	root := newNode(nil, nil, nil, nil) // sentinel Node
	root.Next, root.Prev = root, root

	om := &OrderedMap{
		table: make(map[interface{}]*node),
		root:  root,
	}
	return om
}

// Len returns the number of elements in the Map
func (om *OrderedMap) Len() int {
	return len(om.table)
}

// Set the key value, if the key overwrites an existing entry, the original
// insertion position is left unchanged, otherwise the key is inserted at the end.
func (om *OrderedMap) Set(key interface{}, value interface{}) {
	if node, ok := om.table[key]; !ok {
		// New Node
		root := om.root
		node := newNode(key, value, root, root.Prev)
		root.Prev.Next = node
		root.Prev = node
		om.table[key] = node
	} else {
		// Update existing node value
		node.Value = value
	}
}

// Get the value of an existing key, leaving the map unchanged
func (om *OrderedMap) Get(key interface{}) (value interface{}, ok bool) {
	if node, isOk := om.table[key]; !isOk {
		value, ok = nil, false
	} else {
		value, ok = node.Value, true
	}
	return
}

// GetLast return the key and value for the last element added, leaving
// the map unchanged
func (om *OrderedMap) GetLast() (key interface{}, value interface{}, ok bool) {
	if len(om.table) == 0 {
		key, value, ok = nil, nil, false
	} else {
		node := om.root.Prev
		key, value, ok = node.Key, node.Value, true
	}
	return
}

// GetFirst returns the key and value for the first element, leaving the map unchanged
func (om *OrderedMap) GetFirst() (key interface{}, value interface{}, ok bool) {
	if len(om.table) == 0 {
		key, value, ok = nil, nil, false
	} else {
		node := om.root.Next
		key, value, ok = node.Key, node.Value, true
	}
	return
}

// Delete a key:value pair from the map.
func (om *OrderedMap) Delete(key interface{}) {
	if node, ok := om.table[key]; ok {
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next

		delete(om.table, key)
	}
}

// Pop and return key:value for the newest or oldest element on the OrderedMap
func (om *OrderedMap) Pop(last bool) (key interface{}, value interface{}, ok bool) {
	if last {
		key, value, ok = om.GetLast()
	} else {
		key, value, ok = om.GetFirst()
	}

	if ok {
		om.Delete(key)
	}
	return
}

// PopLast is a shortcut to Pop the last element
func (om *OrderedMap) PopLast() (key interface{}, value interface{}, ok bool) {
	return om.Pop(true)
}

// PopFirst is a shortcut to Pop the first element
func (om *OrderedMap) PopFirst() (key interface{}, value interface{}, ok bool) {
	return om.Pop(false)
}

// Move an existing key to either the end of the OrderedMap
func (om *OrderedMap) Move(key interface{}, last bool) (ok bool) {

	var moved *node

	// Remove from current position
	anode, ok := om.table[key]
	if !ok {
		return false
	}

	anode.Next.Prev = anode.Prev
	anode.Prev.Next = anode.Next
	moved = anode

	// Insert at the start or end
	root := om.root
	if last {
		moved.Next = root
		moved.Prev = root.Prev
		root.Prev.Next = moved
		root.Prev = moved
	} else {
		moved.Prev = root
		moved.Next = root.Next
		root.Next.Prev = moved
		root.Next = moved
	}

	return true
}

// MoveLast is a shortcut to Move a key to the end o the map
func (om *OrderedMap) MoveLast(key interface{}) (ok bool) {
	return om.Move(key, true)
}

// MoveFirst is a shortcut to Move a key to the beginning of the map
func (om *OrderedMap) MoveFirst(key interface{}) (ok bool) {
	return om.Move(key, false)
}

// MapIterator is a iterator over an OrderedMap
type MapIterator struct {
	curr    *node
	root    *node
	reverse bool
}

// Iter creates a map iterator
func (om *OrderedMap) Iter() *MapIterator {
	return &MapIterator{
		curr:    om.root,
		root:    om.root,
		reverse: false,
	}
}

// IterReverse creates a reverse order map iterator
func (om *OrderedMap) IterReverse() *MapIterator {
	return &MapIterator{
		curr:    om.root,
		root:    om.root,
		reverse: true,
	}
}

// Next key:value pair
func (mi *MapIterator) Next() (key interface{}, value interface{}, ok bool) {

	// Already finished
	if mi.curr == nil {
		return nil, nil, false
	}

	// Advance pointer
	if mi.reverse {
		mi.curr = mi.curr.Prev
	} else {
		mi.curr = mi.curr.Next
	}

	// This is the last iteration
	if mi.curr == mi.root {
		mi.curr = nil
		key, value, ok = nil, nil, false
	} else {
		key, value, ok = mi.curr.Key, mi.curr.Value, true
	}

	return
}

// String interface
func (om *OrderedMap) String() string {
	buffer := make([]string, om.Len())

	iter := om.Iter()
	index := 0
	for key, value, ok := iter.Next(); ok; key, value, ok = iter.Next() {
		buffer[index] = fmt.Sprintf("%v:%v, ", key, value)
		index++
	}
	return fmt.Sprintf("OrderedMap%v", buffer)
}
