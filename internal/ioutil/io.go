package ioutil

import (
	"io"
)

// ReaderUnwrapper unpacked the wrapped readers.
type ReaderUnwrapper interface {
	// Unwrap returns the wrapped reader.
	Unwrap() io.Reader
}

// NopCloser is the same as `io.NopCloser` but implements `ReaderUnwrapper`.
func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{
		Reader: r,
	}
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

func (n nopCloser) Unwrap() io.Reader {
	return n.Reader
}
