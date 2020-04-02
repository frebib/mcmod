package api

import (
	"fmt"
	"net/http"
)

type ErrHttpStatus struct {
	Req     *http.Request
	Code    int
	ErrBody string
}

func (e *ErrHttpStatus) Error() string {
	return fmt.Sprintf("%s %s: %s (%s)",
		e.Req.Method,
		e.Req.URL.Path,
		http.StatusText(e.Code),
		e.ErrBody,
	)
}

type ErrNoSuchAddon struct {
	Name string
	ID   int
}

func (e *ErrNoSuchAddon) Error() string {
	if e.Name != "" {
		return fmt.Sprintf("no mod found with name '%s'", e.Name)
	} else if e.ID != 0 {
		return fmt.Sprintf("no mod found with id %d", e.ID)
	} else {
		return "no mod found matching criteria"
	}
}
