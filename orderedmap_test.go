package orderedmap

import (
	"testing"
	"fmt"
)



type KeyValue struct {
	key int
	value int
}


type VisitorFunc func(t *testing.T, om *OrderedMap, iter_num int, key interface{}, value interface{})

//
// Visitor func
// t 
// om
// in -> iteration number starting from 0
func IterVisitor(t *testing.T, om *OrderedMap, visitor VisitorFunc) {
	
	iter_num := 0
	iter := om.Iter()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		visitor(t, om, iter_num, k, v)
		iter_num ++
	}
}

func IterReverseVisitor(t *testing.T, om *OrderedMap, visitor VisitorFunc) {
	
	iter_num := 0
	iter := om.IterReverse()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		visitor(t, om, iter_num, k, v)
		iter_num ++
	}
}




// Helper function to test iter key:values at the same time a function is applied
// ifunc -> function applied during each iteration
// values -> array of values expected during iteration
// revers -> iterate in reverse
func ApplyIterFunc(t *testing.T, om *OrderedMap, ifunc VisitorFunc, values []KeyValue, reverse bool) {
	
	iter_func := func(t *testing.T, om *OrderedMap, iter_num int, key interface{}, value interface{}) {
	
		if iter_num >= len(values) {
			t.Error("Iteration too long")
		}

		if values[iter_num].key != key || values[iter_num].value != value {
			t.Error(fmt.Sprintf("Expecting {key: %v, value: %v}",
								values[iter_num].key, values[iter_num].value))
			t.Error(fmt.Sprintf("Received {key: %v, value: %v}", key, value)) 
		}

		if ifunc != nil {
			ifunc(t, om, iter_num, key, value)
		}
	}
	if reverse {
		IterReverseVisitor(t, om, iter_func)
	} else {
		IterVisitor(t, om, iter_func)
	}
}


func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}




func TestNewOrderedMap (t* testing.T) {
	om := NewOrderedMap()
	om.Set(5, 33)
	om.Set(6, 44)
	om.Set("key", "value")
	
	if val, ok := om.Get(5); val != 33 || !ok {
		t.Error("Value error, expecting 33 received %d", val)
	}
	if val, ok := om.Get(6); val != 44 || !ok {
		t.Error("Value error, expecting 44 received %d", val)
	}
	if val, ok := om.Get("key"); val != "value" || !ok {
		t.Error("Value error, expecting 'value' recived %s", val)
	}

	if val, ok := om.Get("not a key"); ok || val != nil {
		t.Error("Shouldn't have returned %v", val)
	}
}


func TestGetLast(t *testing.T) {
	om := NewOrderedMap()

	if key, value, ok := om.GetLast(); key != nil || value!=nil || ok{
		t.Error(fmt.Sprintf("Expecting nil, nil, false -> Returned %v %v %v", 
							key, value, ok,))
	}

	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.GetLast(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
	}
	
	if key, value, ok := om.GetLast(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
	}

	om.MoveLast("one")
	if key, value, ok := om.GetLast(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
	}

	if om.Len() != 2 {
		t.Error("Somehow popped an item ..")
	}
}


func TestGetFirst(t *testing.T) {
	om := NewOrderedMap()

	if key, value, ok := om.GetFirst(); key != nil || value!=nil || ok{
		t.Error(fmt.Sprintf("Expecting nil, nil, false -> Returned %v %v %v", 
							key, value, ok,))
	}

	om.Set("one", 1)
	om.Set("two", 2)

	if key, value, ok := om.GetFirst(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
	}
	
	if key, value, ok := om.GetFirst(); key != "one" || value != 1 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
	}

	om.MoveLast("one")
	if key, value, ok := om.GetFirst(); key != "two" || value != 2 || !ok {
		t.Error(fmt.Sprintf("Expecting 'two', 2, true -> Returned %v %v %v", 
							key, value, ok,))
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


	// Test key present in OrderedMap
	test_has_key := func (om *OrderedMap, key interface{}, value interface{}) bool{
		
		if v, ok := om.Get(key); v!=value || !ok {
			t.Error(fmt.Sprintf("Get(%v) -> expected %v received %v", key, value, v))
			return false
		} else {
			return true
		}
	}

	// Test key not present in OrderedMap
	test_not_key := func (om *OrderedMap, key interface{}) bool {

		if v, ok:= om.Get(key); v!=nil || ok {
			t.Error(fmt.Sprintf("Get(%v) -> shouldn't have a value", key))
			return true
		} else {
			return false
		}
	}

	// Pop last
	if key, val, ok := om.Pop(true); key!="three" || val!=3 || !ok {
		t.Error("Pop last item error")
	}
	test_has_key(om, "one", 1)
	test_has_key(om, "two", 2)
	test_not_key(om, "three")

	// Pop first
	if key, val, ok := om.Pop(false); key!="one" || val!=1 || !ok {
		t.Error("Pop first item error")
	}
	test_has_key(om, "two", 2)
	test_not_key(om, "one")
	test_not_key(om, "three")

	// Add and Pop again
	om.Set("four", 4)
	om.Set("five", 5)
	om.Set("two", "new 2") //This should only change the value
	test_not_key(om, "one")
	test_not_key(om, "three")
	test_has_key(om, "two", "new 2")
	test_has_key(om, "four", 4)
	test_has_key(om, "five", 5)
	
	// pop first
	if key, val, ok := om.Pop(false); key != "two" || val!= "new 2" || !ok {
		t.Error("Popped ")
	}
	test_not_key(om, "one")
	test_not_key(om, "two")
	test_not_key(om, "three")
	test_has_key(om, "four", 4)
	test_has_key(om, "five", 5)

	if key, val, ok := om.Pop(true); key != "five" || val!= 5 || !ok {
		t.Error("Pop Error ")
	}
	test_not_key(om, "one")
	test_not_key(om, "two")
	test_not_key(om, "three")
	test_not_key(om, "five")
	test_has_key(om, "four", 4)
	
	if key, val, ok := om.Pop(true); key != "four" || val!= 4 || !ok {
		t.Error("%v %v %v", key, val, ok)
		t.Error("Pop Error ")
	}
	test_not_key(om, "one")
	test_not_key(om, "two")
	test_not_key(om, "three")
	test_not_key(om, "four")
	test_not_key(om, "five")

	// Check it is empty
	if key, val, ok := om.Pop(false); key != nil || val != nil || ok {
		t.Error("Pop should be empty")
	}

	// Add a last one and pop it
	om.Set("six", 6)
	test_has_key(om, "six", 6)
	if key, val, ok := om.Pop(true); key != "six" || val!= 6 || !ok {
		t.Error("Pop Error ")
	}
	test_not_key(om, "six")

	// Try to pop an item from an empty amp
	if key, val, ok := om.Pop(true); key != nil || val != nil || ok{
		t.Error("Map should be empty")
	}
}


func TestPopLast(t *testing.T) {	
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)
	
	if key, value, ok := om.PopLast(); key!= "two" || value!=2||!ok {
		t.Error("PopLast didn't pop last element")
	}
}


func TestPopFirst(t *testing.T) {
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)
	
	if key, value, ok := om.PopFirst(); key!= "one" || value!=1||!ok {
		t.Error("PopFirst didn't pop first element")
	}
}


func TestLen(t* testing.T) {
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
		t.Error("Expection 0, returned ", l)
	}
}


func TestDelete(t *testing.T) {
	
	// Test key present in OrderedMap
	test_has_key := func (om *OrderedMap, key interface{}, value interface{}) bool{
		
		if v, ok := om.Get(key); v!=value || !ok {
			t.Error(fmt.Sprintf("Get(%v) -> expected %v received %v", key, value, v))
			return false
		} else {
			return true
		}
	}

	// Test key not present in OrderedMap
	test_not_key := func (om *OrderedMap, key interface{}) bool {

		if v, ok:= om.Get(key); v!=nil || ok {
			t.Error(fmt.Sprintf("Get(%v) -> shouldn't have a value", key))
			return true
		} else {
			return false
		}
	}

	om := NewOrderedMap()
	
	om.Set("one", 1)
	om.Set("two", 2)
	
	om.Delete("one")
	test_not_key(om, "one")
	test_has_key(om, "two", 2)

	om.Delete("two")
	test_not_key(om, "one")
	test_not_key(om, "two")
	
	// Add and delete from empty OrderedMap
	om.Set("three", 3)
	test_has_key(om, "three", 3)

	om.Delete("three")
	test_not_key(om, "three")

	if _, _, ok := om.Pop(true); ok {
		t.Error("Map should be empty")
	}
}


func TestMove(t *testing.T) {	
	
	// Test key present in OrderedMap
	test_has_key := func (om *OrderedMap, key interface{}, value interface{}) {
		if v, ok := om.Get(key); v!=value || !ok {
			t.Error(fmt.Sprintf("Get(%v) -> expected %v received %v", key, value, v))
		}
	}

	// Move last
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)
	om.Set("three", 3)

	// Moving last element to the end should leave everithing uncahnged
	om.Move("three", true)
	
	// Move "one" element to the end
	om.Move("one", true)
	test_has_key(om, "one", 1)
	test_has_key(om, "two", 2)
	test_has_key(om, "three", 3)
	if key, value, ok := om.Pop(true); key!="one" || value != 1 || !ok {
		t.Error("Item was not moved to the end")
	}

	// Try to move unknown key
	if ok := om.Move("unknown", true); ok {
		t.Error("Moved a non-existent element")
	}

	// Move "three" to the beginning
	om.Move("three", false)
	test_has_key(om, "three", 3)
	test_has_key(om, "two", 2)
	if key, value, ok := om.Pop(true); key!="two" || value != 2 || !ok {
		t.Error("Item was not moved to the start")
	}

	// Move when there is a single element
	om.Move("three", false)
	om.Move("three", true)
	test_has_key(om, "three", 3)

	if key, value, ok := om.Pop(false); key!= "three" || value != 3 || !ok {
		t.Error("There was an error while moving the last element")
	}

	if om.Len() != 0 {
		t.Error("The Map should have been empty")
	}

	// Tru to move empty map
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

	if key, value, ok := om.PopLast(); key!= "one" || value!=1||!ok {
		t.Error("MoveLast didn't move to last position")
	}
}

func TestMoveFirst(t *testing.T) {
	om := NewOrderedMap()
	om.Set("one", 1)
	om.Set("two", 2)
	
	om.MoveFirst("two")

	if key, value, ok := om.PopLast(); key!= "one" || value!=1||!ok {
		t.Error("MoveFirst didn't move to the beginning")
	}

}



func TestIterator(t *testing.T) {

	type TestCase struct {
		key interface {}
		value interface {}
	}

	tests := []TestCase {
		TestCase {"zero", 0},
		TestCase {"one", 1},
		TestCase {"two", 2},
		TestCase {"three", 3},
		TestCase {"four", 4},
		TestCase {"five", 5},
	}


	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	// Standard Iteration
	iter := om.Iter()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		if tests[v.(int)].key != k {
			fmt.Printf("%v %v", k, v)
			t.Error("Iteration error", k, v)
		}
	}

	// Reverse Iteration
	iter = om.IterReverse()
	start := len(tests)-1
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		if tests[v.(int)].key != k || start != v {
			t.Error("Iteration error", k, v)
		}
		start--
	}

	// Iterating empty map
	om = NewOrderedMap()
	iterations := 0 
	iter = om.Iter()
	for _, _, ok := iter.Next(); ok; _, _, ok = iter.Next() {
		iterations++
	}
	iter = om.IterReverse()
	for _, _, ok := iter.Next(); ok; _, _, ok = iter.Next() {
		iterations++
	}
	if iterations != 0 {
		t.Error("Iterated an empty map")
	}

	// Try to use Next() one more time after the iteration has finished
	om = NewOrderedMap()
	om.Set(1, 2)
	iter = om.Iter()
	for _, _, ok := iter.Next(); ok; _, _, ok = iter.Next() {
	}

	if _, _, ok := iter.Next(); ok {
		t.Error("Next() returned a value after iterationg was finished")
	}
	
}



// Test Setting existing keys values while iterating
func TestIterSet(t *testing.T) {
	
	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
	}

	// Modify current key while iterating
	// ----------------------------------
	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter := om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		// Set odd values 100
		if k.(int)%2  == 0 {
			om.Set(k, 100)
		}
	}
	
	result1 := []KeyValue {
		KeyValue {0, 100},
		KeyValue {1, 1},
		KeyValue {2, 100},
		KeyValue {3, 3},
		KeyValue {4, 100},
	}
	
	ApplyIterFunc(t, om, nil, result1, false)

	// Modify next key while iterating
	// -------------------------------
	om = NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter = om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2  == 0 && k.(int) < len(tests)-1{
			om.Set(k.(int)+1, 101)
		}
	}
	
	result2 := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 101},
		KeyValue {2, 2},
		KeyValue {3, 101},
		KeyValue {4, 4},
	}
	
	ApplyIterFunc(t, om, nil, result2, false)

	// Modify previous key while iterating
	// -----------------------------------
	om = NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter = om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2  == 0 && k.(int) > 0{
			om.Set(k.(int)-1, 103)
		}
	}
	
	result3 := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 103},
		KeyValue {2, 2},
		KeyValue {3, 103},
		KeyValue {4, 4},
	}
	
	ApplyIterFunc(t, om, nil, result3, false)
}


func TestIterReverseSet(t *testing.T) {
	
	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
	}

	// Modify current key while iterating
	// ----------------------------------
	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter := om.IterReverse()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		// Set odd values 100
		if k.(int)%2  == 0 {
			om.Set(k, 100)
		}
	}
	
	result1 := []KeyValue {
		KeyValue {0, 100},
		KeyValue {1, 1},
		KeyValue {2, 100},
		KeyValue {3, 3},
		KeyValue {4, 100},
	}
	
	ApplyIterFunc(t, om, nil, result1, false)

	// Modify next key while iterating
	// -------------------------------
	om = NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter = om.IterReverse()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2  == 0 && k.(int) < len(tests)-1{
			om.Set(k.(int)+1, 101)
		}
	}
	
	result2 := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 101},
		KeyValue {2, 2},
		KeyValue {3, 101},
		KeyValue {4, 4},
	}
	
	ApplyIterFunc(t, om, nil, result2, false)

	// Modify previous key while iterating
	// -----------------------------------
	om = NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter = om.IterReverse()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2  == 0 && k.(int) > 0{
			om.Set(k.(int)-1, 103)
		}
	}
	
	result3 := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 103},
		KeyValue {2, 2},
		KeyValue {3, 103},
		KeyValue {4, 4},
	}
	
	ApplyIterFunc(t, om, nil, result3, false)
}




// Test Adding new keys while iterating
func TestIterInsert(t *testing.T) {

	// Can insert keys while iterating and will be included in the
	// current iteration, unless the last one has already been reached
	// or it is using IterReverse.
	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
	}

	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter_count := 0 // Count iterations to check added/keys are included
	iter := om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)<len(tests) {
			om.Set(k.(int)+100, 100)
		}
		iter_count++
	}

	if iter_count!=2*len(tests) || iter_count!=om.Len(){
		t.Error("Failed adding new keys")
	}
	
	result := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
		KeyValue {100, 100},
		KeyValue {101, 100},
		KeyValue {102, 100},
		KeyValue {103, 100},
		KeyValue {104, 100},
	}

	ApplyIterFunc(t, om, nil, result, false)
}


func TestIterReverseInsert(t *testing.T) {

	// Can insert keys while iterating but will not appear in the 
	// current iteration	
	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
	}

	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iter_count := 0 // Count iterations to check added/keys are included
	iter := om.IterReverse()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)<len(tests) {
			om.Set(k.(int)+100, 100)
		}
		iter_count++
	}

	if iter_count!=len(tests) || om.Len() != 2 * len(tests){
		t.Error("Failed adding new keys")
	}
	
	result := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
		KeyValue {104, 100},
		KeyValue {103, 100},
		KeyValue {102, 100},
		KeyValue {101, 100},
		KeyValue {100, 100},
	}

	ApplyIterFunc(t, om, nil, result, false)
}


// Test deleting current key while iterating
func TestIterDeleteCurrent(t *testing.T) {
	var tests []KeyValue 
	
	for k := 0; k<100; k++ {
		tests = append(tests, KeyValue{key:k, value:k})
	}

	for num, _ := range tests {
	
		deletef := func(t *testing.T, om *OrderedMap, iter_num int,
						key interface{}, value interface{}){
				if iter_num == num {
					om.Delete(num)
				}
		}
		
		// Initialize map with al test elements
		om := NewOrderedMap()
		for _, test := range tests {
			om.Set(test.key, test.value)
		}
		
		// Check iters through all elements 
		ApplyIterFunc(t, om, deletef, tests, false)

		// Check the element was deleted
		var expected_result []KeyValue
		expected_result = append(expected_result, tests[:num]...)
		expected_result = append(expected_result, tests[num+1:]...)
		ApplyIterFunc(t, om, nil, expected_result, false)
	}
}


//Test deleting all map keys while iterating (only current)
func TestIterDeleteAll(t *testing.T) {
	
	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
		KeyValue {5, 5},
		KeyValue {6, 6},
	}
	
	// Initialize map with al test elements
	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	// Delete even keys
	iter := om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2==0{
			om.Delete(k)
		}
	}

	if om.Len() != 3 {
		t.Error("Failed deleting keys")
	}

	// Check odd keys where not deleted
	result := []KeyValue {
		KeyValue {1, 1},
		KeyValue {3, 3},
		KeyValue {5, 5},
	}

	ApplyIterFunc(t, om, nil, result, false)

	// Delete odd keys
	iter = om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2==1{
			om.Delete(k)
		}
	}

	if om.Len() != 0 {
		t.Error("Failed deleting keys")
	}
}


// Delete random keys while iterating (only one per iteration)
func TestIterDeleteOneRandom(t *testing.T){
	var tests []KeyValue
	var reversed_tests []KeyValue

	nkeys := 10


	for k := 0; k<nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
		reversed_tests = append(reversed_tests, KeyValue{nkeys-1-k, nkeys-1-k})
	}

	// Forward Iteration
	delete_key, pos := 0, 0

	for pos = 0; pos < nkeys; pos++ {
		for delete_key = 0; delete_key < nkeys; delete_key++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			// delete key at delete_pos
			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
							key interface{}, value interface{}){
					if iter_num == pos {
						om.Delete(delete_key)
					}
			}
		
			// Verify expected key/value pairs during iterarion
			var iter_keys []KeyValue
			if delete_key <= pos { // Delete key already visited
				iter_keys = append(iter_keys, tests[:]...)
			} else { // Delete key not yet visited
				iter_keys = append(iter_keys, tests[:delete_key]...)
				iter_keys = append(iter_keys, tests[delete_key+1:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iter_keys, false)
	
			// Verify resulting map
			var expected_result []KeyValue
			expected_result = append(expected_result, tests[:delete_key]...)
			expected_result = append(expected_result, tests[delete_key+1:]...)
			ApplyIterFunc(t, om, nil, expected_result, false)
		}
	}
	
	// Reverse Iteration
	for pos = 0; pos < nkeys; pos++ {
		for delete_key = 0; delete_key < nkeys; delete_key++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
							key interface{}, value interface{}){
					if key.(int) == pos {
						om.Delete(delete_key)
					}
			}
			
			// Verify expected key/value pairs during iterarion
			var iter_keys []KeyValue
			if delete_key < pos { // Delete key not yet visited
				iter_keys = append(iter_keys, reversed_tests[:nkeys-1-delete_key]...)
				iter_keys = append(iter_keys, reversed_tests[nkeys-1-delete_key+1:]...)
			} else { // Deleted visited key
				iter_keys = append(iter_keys, reversed_tests[:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iter_keys, true)

			// verify resulting map
			var expected_result []KeyValue
			expected_result = append(expected_result, reversed_tests[:nkeys-1-delete_key]...)
			expected_result = append(expected_result, reversed_tests[nkeys-1-delete_key+1:]...)
			ApplyIterFunc(t, om, nil, expected_result, true)
		}
	}
}


// Set a new key then delete another (while iterating)
func TestIterInsertDelete(t *testing.T){
	var tests []KeyValue 
	nkeys := 4


	for k := 0; k<nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
	}

	pos, delete_key := 0, 0

	// Also insert a key before deletion
	for pos = 0; pos < nkeys; pos++ {
		for delete_key = 0; delete_key < nkeys; delete_key++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}
			
			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
							key interface{}, value interface{}){
					if iter_num == pos {
						om.Set(pos+1000, pos+1000)
						om.Delete(delete_key)
					}
			}
				
			var iter_keys []KeyValue
			if delete_key <= pos {
				iter_keys = append(iter_keys, tests[:]...)
			} else {
				iter_keys = append(iter_keys, tests[:delete_key]...)
				iter_keys = append(iter_keys, tests[delete_key+1:]...)
			}
		
			iter_keys = append(iter_keys, KeyValue{pos+1000, pos+1000})

			// Check all tests are iterated
			ApplyIterFunc(t, om, deleteFunc, iter_keys, false)
	
			// Check everything went ok
			var expected_result []KeyValue
			expected_result = append(expected_result, tests[:delete_key]...)
			expected_result = append(expected_result, tests[delete_key+1:]...)
			expected_result = append(expected_result, KeyValue{pos+1000, pos+1000})
			ApplyIterFunc(t, om, nil, expected_result, false)
		}
	}

	//TODO: Test Reverse Iteration
}


// Delete key and then insert a new one (while iterating)
func TestIterDelteInsert(t *testing.T){
	var tests []KeyValue 
	nkeys := 4


	for k := 0; k<nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
	}

	pos, delete_key := 0, 0

	// Forward
	for pos = 0; pos < nkeys; pos++ {
		for delete_key = 0; delete_key < nkeys; delete_key++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}
			
			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
							key interface{}, value interface{}){
					if iter_num == pos {
						om.Delete(delete_key)
						om.Set(pos+1000, pos+1000)
					}
			}
				
			var iter_keys []KeyValue
			if delete_key <= pos {
				iter_keys = append(iter_keys, tests[:]...)
			} else {
				iter_keys = append(iter_keys, tests[:delete_key]...)
				iter_keys = append(iter_keys, tests[delete_key+1:]...)
			}
	
			// If the iteration is at its last position and it is removed
			// the iteration is finished, so the key added after is ignored
			if !(pos == nkeys - 1 && delete_key == pos) {
				iter_keys = append(iter_keys, KeyValue{pos+1000, pos+1000})
			}

			// Check all tests are iterated
			ApplyIterFunc(t, om, deleteFunc, iter_keys, false)
	
			// Check everything went ok
			var expected_result []KeyValue
			expected_result = append(expected_result, tests[:delete_key]...)
			expected_result = append(expected_result, tests[delete_key+1:]...)
			expected_result = append(expected_result, KeyValue{pos+1000, pos+1000})
			ApplyIterFunc(t, om, nil, expected_result, false)
		}
	}

	//TODO: Test reverse iteration
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

