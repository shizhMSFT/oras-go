package dag

import (
	"context"
	"errors"
	"io"
)

// ReadContent safely reads the content described by the descriptor.
// The read content is verified against the size and the digest.
func ReadContent(r io.Reader, desc Descriptor) ([]byte, error) {
	// verify while reading
	verifier := desc.Digest.Verifier()
	r = io.TeeReader(r, verifier)
	buf := make([]byte, desc.Size)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if !verifier.Verified() {
		return nil, errors.New("digest verification failed")
	}

	// ensure EOF
	var peek [1]byte
	_, err = io.ReadFull(r, peek[:])
	if err != io.EOF {
		return nil, errors.New("trailing data")
	}

	return buf, nil
}

// FetchContent safely fetches the content described by the descriptor.
// The fetched content is verified against the size and the digest.
func FetchContent(ctx context.Context, fetcher Fetcher, desc Descriptor) ([]byte, error) {
	rc, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ReadContent(rc, desc)
}
