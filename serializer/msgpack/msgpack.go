package msgpack

import (
	"github.com/navisot/go-url-shortener/shortener"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

// Decode decodes MSGPACK
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	red := &shortener.Redirect{}

	if err := msgpack.Unmarshal(input, red); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return red, nil
}

// Encode encodes MSGPACK
func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {

	enc, err := msgpack.Marshal(redirect)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return enc, nil
}
