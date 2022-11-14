package rand

import (
	"encoding/json"
	"errors"
)

type MarshallableSource interface {
	Source

	Marshall() ([]byte, error)
	Unmarshall([]byte) error
}

var errNotAMarshallableSource = errors.New("not a marshallable source")

// Marshall marshalls the default Source.
func Marshall() ([]byte, error) {
	s, ok := globalRand.src.(MarshallableSource)
	if !ok {
		return nil, errNotAMarshallableSource
	}
	return s.Marshall()
}

// Unmarshall resets the default Source's internal state to the input's contents.
func Unmarshall(input []byte) error {
	s, ok := globalRand.src.(MarshallableSource)
	if !ok {
		return errNotAMarshallableSource
	}
	return s.Unmarshall(input)
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

	return r.src.Marshall()
}

func (r *lockedSource) Unmarshall(input []byte) error {
	r.lk.Lock()
	defer r.lk.Unlock()

	return r.src.Unmarshall(input)
}
