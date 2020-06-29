package common

import (
	"fmt"
	dbCommot "github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
)

type ErrorType int

const (
	NotFound ErrorType = iota
	Conflict
	InternalError
)

type Error struct {
	Type        ErrorType
	Description string
	Cause       error
}

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %v", e.Description, e.Cause)
	} else {
		return e.Description
	}
}

func (e Error) Unwrap() error {
	return e.Cause
}

func NewNotFountError() error {
	return Error{Type: NotFound, Description: "not found"}
}

func NewConflictError(description string) error {
	return Error{Type: Conflict, Description: description}
}

func NewInternalError(description string, cause error) error {
	return Error{Type: InternalError, Description: description, Cause: cause}
}

func MaybeNewInternalError(description string, err error) error {
	if err != nil {
		return NewInternalError(description, err)
	}
	return err
}

func MaybeNewNotFoundOrInternalError(description string, err error) error {
	if err != nil {
		return NewNotFoundOrInternalError(description, err)
	}
	return err
}

func NewNotFoundOrInternalError(description string, err error) error {
	if dbCommot.IsNoRowsError(err) {
		return NewNotFountError()
	} else {
		return NewInternalError(description, err)
	}
}
