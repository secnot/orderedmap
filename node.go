package ordered_map


// An element of an OrderedDict, forms a linked list ordered by insertion time
type node struct {
	Key   interface{}
	Value interface{}
	Next *node
	Prev *node
}

// Create new node
func newNode(key interface{}, value interface{}, next *node, prev *node) *node {
	return &node{key, value, next, prev}
}
