package json

import (
	"encoding/json"
	"github.com/navisot/go-url-shortener/shortener"
	"github.com/pkg/errors"
)

type Redirect struct{}

// Decode decodes JSON
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	red := &shortener.Redirect{}

	if err := json.Unmarshal(input, red); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return red, nil
}

// Encode encodes JSON
func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {

	enc, err := json.Marshal(redirect)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return enc, nil
}
