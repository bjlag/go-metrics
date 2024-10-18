package model

import "errors"

var (
	ErrInvalidID    = errors.New("metric ID not specified")
	ErrInvalidType  = errors.New("metric type is invalid")
	ErrInvalidValue = errors.New("metric value is invalid")
)
