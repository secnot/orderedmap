package orderedmap

import (
	"fmt"
	"testing"
)

// Test key present in OrderedMap
func mapHasKey(t *testing.T, om *OrderedMap, key interface{}, value interface{}) bool {

	if v, ok := om.Get(key); v != value || !ok {
		t.Error(fmt.Sprintf("Get(%v) -> expected %v received %v", key, value, v))
		return false
	}
	return true
}

// Test key not present in OrderedMap
func mapNotKey(t *testing.T, om *OrderedMap, key interface{}) bool {

	if v, ok := om.Get(key); v != nil || ok {
		t.Error(fmt.Sprintf("Get(%v) -> shouldn't have a value", key))
		return true
	}
	return false
}

// Test map is empty
func mapIsEmpty(t *testing.T, om *OrderedMap) bool {

	if om.Len() != 0 {
		t.Error("Map is not empty")
		return false
	}

	if key, value, ok := om.PopLast(); key != nil || value != nil || ok {
		t.Error("Map length is 0 but Pop() returned and item")
		return false
	}

	return true
}

func TestNewOrderedMap(t *testing.T) {
	om := NewOrderedMap()
	om.Set(5, 33)
	om.Set(6, 44)
	om.Set("key", "value")

	if val, ok := om.Get(5); val != 33 || !ok {
		t.Error(fmt.Sprintf("Value error, expecting 33 received %v", val))
	}
	if val, ok := om.Get(6); val != 44 || !ok {
		t.Error(fmt.Sprintf("Value error, expecting 44 received %v", val))
	}
	if val, ok := om.Get("key"); val != "value" || !ok {
		t.Error(fmt.Sprintf("Value error, expecting 'value' received %v", val))
	}

	if val, ok := om.Get("not a key"); ok || val != nil {
		t.Error(fmt.Sprintf("Shouldn't have returned %v", val))
	}
}

func TestGetLast(t *testing.T) {
	om := NewOrderedMap()

	if key, value, ok := om.GetLast(); key != nil || value != nil || ok {
		t.Error(fmt.Sprintf("Expecting nil, nil, false -> Returned %v %v %v",
			key, value, ok))
	}

	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.GetLast(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	if key, value, ok := om.GetLast(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	om.MoveLast("one")
	if key, value, ok := om.GetLast(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	if om.Len() != 2 {
		t.Error("Somehow popped an item ..")
	}
}

func TestGetFirst(t *testing.T) {
	om := NewOrderedMap()

	if key, value, ok := om.GetFirst(); key != nil || value != nil || ok {
		t.Error(fmt.Sprintf("Expecting nil, nil, false -> Returned %v %v %v",
			key, value, ok))
	}

	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.GetFirst(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	if key, value, ok := om.GetFirst(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	om.MoveLast("one")
	if key, value, ok := om.GetFirst(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v",
			key, value, ok))
	}

	if om.Len() != 2 {
		t.Error("Somehow popped an item ..")
	}

}

func TestPop(t *testing.T) {
	om := NewOrderedMap()

	om.Set("one", 1)
	om.Set("two", 2)
	om.Set("three", 3)

	// Pop last
	if key, val, ok := om.Pop(true); key != "three" || val != 3 || !ok {
		t.Error("Pop last item error")
	}
	mapHasKey(t, om, "one", 1)
	mapHasKey(t, om, "two", 2)
	mapNotKey(t, om, "three")

	// Pop first
	if key, val, ok := om.Pop(false); key != "one" || val != 1 || !ok {
		t.Error("Pop first item error")
	}
	mapHasKey(t, om, "two", 2)
	mapNotKey(t, om, "one")
	mapNotKey(t, om, "three")

	// Add and Pop again
	om.Set("four", 4)
	om.Set("two", "new 2") //This should only change the value
	mapNotKey(t, om, "one")
	mapNotKey(t, om, "three")
	mapHasKey(t, om, "two", "new 2")
	mapHasKey(t, om, "four", 4)

	// Pop last
	om.Pop(false)
	mapNotKey(t, om, "one")
	mapNotKey(t, om, "two")
	mapNotKey(t, om, "three")
	mapHasKey(t, om, "four", 4)

	// Pop first
	om.Pop(true)
	mapNotKey(t, om, "one")
	mapNotKey(t, om, "two")
	mapNotKey(t, om, "three")
	mapNotKey(t, om, "four")

	// Check it is empty
	mapIsEmpty(t, om)

	// Add a last one and pop it
	om.Set("six", 6)
	mapHasKey(t, om, "six", 6)
	om.Pop(true)
	mapNotKey(t, om, "six")

	// Try to pop an item from an empty map
	mapIsEmpty(t, om)
}

func TestPopLast(t *testing.T) {
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.PopLast(); key != "two" || value != 2 || !ok {
		t.Error("PopLast didn't pop last element")
	}
}

func TestPopFirst(t *testing.T) {
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.PopFirst(); key != "one" || value != 1 || !ok {
		t.Error("PopFirst didn't pop first element")
	}
}

func TestLen(t *testing.T) {
	om := NewOrderedMap()
	if l := om.Len(); l != 0 {
		t.Error("Expecting 0, returned ", l)
	}

	om.Set("one", 1)
	if l := om.Len(); l != 1 {
		t.Error("Expecting 1, returned ", l)
	}

	om.Set("two", 2)
	if l := om.Len(); l != 2 {
		t.Error("Expecting 2, returned ", l)
	}

	om.Pop(true)
	if l := om.Len(); l != 1 {
		t.Error("Expenting 1, returned ", l)
	}
	om.Pop(true)
	if l := om.Len(); l != 0 {
		t.Error("Expecting 0, returned ", l)
	}
}

func TestDelete(t *testing.T) {

	om := NewOrderedMap()

	om.Set("one", 1)
	om.Set("two", 2)

	om.Delete("one")
	mapNotKey(t, om, "one")
	mapHasKey(t, om, "two", 2)

	om.Delete("two")
	mapNotKey(t, om, "one")
	mapNotKey(t, om, "two")

	// Add and delete from empty OrderedMap
	om.Set("three", 3)
	mapHasKey(t, om, "three", 3)

	om.Delete("three")
	mapNotKey(t, om, "three")

	mapIsEmpty(t, om)
}

func TestMove(t *testing.T) {

	// Move last
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)
	om.Set("three", 3)

	// Moving last element to the end should leave everithing uncahnged
	om.Move("three", true)

	// Move "one" element to the end
	om.Move("one", true)
	mapHasKey(t, om, "one", 1)
	mapHasKey(t, om, "two", 2)
	mapHasKey(t, om, "three", 3)
	if key, value, ok := om.Pop(true); key != "one" || value != 1 || !ok {
		t.Error("Item was not moved to the end")
	}

	// Try to move unknown key
	if ok := om.Move("unknown", true); ok {
		t.Error("Moved a non-existent element")
	}

	// Move "three" to the beginning
	om.Move("three", false)
	mapHasKey(t, om, "three", 3)
	mapHasKey(t, om, "two", 2)
	if key, value, ok := om.Pop(true); key != "two" || value != 2 || !ok {
		t.Error("Item was not moved to the start")
	}

	// Move when there is a single element
	om.Move("three", false)
	om.Move("three", true)
	mapHasKey(t, om, "three", 3)

	if key, value, ok := om.Pop(false); key != "three" || value != 3 || !ok {
		t.Error("There was an error while moving the last element")
	}

	mapIsEmpty(t, om)

	// Try to move empty map
	if ok := om.Move("three", true); ok {
		t.Error("Somehow moved a key in an empy map")
	}
}

// Just test it MoveLast calls Move with correct option
func TestMoveLast(t *testing.T) {

	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)

	om.MoveLast("one")

	if key, value, ok := om.PopLast(); key != "one" || value != 1 || !ok {
		t.Error("MoveLast didn't move to last position")
	}
}

func TestMoveFirst(t *testing.T) {
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)

	om.MoveFirst("two")

	if key, value, ok := om.PopLast(); key != "one" || value != 1 || !ok {
		t.Error("MoveFirst didn't move to the beginning")
	}

}

// Test string interface
func TestString(t *testing.T) {

	om := NewOrderedMap()

	if fmt.Sprintf("%v", om) != "OrderedMap[]" {
		t.Error("Invalid empty OrderedMap representation")
	}

	om.Set(1, 2)
	if fmt.Sprintf("%v", om) != "OrderedMap[1:2, ]" {
		t.Error("Invalid OrderedMap representation")
	}

}
