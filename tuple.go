/*
Features:
 - support for python-style negative indexes
 - python-style comparison
 - helper functions like popleft, popright, reverse

TODO:
// standardize terminology "items" vs. "elems"
// rename to golang-list? it's more like a python list than python tuple since
//   it's mutable
// TODO:
// func (this *Tuple) Zip(...*Tuple) Tuple {}
// func (this *Tuple) Map(func) ? or Apply?
// Dict() convert to Dict which I'm also working on
// append - or better pushleft, pushright
// insert/remove
// flatten?
// group?
// chunk?
// coalesce
// append, extend

*/
package tuple

import (
	"fmt"
	_ "math"
	"reflect"
)

type Tuple struct {
	data []interface{}
}

func max(lhs int, rhs int) int {
	if lhs > rhs {
		return lhs
	} else {
		return rhs
	}
}

func min(lhs int, rhs int) int {
	if lhs < rhs {
		return lhs
	} else {
		return rhs
	}
}

// Creates a new empty Tuple of length n
func NewTuple(n int) *Tuple {
	t := &Tuple{}
	t.data = make([]interface{}, n)
	return t
}

// Creates a new Tuple from an existing slice
func NewTupleFromSlice(slice []interface{}) *Tuple {
	t := &Tuple{}
	t.data = slice
	return t
}

// Creates a new tuple from a literal sequence of items
func NewTupleFromItems(items ...interface{}) *Tuple {
	t := NewTuple(len(items))
	for i, item := range items {
		t.Set(i, item)
	}
	return t
}

// Returns a new tuple with a copy of this tuple's data
func (this *Tuple) Copy() *Tuple {
	t := NewTuple(this.Len())
	copy(t.data, this.data)
	return t
}

// Returns the number of elements in the Tuple
func (this *Tuple) Len() int {
	return len(this.data)
}

// Returns the internal slice
func (this *Tuple) Data() []interface{} {
	return this.data
}

// Returns a new tuple with copy of n leftmost items
func (this *Tuple) Left(n int) *Tuple {
	return this.Slice(0, n)
}

// Returns a new tuple with copy of n rightmost items
func (this *Tuple) Right(n int) *Tuple {
	length := this.Len()
	n = max(0, length-n)
	return this.Slice(n, length)
}

// Returns a new tuple with slice of this tuple's data
func (this *Tuple) Slice(start int, end int) *Tuple {
	c := this.Copy()
	max := this.Len()
	start = min(c.Offset(start), max)
	end = min(c.Offset(end), max)
	c.data = c.data[start:end]
	return c
}

// Convert n to an index into the internal slice.
// Negative numbers are supported, e.g. -1 points to the last item
func (this *Tuple) Offset(n int) int {
	// allow negative indexing as in Python
	if n < 0 {
		n = this.Len() + n
	}
	return n
}

// Set the item at index n
func (this *Tuple) Set(n int, item interface{}) {
	this.data[this.Offset(n)] = item
}

// Get the item at index n
func (this *Tuple) Get(n int) interface{} {
	item := this.data[this.Offset(n)]
	return item
}

// Returns a string representation of the Tuple
func (this *Tuple) String() string {
	return fmt.Sprintf("%v", this.data)
}

// Pops the leftmost item from the Tuple and
// returns it
func (this *Tuple) PopLeft() interface{} {
	if this.Len() < 1 {
		return nil
	}
	ret := this.data[0]
	this.data = this.data[1:]
	return ret
}

// Pops the rightmost item from the Tuple and
// returns it
func (this *Tuple) PopRight() interface{} {
	if this.Len() < 1 {
		return nil
	}
	idx := this.Offset(-1)
	ret := this.data[idx]
	this.data = this.data[:idx]
	return ret
}

// Reverses the Tuple in place
func (this *Tuple) Reverse() {
	for i, j := 0, this.Len()-1; i < j; i, j = i+1, j-1 {
		this.data[i], this.data[j] = this.data[j], this.data[i]
	}
}

// Returns true if the two items are logically "equal"
func TupleElemEq(lhsi interface{}, rhsi interface{}) bool {
	lhsv, rhsv := reflect.ValueOf(lhsi), reflect.ValueOf(rhsi)
	// IsNil panics if type is not interface-y, so use IsValid instead
	if lhsv.IsValid() != rhsv.IsValid() {
		return false
	}
	// TODO: this currently blows up if lhs can't be converted to same
	// type as rhs (e.g. int vs. string)
	switch lhsi.(type) {
	case nil:
		if rhsv.IsValid() {
			return false
		}
	case string:
		if lhsi.(string) != rhsi.(string) {
			return false
		}
	case int, int8, int16, int32, int64:
		if lhsv.Int() != rhsv.Int() {
			return false
		}
	case uint, uintptr, uint8, uint16, uint32, uint64:
		if lhsv.Uint() != rhsv.Uint() {
			return false
		}
	case float32, float64:
		if lhsv.Float() != rhsv.Float() {
			return false
		}
	case *Tuple:
		if lhsi.(*Tuple).Ne(rhsi.(*Tuple)) {
			return false
		}
	default:
		//if !lhsv.IsValid() && !rhsv.IsValid() {
		//return false
		//}
		// TODO: allow user-defined callback for unsupported types
		panic(fmt.Sprintf("unsupported type %#v for Eq in Tuple", lhsi))
	}
	return true
}

// Returns True if this Tuple is elementwise equal to other
func (this *Tuple) Eq(other *Tuple) bool {
	if this.Len() != other.Len() {
		return false
	}
	//return reflect.DeepEqual(this.data, other.data)
	for i := 0; i < this.Len(); i++ {
		lhsi, rhsi := this.Get(i), other.Get(i)
		if !TupleElemEq(lhsi, rhsi) {
			return false
		}
	}
	return true
}

// Returns True if this Tuple is not elementwise equal to other
func (this *Tuple) Ne(other *Tuple) bool {
	return !this.Eq(other)
}

// Support for sort.Sort
func (this *Tuple) Less(i, j int) bool {
	return TupleElemLt(this.Get(i), this.Get(j))
}

// Support for sort.Sort
func (this *Tuple) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

// Support for sorting slices of Tuples
type ByElem []*Tuple

func (a ByElem) Len() int           { return len(a) }
func (a ByElem) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByElem) Less(i, j int) bool { return a[i].Lt(a[j]) }

// Returns true if the item lhsi is logically less than rhsi
func TupleElemLt(lhsi interface{}, rhsi interface{}) bool {
	lhsv, rhsv := reflect.ValueOf(lhsi), reflect.ValueOf(rhsi)
	if lhsv.IsValid() && !rhsv.IsValid() {
		// zero value is considered least
		return false
	}
	switch lhsi.(type) {
	case nil:
		if rhsv.IsValid() {
			return true
		}
	case string:
		if lhsi.(string) < rhsi.(string) {
			return true
		}
	case int, int8, int16, int32, int64:
		if lhsv.Int() < rhsv.Int() {
			return true
		}
	case uint, uintptr, uint8, uint16, uint32, uint64:
		if lhsv.Uint() < rhsv.Uint() {
			return true
		}
	case float32, float64:
		if lhsv.Float() < rhsv.Float() {
			return true
		}
	case *Tuple:
		if lhsi.(*Tuple).Lt(rhsi.(*Tuple)) {
			return true
		}
	default:
		// TODO: allow user-defined callback for unsupported types
		panic(fmt.Sprintf("unsupported type %#v for Lt in Tuple", lhsi))
	}
	return false
}

// Returns True if this Tuple is elementwise less than other
func (this *Tuple) Lt(other *Tuple) bool {
	tlen, olen := this.Len(), other.Len()
	var n int
	if tlen < olen {
		n = tlen
	} else {
		n = olen
	}
	for i := 0; i < n; i++ {
		lhsi, rhsi := this.Get(i), other.Get(i)
		if TupleElemLt(lhsi, rhsi) {
			return true
		} else if !TupleElemEq(lhsi, rhsi) {
			return false
		}
	}
	// if we get here then they matched up to n
	if tlen < olen {
		return true
	}
	return false
}

// Returns True if this Tuple is elementwise less than
// or equal to other
func (this *Tuple) Le(other *Tuple) bool {
	return this.Lt(other) || this.Eq(other)
}

// Returns True if this Tuple is elementwise greater than other
func (this *Tuple) Gt(other *Tuple) bool {
	return !this.Le(other)
}

// Returns True if this Tuple is elementwise greater than
// or equal to other
func (this *Tuple) Ge(other *Tuple) bool {
	return !this.Lt(other)
}

// Returns the position of item in the tuple. The search
// starts from position start. Returns -1 if item is not found
func (this *Tuple) Index(item interface{}, start int) int {
	for i := start; i < this.Len(); i++ {
		if TupleElemEq(this.Get(i), item) {
			return i
		}
	}
	return -1
}

// Returns the number of occurrences of item in the tuple,
// given starting position start.
func (this *Tuple) Count(item interface{}, start int) int {
	ctr := 0
	for i := start; i < this.Len(); i++ {
		if TupleElemEq(this.Get(i), item) {
			ctr += 1
		}
	}
	return ctr
}

// Inserts data from other table into this, starting at offset start
func (this *Tuple) Insert(start int, other *Tuple) {
	this.InsertItems(start, other.data...)
}

// Inserts items into this tuple, starting from offset start
func (this *Tuple) InsertItems(start int, items ...interface{}) {
	start = this.Offset(start)
	rhs := this.Copy().data[start:]
	this.data = append(this.data[:start], items...)
	this.data = append(this.data, rhs...)
}

// Appends all elements from other tuple to this
func (this *Tuple) Append(other *Tuple) {
	this.AppendItems(other.data...)
}

// Appends one or more items to end of data
func (this *Tuple) AppendItems(items ...interface{}) {
	this.data = append(this.data, items...)
}
