// Package bitmaptable implements an in-memory bitmap table that stores a fixed
// amount of bits for a fixed amount of rows.
//
// Example
//
// Lets assume that humans are given an incremental identifier when they are born.
// We want to create a bitmap table that keeps an overview of whether a human is
// still alive, as well as it's gender.
//
//     const (
//         MAN     bool = false
//         WOMAN   bool = true
//
//         ALIVE   bool = false
//         DEAD    bool = true
//     )
//
//     bm := bitmaptable.New(110*1000*1000*1000, 2);
//
//     bm.Set(123456, 0, WOMAN)
//     bm.Set(123456, 1, DEAD)
//
//     woman, _ := bm.Get(123456, 0)
//     dead, _ := bm.Get(123456, 1)
//
// This bitmap table indexes the gender and aliveness of all ~110 billion human
// beings to have ever walked this earth, and needs give or take 26 GB of memory.
package bitmaptable

import (
	"errors"

	"github.com/boljen/go-bitmap"
)

// These are errors that can be returned by the Bitmaptable.
var (
	ErrIllegalIndex = errors.New("Bitmaptable: Illegal identifier or position")
	ErrIllegalWidth = errors.New("Bitmaptable: Illegal value width, must be between 1 and 64")
)

// Bitmaptable is the basic bitmap table on which all other tables are built.
// The bitmap table stores column-based bit information on a per-row basis.
type Bitmaptable interface {
	// Data returns the underlying data of the bitmap.
	// If copy is true it will copy all the data into a new byteslice.
	Data(copy bool) []byte

	// Rows returns the amount of rows inside this bitmap table.
	Rows() int

	// Columns returns the amount of columns inside this bitmap table.
	Columns() int

	// Get gets the value for the provided row and column tuple.
	Get(row int, column int) (bool, error)

	// Set sets the value for the provided row and column tuple.
	Set(row int, column int, value bool) error
}

// New creates a new Bitmaptable instance.
// Remember that this will allocate rows * columns bits of memory.
// There can be 2^64 rows and 2^16 columns per row, theoretically.
func New(rows, columns int) Bitmaptable {
	return newNTS(rows, columns)
}

// NewTS creates a new thread-safe Bitmaptable instance.
func NewTS(rows, columns int) Bitmaptable {
	return newTS(rows, columns)
}

func newNTS(rows, columns int) *bitmaptable {
	return &bitmaptable{
		rows:    rows,
		columns: columns,
		bitmap:  bitmap.New(columns * rows),
	}
}

type bitmaptable struct {
	rows    int           // Amount of rows.
	columns int           // Amount of columns per row.
	bitmap  bitmap.Bitmap // The actual bitmap
}

// Rows implements Bitmaptable.Rows
func (b *bitmaptable) Rows() int {
	return b.rows
}

// Columns implements Bitmaptable.Columns
func (b *bitmaptable) Columns() int {
	return b.columns
}

// Data implements Bitmaptable.Data
func (b *bitmaptable) Data(c bool) []byte {
	return b.bitmap.Data(c)
}

// Get implements Bitmaptable.Get
func (b *bitmaptable) Get(row int, column int) (bool, error) {
	if column >= b.columns || row >= b.rows {
		return false, ErrIllegalIndex
	}
	return b.bitmap.Get(row*b.columns + column), nil
}

// Set implements Bitmaptable.Set
func (b *bitmaptable) Set(row int, column int, value bool) error {
	if column >= int(b.columns) || row >= int(b.rows) {
		return ErrIllegalIndex
	}
	b.bitmap.Set(row*b.columns+column, value)
	return nil
}
