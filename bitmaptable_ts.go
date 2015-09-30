package bitmaptable

import "sync"

// ts is a Thread-Safe implementation of the Bitmaptable struct.
type ts struct {
	mu *sync.Mutex
	b  *bitmaptable
}

func newTS(rows, columns int) *ts {
	return &ts{
		mu: new(sync.Mutex),
		b:  newNTS(rows, columns),
	}
}

// Rows implements Bitmaptable.Rows
func (t *ts) Rows() int {
	return t.b.Rows()
}

// Columns implements Bitmaptable.Columns
func (t *ts) Columns() int {
	return t.b.Columns()
}

// Data implements Bitmaptable.Data
func (t *ts) Data(c bool) []byte {
	t.mu.Lock()
	data := t.b.Data(c)
	t.mu.Unlock()
	return data
}

// Get implements Bitmaptable.Get
func (t *ts) Get(row int, column int) (bool, error) {
	return t.b.Get(row, column)
}

// Set implements Bitmaptable.Set
func (t *ts) Set(row int, column int, value bool) error {
	t.mu.Lock()
	err := t.b.Set(row, column, value)
	t.mu.Unlock()
	return err
}
