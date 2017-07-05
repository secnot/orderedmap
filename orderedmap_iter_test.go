package orderedmap

import (
	"fmt"
	"testing"
)

type KeyValue struct {
	key   int
	value int
}

type VisitorFunc func(t *testing.T, om *OrderedMap, iter_num int, key interface{}, value interface{})

//
// Visitor func
// t
// om
// in -> iteration number starting from 0
func IterVisitor(t *testing.T, om *OrderedMap, visitor VisitorFunc) {

	iterNum := 0
	iter := om.Iter()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		visitor(t, om, iterNum, k, v)
		iterNum++
	}
}

func IterReverseVisitor(t *testing.T, om *OrderedMap, visitor VisitorFunc) {

	iterNum := 0
	iter := om.IterReverse()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		visitor(t, om, iterNum, k, v)
		iterNum++
	}
}

// Helper function to test iter key:values at the same time a function is applied
// ifunc -> function applied during each iteration
// values -> array of values expected during iteration
// revers -> iterate in reverse
func ApplyIterFunc(t *testing.T, om *OrderedMap, ifunc VisitorFunc, values []KeyValue, reverse bool) {

	iterFunc := func(t *testing.T, om *OrderedMap, iter_num int, key interface{}, value interface{}) {

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
		IterReverseVisitor(t, om, iterFunc)
	} else {
		IterVisitor(t, om, iterFunc)
	}
}

func TestIterator(t *testing.T) {

	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{5, 5},
	}

	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	// Standard Iteration
	current := 0
	iter := om.Iter()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		if tests[v.(int)].key != current {
			fmt.Printf("%v %v", k, v)
			t.Error("Iteration error", k, v)
		}
		current++
	}

	// Reverse Iteration
	iter = om.IterReverse()
	start := len(tests) - 1
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

	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
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
		if k.(int)%2 == 0 {
			om.Set(k, 100)
		}
	}

	result1 := []KeyValue{
		{0, 100},
		{1, 1},
		{2, 100},
		{3, 3},
		{4, 100},
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
		if k.(int)%2 == 0 && k.(int) < len(tests)-1 {
			om.Set(k.(int)+1, 101)
		}
	}

	result2 := []KeyValue{
		{0, 0},
		{1, 101},
		{2, 2},
		{3, 101},
		{4, 4},
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
		if k.(int)%2 == 0 && k.(int) > 0 {
			om.Set(k.(int)-1, 103)
		}
	}

	result3 := []KeyValue{
		{0, 0},
		{1, 103},
		{2, 2},
		{3, 103},
		{4, 4},
	}

	ApplyIterFunc(t, om, nil, result3, false)
}

func TestIterReverseSet(t *testing.T) {

	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
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
		if k.(int)%2 == 0 {
			om.Set(k, 100)
		}
	}

	result1 := []KeyValue{
		{0, 100},
		{1, 1},
		{2, 100},
		{3, 3},
		{4, 100},
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
		if k.(int)%2 == 0 && k.(int) < len(tests)-1 {
			om.Set(k.(int)+1, 101)
		}
	}

	result2 := []KeyValue{
		{0, 0},
		{1, 101},
		{2, 2},
		{3, 101},
		{4, 4},
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
		if k.(int)%2 == 0 && k.(int) > 0 {
			om.Set(k.(int)-1, 103)
		}
	}

	result3 := []KeyValue{
		{0, 0},
		{1, 103},
		{2, 2},
		{3, 103},
		{4, 4},
	}

	ApplyIterFunc(t, om, nil, result3, false)
}

// Test Adding new keys while iterating
func TestIterInsert(t *testing.T) {

	// Can insert keys while iterating and will be included in the
	// current iteration, unless the last one has already been reached
	// or it is using IterReverse.
	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
	}

	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iterCount := 0 // Count iterations to check added/keys are included
	iter := om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int) < len(tests) {
			om.Set(k.(int)+100, 100)
		}
		iterCount++
	}

	if iterCount != 2*len(tests) || iterCount != om.Len() {
		t.Error("Failed adding new keys")
	}

	result := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{100, 100},
		{101, 100},
		{102, 100},
		{103, 100},
		{104, 100},
	}

	ApplyIterFunc(t, om, nil, result, false)
}

func TestIterReverseInsert(t *testing.T) {

	// Can insert keys while iterating but will not appear in the
	// current iteration
	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
	}

	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	iterCount := 0 // Count iterations to check added/keys are included
	iter := om.IterReverse()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int) < len(tests) {
			om.Set(k.(int)+100, 100)
		}
		iterCount++
	}

	if iterCount != len(tests) || om.Len() != 2*len(tests) {
		t.Error("Failed adding new keys")
	}

	result := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{104, 100},
		{103, 100},
		{102, 100},
		{101, 100},
		{100, 100},
	}

	ApplyIterFunc(t, om, nil, result, false)
}

// Test deleting current key while iterating
func TestIterDeleteCurrent(t *testing.T) {
	var tests []KeyValue

	for k := 0; k < 100; k++ {
		tests = append(tests, KeyValue{key: k, value: k})
	}

	for num := range tests {

		deletef := func(t *testing.T, om *OrderedMap, iter_num int,
			key interface{}, value interface{}) {
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
		var expectedResult []KeyValue
		expectedResult = append(expectedResult, tests[:num]...)
		expectedResult = append(expectedResult, tests[num+1:]...)
		ApplyIterFunc(t, om, nil, expectedResult, false)
	}
}

//Test deleting all map keys while iterating (only current)
func TestIterDeleteAll(t *testing.T) {

	tests := []KeyValue{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{5, 5},
		{6, 6},
	}

	// Initialize map with al test elements
	om := NewOrderedMap()
	for _, test := range tests {
		om.Set(test.key, test.value)
	}

	// Delete even keys
	iter := om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2 == 0 {
			om.Delete(k)
		}
	}

	if om.Len() != 3 {
		t.Error("Failed deleting keys")
	}

	// Check odd keys where not deleted
	result := []KeyValue{
		{1, 1},
		{3, 3},
		{5, 5},
	}

	ApplyIterFunc(t, om, nil, result, false)

	// Delete odd keys
	iter = om.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if k.(int)%2 == 1 {
			om.Delete(k)
		}
	}

	if om.Len() != 0 {
		t.Error("Failed deleting keys")
	}
}

// Delete random keys while iterating (only one per iteration)
func TestIterDeleteOneRandom(t *testing.T) {
	var tests []KeyValue
	var reversedTests []KeyValue

	nkeys := 10

	for k := 0; k < nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
		reversedTests = append(reversedTests, KeyValue{nkeys - 1 - k, nkeys - 1 - k})
	}

	// Forward Iteration
	deleteKey, pos := 0, 0

	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			// delete key at delete_pos
			deleteFunc := func(t *testing.T, om *OrderedMap, iterNum int,
				key interface{}, value interface{}) {
				if iterNum == pos {
					om.Delete(deleteKey)
				}
			}

			// Verify expected key/value pairs during iterarion
			var iterKeys []KeyValue
			if deleteKey <= pos { // Delete key already visited
				iterKeys = append(iterKeys, tests[:]...)
			} else { // Delete key not yet visited
				iterKeys = append(iterKeys, tests[:deleteKey]...)
				iterKeys = append(iterKeys, tests[deleteKey+1:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iterKeys, false)

			// Verify resulting map
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, tests[:deleteKey]...)
			expectedResult = append(expectedResult, tests[deleteKey+1:]...)
			ApplyIterFunc(t, om, nil, expectedResult, false)
		}
	}

	// Reverse Iteration
	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iterNum int,
				key interface{}, value interface{}) {
				if key.(int) == pos {
					om.Delete(deleteKey)
				}
			}

			// Verify expected key/value pairs during iterarion
			var iterKeys []KeyValue
			if deleteKey < pos { // Delete key not yet visited
				iterKeys = append(iterKeys, reversedTests[:nkeys-1-deleteKey]...)
				iterKeys = append(iterKeys, reversedTests[nkeys-1-deleteKey+1:]...)
			} else { // Deleted visited key
				iterKeys = append(iterKeys, reversedTests[:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iterKeys, true)

			// verify resulting map
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, reversedTests[:nkeys-1-deleteKey]...)
			expectedResult = append(expectedResult, reversedTests[nkeys-1-deleteKey+1:]...)
			ApplyIterFunc(t, om, nil, expectedResult, true)
		}
	}
}

// Set a new key then delete another (while iterating)
func TestIterInsertDelete(t *testing.T) {
	var tests []KeyValue         // Tests for normal Iter
	var reversedTests []KeyValue // Tests for IterReverse
	nkeys := 10

	for k := 0; k < nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
		reversedTests = append(reversedTests, KeyValue{nkeys - 1 - k, nkeys - 1 - k})
	}

	pos, deleteKey := 0, 0

	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
				key interface{}, value interface{}) {
				if iter_num == pos {
					om.Set(pos+1000, pos+1000)
					om.Delete(deleteKey)
				}
			}

			// Because it is iterating from start to finish, it will iterate
			// over the newly added keys
			var iterKeys []KeyValue
			if deleteKey <= pos {
				iterKeys = append(iterKeys, tests[:]...)
			} else {
				iterKeys = append(iterKeys, tests[:deleteKey]...)
				iterKeys = append(iterKeys, tests[deleteKey+1:]...)
			}

			iterKeys = append(iterKeys, KeyValue{pos + 1000, pos + 1000})

			// Check all tests are iterated
			ApplyIterFunc(t, om, deleteFunc, iterKeys, false)

			// Check everything went ok
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, tests[:deleteKey]...)
			expectedResult = append(expectedResult, tests[deleteKey+1:]...)
			expectedResult = append(expectedResult, KeyValue{pos + 1000, pos + 1000})
			ApplyIterFunc(t, om, nil, expectedResult, false)
		}
	}

	// Test Reverse Iteration
	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
				key interface{}, value interface{}) {
				if key.(int) == pos {
					om.Set(pos+1000, pos+1000)
					om.Delete(deleteKey)
				}
			}
			// The inserted keys should not appear during reverse iteration because
			// they are added at the end of the OrderedMap and it's iterating towards
			// the start
			var iterKeys []KeyValue
			if deleteKey < pos { // Delete key not yet visited
				iterKeys = append(iterKeys, reversedTests[:nkeys-1-deleteKey]...)
				iterKeys = append(iterKeys, reversedTests[nkeys-1-deleteKey+1:]...)
			} else { // Deleted visited key
				iterKeys = append(iterKeys, reversedTests[:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iterKeys, true)

			// The new keys should appear in the next iteration
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, KeyValue{pos + 1000, pos + 1000})
			expectedResult = append(expectedResult, reversedTests[:nkeys-1-deleteKey]...)
			expectedResult = append(expectedResult, reversedTests[nkeys-1-deleteKey+1:]...)
			ApplyIterFunc(t, om, nil, expectedResult, true)
		}
	}

}

// Delete key and then insert a new one (while iterating)
func TestIterDelteInsert(t *testing.T) {
	var tests []KeyValue
	var reversedTests []KeyValue
	nkeys := 10

	for k := 0; k < nkeys; k++ {
		tests = append(tests, KeyValue{k, k})
		reversedTests = append(reversedTests, KeyValue{nkeys - 1 - k, nkeys - 1 - k})
	}

	pos, deleteKey := 0, 0

	// Forward
	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
				key interface{}, value interface{}) {
				if iter_num == pos {
					om.Delete(deleteKey)
					om.Set(pos+1000, pos+1000)
				}
			}

			var iterKeys []KeyValue
			if deleteKey <= pos {
				iterKeys = append(iterKeys, tests[:]...)
			} else {
				iterKeys = append(iterKeys, tests[:deleteKey]...)
				iterKeys = append(iterKeys, tests[deleteKey+1:]...)
			}

			// If the iteration is at its last position and it is removed
			// the iteration is finished, so the key added after is ignored
			if !(pos == nkeys-1 && deleteKey == pos) {
				iterKeys = append(iterKeys, KeyValue{pos + 1000, pos + 1000})
			}

			// Check all tests are iterated
			ApplyIterFunc(t, om, deleteFunc, iterKeys, false)

			// Check everything went ok
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, tests[:deleteKey]...)
			expectedResult = append(expectedResult, tests[deleteKey+1:]...)
			expectedResult = append(expectedResult, KeyValue{pos + 1000, pos + 1000})
			ApplyIterFunc(t, om, nil, expectedResult, false)
		}
	}

	// Test Reverse Iteration
	for pos = 0; pos < nkeys; pos++ {
		for deleteKey = 0; deleteKey < nkeys; deleteKey++ {
			// Initialize map with al test elements
			om := NewOrderedMap()
			for _, test := range tests {
				om.Set(test.key, test.value)
			}

			deleteFunc := func(t *testing.T, om *OrderedMap, iter_num int,
				key interface{}, value interface{}) {
				if key.(int) == pos {
					om.Delete(deleteKey)
					om.Set(pos+1000, pos+1000)
				}
			}
			// The inserted keys should not appear during reverse iteration because
			// they are added at the end of the OrderedMap and it's iterating towards
			// the start
			var iterKeys []KeyValue
			if deleteKey < pos { // Delete key not yet visited
				iterKeys = append(iterKeys, reversedTests[:nkeys-1-deleteKey]...)
				iterKeys = append(iterKeys, reversedTests[nkeys-1-deleteKey+1:]...)
			} else { // Deleted visited key
				iterKeys = append(iterKeys, reversedTests[:]...)
			}
			ApplyIterFunc(t, om, deleteFunc, iterKeys, true)

			// The new keys should appear in the next iteration
			var expectedResult []KeyValue
			expectedResult = append(expectedResult, KeyValue{pos + 1000, pos + 1000})
			expectedResult = append(expectedResult, reversedTests[:nkeys-1-deleteKey]...)
			expectedResult = append(expectedResult, reversedTests[nkeys-1-deleteKey+1:]...)
			ApplyIterFunc(t, om, nil, expectedResult, true)
		}
	}
}
