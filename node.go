package ordered_map

type node struct {
	Key   interface{}
	Value interface{}
	Next *node
	Prev *node
}

func newNode(key interface{}, value interface{}, next *node, prev *node) *node {
	return &node{key, value, next, prev}
}
