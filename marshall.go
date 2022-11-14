package rand

import (
	"encoding/json"
	"math/rand"
)

type MarshallableSource interface {
	rand.Source

	Marshall() ([]byte, error)
	Unmarshall([]byte) error
}

func NewSource(seed int64) MarshallableSource {
	var rng rngSource
	rng.Seed(seed)
	return &rng
}

var (
	globalSource = NewSource(1).(*rngSource)
	globalRand   = rand.New(&lockedSource{rngSource: globalSource})
)

// Marshall marshalls the default Source.
func Marshall() ([]byte, error) {
	return globalSource.Marshall()
}

// Unmarshall resets the default Source's internal state to the input's contents.
func Unmarshall(input []byte) error {
	return globalSource.Unmarshall(input)
}

type serializableRngSource struct {
	Tap  int           `json:"tap"`
	Feed int           `json:"feed"`
	Vec  [rngLen]int64 `json:"vec"`
}

func (rng *rngSource) Marshall() ([]byte, error) {
	return json.Marshal(serializableRngSource{
		Tap:  rng.tap,
		Feed: rng.feed,
		Vec:  rng.vec,
	})
}

func (rng *rngSource) Unmarshall(input []byte) error {
	var source serializableRngSource
	if err := json.Unmarshal(input, &source); err != nil {
		return err
	}

	rng.tap = source.Tap
	rng.feed = source.Feed
	rng.vec = source.Vec

	return nil
}

func (r *lockedSource) Marshall() ([]byte, error) {
	r.lk.Lock()
	defer r.lk.Unlock()

	return r.rngSource.Marshall()
}

func (r *lockedSource) Unmarshall(input []byte) error {
	r.lk.Lock()
	defer r.lk.Unlock()

	return r.rngSource.Unmarshall(input)
}
