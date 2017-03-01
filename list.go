package collection

import (
	"bytes"
	"fmt"
	"sync"
)

// List is a dynamically sized list akin to List in the .NET world,
// ArrayList in the Java world, or vector in the C++ world.
type List struct {
	underlyer []interface{}
	key       sync.RWMutex
}

// NewList creates a new list which contains the elements provided.
func NewList(entries ...interface{}) *List {
	return &List{
		underlyer: entries,
	}
}

// Add appends an entry to the logical end of the List.
func (l *List) Add(entries ...interface{}) {
	l.key.Lock()
	defer l.key.Unlock()
	l.underlyer = append(l.underlyer, entries...)
}

// AddAt injects values beginning at `position`. If multiple
// values are provided in `entries` they are placed in the same
// order they are provided.
func (l *List) AddAt(position uint, entries ...interface{}) {
	l.key.Lock()
	defer l.key.Unlock()
	originalLength := len(l.underlyer)
	l.underlyer = append(l.underlyer, entries...)
	for i := 0; i < len(entries); i++ {
		l.swap(position+uint(i), uint(originalLength+i))
	}
}

// Enumerate lists each element present in the collection
func (l *List) Enumerate() Enumerable {
	retval := make(chan interface{})

	go func() {
		l.key.RLock()
		defer l.key.RUnlock()

		for _, entry := range l.underlyer {
			retval <- entry
		}
	}()

	return retval
}

// Get retreives the value stored in a particular position of the list.
// If no item exists at the given position, the second parameter will be
// returned as false.
func (l *List) Get(pos uint) (interface{}, bool) {
	l.key.RLock()
	defer l.key.RUnlock()

	if pos > uint(len(l.underlyer)) {
		return nil, false
	}
	return l.underlyer[pos], true
}

// IsEmpty tests to see if this List has any elements present.
func (l *List) IsEmpty() bool {
	l.key.RLock()
	defer l.key.RUnlock()
	return 0 == len(l.underlyer)
}

// Length returns the number of elements in the List.
func (l *List) Length() uint {
	l.key.RLock()
	l.key.RUnlock()
	return uint(len(l.underlyer))
}

// String generates a textual representation of the List for the sake of debugging.
func (l *List) String() string {
	l.key.RLock()
	defer l.key.RUnlock()

	builder := bytes.NewBufferString("[")

	for i, entry := range l.underlyer {
		if i >= 15 {
			builder.WriteString("... ")
			break
		}
		builder.WriteString(fmt.Sprintf("%v ", entry))
	}
	builder.Truncate(builder.Len() - 1)
	builder.WriteRune(']')
	return builder.String()
}

// Swap switches the values that are stored at positions `x` and `y`
func (l *List) Swap(x, y uint) bool {
	l.key.Lock()
	defer l.key.Unlock()
	return l.swap(x, y)
}

func (l *List) swap(x, y uint) bool {
	count := uint(len(l.underlyer))
	if x < count && y < count {
		temp := l.underlyer[x]
		l.underlyer[x] = l.underlyer[y]
		l.underlyer[y] = temp
		return true
	}
	return false
}
