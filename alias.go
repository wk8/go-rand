package rand

import "math/rand"

// This file aliases what we need from math/rand to make this package a complete drop-in replacement.

var (
	New     = rand.New
	NewZipf = rand.NewZipf
)

type Source interface{ rand.Source }
type Source64 interface{ rand.Source64 }
type Rand rand.Rand
type Zipf rand.Zipf
