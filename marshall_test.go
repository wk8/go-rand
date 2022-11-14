package rand

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRngMarshallUnmarshall(t *testing.T) {
	nIterations := 100
	nSubtests := 100

	source := &rngSource{}
	source.Seed(time.Now().UnixNano())

	// let's unroll the source, and store the state for each iteration
	values := make([]int, nIterations)
	states := make([][]byte, nIterations)

	r := New(source)
	for i := 0; i < nIterations; i++ {
		states[i] = marshall(t, source)
		values[i] = r.Int()
	}

	var wg sync.WaitGroup
	wg.Add(nSubtests)

	for i := 0; i < nSubtests; i++ {
		go func(permutationSeed int64) {
			indices := randomPermutation(nIterations, permutationSeed)

			s := &rngSource{}
			r := New(s)

			for _, j := range indices {
				unmarshall(t, s, states[j])
				actual := r.Int()
				expected := values[j]

				if actual != expected {
					t.Errorf("iteration %d, expected %d, got %d", j, expected, actual)
				}
			}

			wg.Done()
		}(time.Now().UnixNano())
	}

	wg.Wait()
}

// Tests that our rng generator generates the same value as math/rand's
func TestSameAsMathRand(t *testing.T) {
	nIterations := 100
	nSubtests := 100

	for i := 0; i < nSubtests; i++ {
		seed := time.Now().UnixNano()

		t.Run(fmt.Sprintf("with seed %d", seed), func(t *testing.T) {
			r := New(NewSource(seed))
			original := rand.New(rand.NewSource(seed))

			for j := 0; j < nIterations; j++ {
				actual := r.Int()
				expected := original.Int()

				if actual != expected {
					t.Errorf("iteration %d, expected %d, got %d", j, expected, actual)
				}
			}
		})
	}
}

// Helpers below

func marshall(t *testing.T, source MarshallableSource) []byte {
	data, err := source.Marshall()
	if err != nil {
		t.Errorf("marshall error: %v", err)
	}
	return data
}

func unmarshall(t *testing.T, source MarshallableSource, input []byte) {
	if err := source.Unmarshall(input); err != nil {
		t.Errorf("unmarshall error: %v", err)
	}
}

// randomPermutation returns a random permutation of [0 .. n) given a random seed
func randomPermutation(n int, seed int64) []int {
	r := rand.New(rand.NewSource(seed))

	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i
	}
	r.Shuffle(n, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}
