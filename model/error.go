package model

import "errors"

var (
	ErrMissingDependency = errors.New("missing dependency")
	ErrWrongRedirect     = errors.New("the m.Redirect has not a valid protocol")
	ErrMethodNotAllowed  = errors.New("the method is not allowed")
)
