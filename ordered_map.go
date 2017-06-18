package ordered_map

import "fmt"

type OrderedMap struct {
	table map[interface{}]*node
	root *node	// Point to the newest node
}


// Create an empty OrderedMap
func NewOrderedMap() *OrderedMap{
	root := newNode(nil, nil, nil, nil) // sentinel Node
	root.Next, root.Prev = root, root

	om := &OrderedMap{
		table: make(map[interface{}]*node),
		root: root,
	}
	return om
}

// Len computes the number of elements in an OrderedMap
func (om *OrderedMap) Len() int {
	return len(om.table)
}


// Sets the key value, if the key overwrites an existing entry, the original
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
func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if node, ok := om.table[key]; !ok {
		return nil, false
	} else {
		return node.Value, true
	}
}

// Get the key and value for the last element added, leaving the map unchanged
func (om *OrderedMap) GetLast()(interface{}, interface{}, bool) {
	if len(om.table) == 0 {
		return nil, nil, false
	}
	node := om.root.Prev
	return node.Key, node.Value, true
}

// Get the key value for the beginning element, leaving the map unchanged
func (om *OrderedMap) GetFirst()(interface{}, interface{}, bool) {
	if len(om.table) == 0 {
		return nil, nil, false
	}

	node := om.root.Next
	return node.Key, node.Value, true
}


// Delete a key:value pair from the map.
func (om *OrderedMap) Delete(key interface{}) {
	if node, ok := om.table[key]; ok{
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next
		
		delete(om.table, key)
	}
}


// Pop and return key:value for the newest or oldest element on the OrderedMap
// returns key, value, ok
// last = false -> FIFO
// last = true  -> LIFO
func (om *OrderedMap) Pop(last bool) (key interface{}, value interface{}, ok bool){
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


// Shortcut to Pop the last element
func (om *OrderedMap) PopLast()(key interface{}, value interface{}, ok bool){
	return om.Pop(true)
}


// Shortcut to Pop the first element
func (om *OrderedMap) PopFirst()(key interface{}, value interface{}, ok bool){
	return om.Pop(false)
}


// Move an existing key to either the end of the OrderedMap
func (om *OrderedMap) Move(key interface{}, last bool) bool{
	
	var moved *node

	// Remove from current position
	if node, ok := om.table[key]; !ok {
		return false
	} else {
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next	
		moved = node
	}

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


// Shortcut to Move an element to the end
func (om *OrderedMap) MoveLast(key interface{}) bool{
	return om.Move(key, true)
}

// Shortcut to Move an element to the beginning
func (om *OrderedMap) MoveFirst(key interface{}) bool{
	return om.Move(key, false)
}


// OrderedMap Iterator
type MapIterator struct {
	curr *node
	root *node
	reverse bool
}

// Create a map iterator
func (om *OrderedMap) Iter() *MapIterator{
	return &MapIterator{
		curr: om.root, 
		root: om.root,
		reverse: false,
	}
}

// Create a reverse ordered map iterator
func (om *OrderedMap) IterReverse() *MapIterator{
	return &MapIterator{
		curr: om.root, 
		root: om.root,
		reverse: true,
	}
}


// Return next key:value pair
func (mi *MapIterator) Next() (interface {}, interface {}, bool) {

	// Already finished
	if mi.curr == nil {
		return nil, nil, false
	}
	
	if mi.reverse {
		mi.curr = mi.curr.Prev
	} else { 
		mi.curr = mi.curr.Next
	}

	// This is the last iteration
	if  mi.curr == mi.root {
		mi.curr = nil
		return nil, nil, false
	}
	
	return mi.curr.Key, mi.curr.Value, true
}


// Stringer interface
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
