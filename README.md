[![CircleCI](https://dl.circleci.com/status-badge/img/gh/wk8/go-rand/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/wk8/go-rand/tree/main)

# Marshallable golang random generator

## Why?

This package is a drop-in replacement for [`math/rand`](https://pkg.go.dev/math/rand), except it also allows you to serialize the random sources' internal state, and restore it at will.

That can come in handy when debugging flaky tests: if you know an initial seed for your whole suite that makes some test fail, you don't need to replay all your tests up to the one that fails any more: just dump the random generator's internal state at the beginning of the failing test once, and you can then load that state to only run the offending test.

## Usage

It's a drop-in replacement for `math/rand` with just a couple of methods on top:

```go
import (
"github.com/wk8/go-rand"
)

source := rand.NewSource(time.Now().UnixNano())
state, err := source.Marshall()
if err != nil {
// ...
}
r := rand.New(source)

// say the next call to r.Int() yields 12

if err := source.Unmarshall(); err != nil {
// ...
}

// your source is now re-set to what it was when it was marshalled
// in particular the next call to r.Int() will yield 12 again
```

or if you prefer the global version:
```go
state, err := rand.Marshall()
// ...
err = rand.Unmarshall(state)
```

## How do I know I can trust this repo?

It uses the exact same RNG as `math/rand`; in fact most of this repo is copy-pasted directly from there:
* [`rng.go` is copy-pasted entirely from `math/rand/`](https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/math/rand/rng.go)
* [so is rand.go](https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/math/rand/rand.go)

That's also why given the same seed, it will give you the exact same pseudo-random values as `math/rand`'s generator.

And for this very same reason, you shouldn't trust this repo any more than `math/rand`'s generator, which is not crypto-secure
