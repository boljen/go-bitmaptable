package bitmap

import "testing"

func TestNew(t *testing.T) {
	b := New(1001, 4)
	if bm, ok := b.(*bitmap); !ok {
		t.Fatal("wrong bitmap type")
	} else if len(bm.a) != 501 {
		t.Fatal("wrong length")
	}
}

func TestBitmapGetSet(t *testing.T) {
	b := New(1000, 12)
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

func TestBitmapTS(t *testing.T) {
	b := NewTS(1000, 4)
	if err := b.Set(50, 1, true); err != nil {
		t.Fatal(err)
	}

	if v, err := b.Get(50, 1); err != nil || !v {
		t.Fatal("wrong value or unexpected error", v, err)
	}
}
