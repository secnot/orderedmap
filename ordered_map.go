package ordered_map

import "fmt"

type OrderedMap struct {
	table map[interface{}]*node
	root *node	// Point to the newest node
}

func NewOrderedMap() *OrderedMap{
	root := newNode(nil, nil, nil, nil) // sentinel Node
	root.Next, root.Prev = root, root

	om := &OrderedMap{
		table: make(map[interface{}]*node),
		root: root,
	}
	return om
}


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


func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if node, ok := om.table[key]; !ok {
		return nil, false
	} else {
		return node.Value, true
	}
}


// Return key/value/ok for the last node
func (om *OrderedMap) GetLast()(interface{}, interface{}, bool) {
	if len(om.table) == 0 {
		return nil, nil, false
	}
	node := om.root.Prev
	return node.Key, node.Value, true
}

// Return key/value/ok for the first node
func (om *OrderedMap) GetFirst()(interface{}, interface{}, bool) {
	if len(om.table) == 0 {
		return nil, nil, false
	}

	node := om.root.Next
	return node.Key, node.Value, true
}



func (om *OrderedMap) Delete(key interface{}) {
	if node, ok := om.table[key]; ok{
		node.Next.Prev = node.Prev
		node.Prev.Next = node.Next
		
		// TODO: Commented while testing
		//node.Next, node.Prev = om.root, om.root //Stop iteration gracefully
		
		delete(om.table, key)
	}
}


func (om *OrderedMap) Len() int {
	return len(om.table)
}


// Pop newest or oldest element from the OrderedMap
// returns key, value, ok
// last = false -> FIFO
// last = true  -> LIFO
func (om *OrderedMap) Pop(last bool) (key interface{}, value interface{}, ok bool){

	var anode *node

	// Empty OrderedMap
	if om.root == om.root.Next {
		return nil, nil, false
	}

	if last {
		// Pop newest node
		anode = om.root.Prev
	} else {
		// Pop oldest node
		anode = om.root.Next
	}
	
	anode.Next.Prev = anode.Prev
	anode.Prev.Next = anode.Next
	anode.Next, anode.Prev = om.root, om.root //Stop iteration gracefully
	delete(om.table, anode.Key)

	return anode.Key, anode.Value, true
}

func (om *OrderedMap) PopLast()(key interface{}, value interface{}, ok bool){
	return om.Pop(true)
}

func (om *OrderedMap) PopFirst()(key interface{}, value interface{}, ok bool){
	return om.Pop(false)
}


// Move an existing key to either end of an OrderedMap
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

func (om *OrderedMap) MoveLast(key interface{}) bool{
	return om.Move(key, true)
}

func (om *OrderedMap) MoveFirst(key interface{}) bool{
	return om.Move(key, false)
}



type MapIterator struct {
	curr *node
	root *node
	reverse bool
}


// Next element of the iteration
// return key, value, ok
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


func (om *OrderedMap) Iter() *MapIterator{
	return &MapIterator{
		curr: om.root, 
		root: om.root,
		reverse: false,
	}
}

func (om *OrderedMap) IterReverse() *MapIterator{
	return &MapIterator{
		curr: om.root, 
		root: om.root,
		reverse: true,
	}
}

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
