package bitmaptable

import "testing"

func TestNew(t *testing.T) {
	New(10, 5)
	NewTS(10, 5)
}

func TestBitmaptable(t *testing.T) {
	bm := newNTS(10, 5)
	if bm.rows != 10 || bm.columns != 5 || len(bm.bitmap) != 7 {
		t.Fatal("wrong configuration")
	}
	if bm.Rows() != 10 || bm.Columns() != 5 {
		t.Fatal("wrong rows and/or columns")
	}

	data := bm.Data(false)
	data[1] = 123
	if bm.bitmap[1] != 123 {
		t.Fatal("didn't return the same slice")
	}

	data2 := bm.Data(true)
	if data2[1] != 123 {
		t.Fatal("wrong copy?")
	}
	data2[1] = 111
	if data[1] == 111 || bm.bitmap[1] == 111 {
		t.Fatal("wrong copy")
	}
}

func TestBitmaptableGetSet(t *testing.T) {
	b := newNTS(1000, 12)
	if err := b.Set(1001, 0, true); err != ErrIllegalIndex {
		t.Fatal("illegal index must be returned")
	}
	if err := b.Set(5, 11, true); err != nil {
		t.Fatal("unexpected error", err)
	}
	if err := b.Set(5, 10, true); err != nil {
		t.Fatal("unexpected error", err)
	}
	if v, err := b.Get(5, 11); err != nil || !v {
		t.Fatal("wrong return")
	}

	if err := b.Set(5, 11, false); err != nil {
		t.Fatal("unexpected error", err)
	}

	if v, err := b.Get(5, 11); err != nil || v {
		t.Fatal("wrong return")
	}

	if _, err := b.Get(1001, 0); err != ErrIllegalIndex {
		t.Fatal("illegal index")
	}
}
