# Bitmaptable (Go)

Package bitmaptable implements an in-memory bitmap table that stores a fixed
amount of bits for a fixed amount of rows.

## Example

Lets assume that humans are given an incremental identifier when they are born.
We want to create a bitmap table that keeps an overview of a the human is
still alive, as well as it's gender.

    const (
        MAN     bool = false
        WOMAN   bool = true

        ALIVE   bool = false
        DEAD    bool = true
    )

    bm := bitmaptable.New(110*1000*1000*1000, 2);

    bm.Set(123456, 0, WOMAN)
    bm.Set(123456, 1, DEAD)

    woman, _ := bm.Get(123456, 0)
    dead, _ := bm.Get(123456, 1)

This bitmap table indexes the gender and aliveness of all ~110 billion human
beings to have ever walked this earth, and needs give or take 26 GB of memory.

## Documentation

See [godoc](https://godoc.org/github.com/boljen/go-bitmap) for more information.

## License

Released under the MIT license.
