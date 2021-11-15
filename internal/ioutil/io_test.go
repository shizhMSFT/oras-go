package ioutil

import "testing"

func TestNopCloserInterface(t *testing.T) {
	if _, ok := NopCloser(nil).(ReaderUnwrapper); !ok {
		t.Error("nopCloser{} does not conform ReaderUnwrapper")
	}
}
