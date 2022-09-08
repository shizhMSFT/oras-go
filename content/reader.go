/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package content

import (
	"errors"
	"fmt"
	"io"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/internal/ioutil"
)

var (
	// ErrInvalidDescriptorSize is returned by ReadAll() when
	// the descriptor has an invalid size.
	ErrInvalidDescriptorSize = errors.New("invalid descriptor size")

	// ErrMismatchedDigest is returned by ReadAll() when
	// the descriptor has an invalid digest.
	ErrMismatchedDigest = errors.New("mismatched digest")

	// ErrTrailingData is returned by ReadAll() when
	// there exists trailing data unread when the read terminates.
	ErrTrailingData = errors.New("trailing data")
)

type reader struct {
	base      io.Reader       // underlying reader
	remaining int64           // bytes remaining
	verifier  digest.Verifier // integrity verifier
	err       error           // last error
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}

	if r.remaining > 0 {
		if int64(len(p)) > r.remaining {
			p = p[0:r.remaining]
		}
		n, err = r.base.Read(p)
		r.remaining -= int64(n)

		// detect early EOF
		if err == io.EOF && r.remaining > 0 {
			err = io.ErrUnexpectedEOF
			r.err = err
			return
		}
	}

	// complete reading
	if r.remaining <= 0 {
		// verify the size of the read content
		if err = ioutil.EnsureEOF(r.base); err != nil {
			err = ErrTrailingData
			r.err = err
			return
		}
		// verify the digest of the read content
		if !r.verifier.Verified() {
			err = ErrMismatchedDigest
			r.err = err
			return
		}
	}
	return
}

func Reader(r io.Reader, desc ocispec.Descriptor) io.Reader {
	verifier := desc.Digest.Verifier()
	return &reader{
		base:      io.TeeReader(r, verifier), // verify while reading
		remaining: desc.Size,
		verifier:  verifier,
	}
}

// ReadAll safely reads the content described by the descriptor.
// The read content is verified against the size and the digest.
func ReadAll(r io.Reader, desc ocispec.Descriptor) ([]byte, error) {
	if desc.Size < 0 {
		return nil, ErrInvalidDescriptorSize
	}
	buf := make([]byte, desc.Size)

	// verify while reading
	r = Reader(r, desc)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}
	return buf, nil
}
