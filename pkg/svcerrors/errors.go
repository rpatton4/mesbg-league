package svcerrors

import "errors"

// NotFound is returned when a requested resource is not found.
var NotFound = errors.New("not found") // I disagree that this is an error, but the GoLang mafia says it is idiomatic to make this mistake

// InvalidID is returned when an ID is expected but what is provided cannot be used
var InvalidID = errors.New("ID is either empty or not a valid identifier")
