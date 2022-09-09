package content

import (
	"errors"
	"io"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/internal/ioutil"
)

type Verifier interface {
	Verify() error
}

type VerifyReader interface {
	Verifier
	io.Reader
}

type verifyReader struct {
	base     *io.LimitedReader
	verifier digest.Verifier
	verified bool
	err      error
}

func (vr *verifyReader) Verify() error {
	if vr.verified {
		return nil
	}
	if vr.err == nil {
		if vr.base.N > 0 {
			return errors.New("early verify")
		}
	} else if vr.err != io.EOF {
		return vr.err
	}

	if err := ioutil.EnsureEOF(vr.base.R); err != nil {
		vr.err = ErrTrailingData
		return vr.err
	}
	if !vr.verifier.Verified() {
		vr.err = ErrMismatchedDigest
		return vr.err
	}

	vr.verified = true
	vr.err = io.EOF
	return nil
}

func (vr *verifyReader) Read(p []byte) (n int, err error) {
	if vr.err != nil {
		return 0, vr.err
	}

	n, err = vr.base.Read(p)
	if err != nil {
		if err == io.EOF && vr.base.N > 0 {
			err = io.ErrUnexpectedEOF
		}
		vr.err = err
	}
	return
}

func NewVerifyReader(r io.Reader, desc ocispec.Descriptor) VerifyReader {
	verifier := desc.Digest.Verifier()
	lr := &io.LimitedReader{
		R: io.TeeReader(r, verifier),
		N: desc.Size,
	}
	return &verifyReader{
		base:     lr,
		verifier: verifier,
	}
}
