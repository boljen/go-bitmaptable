// Package bitmap implements a configurable bitmap that can efficiently store
// in-memory boolean metadata.
// To put things into perspective; It should be capable of storing the gender
// for the entire world population inside less than 1GB of memory.
package bitmap

import (
	"errors"
	"sync"
)

// ErrIllegalIndex is returned by the Set and Get methods when the defined
// identifier or position index is invalid.
var ErrIllegalIndex = errors.New("Bitmap: Illegal identifier or position")

// Values for bitwise operations.
var tA = [8]byte{1, 2, 4, 8, 16, 32, 64, 128}
var tB = [8]byte{254, 253, 251, 247, 239, 223, 191, 127}

// setBit sets bit p of byte b to value v.
// It doesn't throw an error if the position is invalid.
func setBit(b byte, p int, v bool) byte {
	if v {
		return b | tA[p]
	}
	return b & tB[p]
}

// getBit gets the value of bit p inside byte b.
func getBit(b byte, p int) bool {
	return b&tA[p] != 0
}

// Bitmap is a bitmap in the literal sense of the word.
// It's created for a fixed amount of items and for a
// fixed amount of boolean properties for those items.
type Bitmap interface {
	// Get returns the boolean data associated with the position of the
	// specifier identifier. It returns ErrIllegalIndex if the provided
	// identifier or position is illegal.
	Get(id int, position int) (bool, error)

	// Set sets the boolean data associated with the identifier and position.
	// It returns ErrIllegalIndex if the provided identifier or position is
	// illegal.
	Set(id int, position int, value bool) error
}

// New creates a new Bitmap instance.
// It will allocate (size * width) bits of data.
func New(size int, width int) Bitmap {
	as := size * width
	rest := as % 8
	as = as / 8
	if rest != 0 {
		as++
	}
	return &bitmap{
		i: size,
		w: width,
		a: make([]byte, as),
	}
}

type bitmap struct {
	i int    // Amount of identities.
	w int    // Amount of properties per identity.
	a []byte // Array containing the actual data.
}

func (s *bitmap) Set(id int, pos int, value bool) error {
	if pos >= s.w || id > s.i {
		return ErrIllegalIndex
	}

	// Identify the specific byte that contains the bit that must be changed.
	l := id*s.w + pos
	by := l / 8
	data := s.a[by]

	s.a[by] = setBit(data, l%8, value)
	return nil
}

func (s *bitmap) Get(id int, pos int) (bool, error) {
	if pos >= s.w || id > s.i {
		return false, ErrIllegalIndex
	}

	// Identify the specific byte that contains the bit that must be returned.
	l := id*s.w + pos
	by := l / 8
	data := s.a[by]

	return getBit(data, l%8), nil
}

// NewTS creates a new Thread-safe bitmap instance.
func NewTS(size int, width int) Bitmap {
	return &bitmapTS{
		mu:     new(sync.Mutex),
		Bitmap: New(size, width),
	}
}

type bitmapTS struct {
	Bitmap
	mu *sync.Mutex
}

func (s *bitmapTS) Set(id int, pos int, value bool) error {
	s.mu.Lock()
	e := s.Bitmap.Set(id, pos, value)
	s.mu.Unlock()
	return e
}

func (s *bitmapTS) Get(id int, pos int) (bool, error) {
	s.mu.Lock()
	v, e := s.Bitmap.Get(id, pos)
	s.mu.Unlock()
	return v, e
}
