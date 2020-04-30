package cmd

import (
	"fmt"
)

type ErrNoMatch struct {
	Filter ModFilter
}

func (nm *ErrNoMatch) Error() string {
	return fmt.Sprintf("no match found for %s", nm.Filter.String())
}
