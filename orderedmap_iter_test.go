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


func TestIterator(t *testing.T) {

	tests := []KeyValue {
		KeyValue {0, 0},
		KeyValue {1, 1},
		KeyValue {2, 2},
		KeyValue {3, 3},
		KeyValue {4, 4},
		KeyValue {5, 5},
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


