package minecraft

import "fmt"

type ErrInvalidVersion struct {
	Version string
}

func (e ErrInvalidVersion) Error() string {
	return fmt.Sprintf("invalid minecraft version '%s'", e.Version)
}

var _ error = &ErrInvalidVersion{}
