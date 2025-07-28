package svcerrors

import "errors"

// ErrNotFound is returned when a requested resource is not found.
var ErrNotFound = errors.New("not found") // I disagree that this is an error, but the GoLang mafia says it is idiomatic to make this mistake

// ErrInvalidID is returned when an ID is expected but what is provided cannot be used
var ErrInvalidID = errors.New("ID is either empty or not a valid identifier")
