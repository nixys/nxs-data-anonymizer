package progressreader

import "io"

type ProgressReader struct {

	// Source reader
	r io.Reader

	// Total read bytes
	tb int64
}

func Init(r io.Reader) *ProgressReader {
	return &ProgressReader{
		r:  r,
		tb: 0,
	}
}

func (pr *ProgressReader) Read(dst []byte) (int, error) {

	// Proxy data from source reader to anonymizer
	n, err := pr.r.Read(dst)

	// Save read bytes count
	pr.tb += int64(n)

	return n, err
}

// Bytes returns total read bytes from source reader
func (pr *ProgressReader) Bytes() int64 {
	return pr.tb
}
