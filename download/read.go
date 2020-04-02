package download

import "io"

type ReadCounter interface {
	io.ReadCloser

	// Count returns the total sum of bytes read so far
	Count() uint64
	// Total returns the expected final sum count of bytes
	ExpectedTotal() uint64
	// IsFinished returns true when the sum of read bytes equals that of the
	// expected final count
	IsFinished() bool
}

type CountingReader struct {
	io.Reader
	// Total is the expected final count of the Reader
	Total uint64
	count uint64
}

func (c *CountingReader) Read(p []byte) (n int, err error) {
	read, err := c.Reader.Read(p)
	c.count += uint64(read)
	return read, err
}

// Close pretends that an io.Reader can be an io.ReadCloser, even when it's not
// I don't think this should ever cause a problem, it would just do nothing
func (c *CountingReader) Close() error {
	if clsr, ok := c.Reader.(io.Closer); ok {
		return clsr.Close()
	}
	return nil
}

func (c *CountingReader) Count() uint64 {
	return c.count
}

func (c *CountingReader) ExpectedTotal() uint64 {
	return c.Total
}

func (c *CountingReader) IsFinished() bool {
	return c.count >= c.Total
}

var _ io.ReadCloser = &CountingReader{}
var _ ReadCounter = &CountingReader{}
