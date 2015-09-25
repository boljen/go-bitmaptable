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

import "errors"

// These are errors that can be returned by the Bitmaptable.
var (
	ErrIllegalIndex = errors.New("Bitmaptable: Illegal identifier or position")
	ErrIllegalWidth = errors.New("Bitmaptable: Illegal value width, must be between 1 and 64")
)

// Values for bitwise operations.
var tA = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}
var tB = [8]byte{254, 253, 251, 247, 239, 223, 191, 127}

// setBit sets bit "bit" of byte b to value v.
// It doesn't throw an error if the position is invalid.
func setBit(by byte, bit int, v bool) byte {
	if v {
		return by | tA[bit]
	}
	return by & tB[bit]
}

// getBit gets the value of bit p inside byte b.
func getBit(by byte, bit int) bool {
	return by&tA[bit] != 0
}

func calculateSize(rows, columns int) int {
	size := (rows * columns) / 8
	if (rows*columns)%8 != 0 {
		size++
	}
	return size
}

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
		rows:    uint64(rows),
		columns: uint16(columns),
		bitmap:  make([]byte, calculateSize(rows, columns)),
	}
}

type bitmaptable struct {
	rows    uint64 // Amount of rows.
	columns uint16 // Amount of columns per row.
	bitmap  []byte // The actual bitmap
}

// getPos returns the position of the bit.
// The first integer is the actual byte index.
// The second integer is the bit inside the byte.
func (b *bitmaptable) getPos(row, column int) (int, int) {
	pos := row*int(b.columns) + column
	return pos / 8, pos % 8
}

// Rows implements Bitmaptable.Rows
func (b *bitmaptable) Rows() int {
	return int(b.rows)
}

// Columns implements Bitmaptable.Columns
func (b *bitmaptable) Columns() int {
	return int(b.columns)
}

// Data implements Bitmaptable.Data
func (b *bitmaptable) Data(c bool) []byte {
	if !c {
		return b.bitmap
	}
	data := make([]byte, len(b.bitmap))
	copy(data, b.bitmap)
	return data
}

// Get implements Bitmaptable.Get
func (b *bitmaptable) Get(row int, column int) (bool, error) {
	if column >= int(b.columns) || row >= int(b.rows) {
		return false, ErrIllegalIndex
	}
	by, bit := b.getPos(row, column)
	return getBit(b.bitmap[by], bit), nil
}

// Set implements Bitmaptable.Set
func (b *bitmaptable) Set(row int, column int, value bool) error {
	if column >= int(b.columns) || row >= int(b.rows) {
		return ErrIllegalIndex
	}
	by, bit := b.getPos(row, column)
	b.bitmap[by] = setBit(b.bitmap[by], bit, value)
	return nil
}

/*

TODO: Impelement these methods using efficient bitwise operations.

// SetMulti sets a value accross multiple bit columns.
// Values that cross up to 64 bits can be set.
func (b *Bitmaptable) SetMulti(row int, column int, width int, value uint64) error {
	if column+width >= int(b.columns) || row >= int(b.rows) {
		return ErrIllegalIndex
	} else if width > 64 || width <= 0 {
		return ErrIllegalWidth
	}

	// First byte might not start at the first bit
	// Other bytes will start after the first bit.

	by, bit := b.getPos(row, column)
	fmt.Println(by, bit)
	return nil
}

// GetMulti gets a value accross multiple bit columns.
// Values that cross up to 64 bits can be retrieved.
func (b *Bitmaptable) GetMulti(row int, column int, width int) (uint64, error) {
	if column+width >= int(b.columns) || row >= int(b.rows) {
		return 0, ErrIllegalIndex
	} else if width > 64 || width <= 0 {
		return 0, ErrIllegalWidth
	}

	by, bit := b.getPos(row, column)
	fmt.Println(by, bit)

	return 0, nil
}
*/
