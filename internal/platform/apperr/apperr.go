// Package apperr define los errores de negocio que la capa HTTP traduce a códigos de estado.
package apperr

import "fmt"

type Kind int

const (
	KindInvalid  Kind = iota + 1 // 400
	KindNotFound                 // 404
	KindConflict                 // 409
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string { return e.Message }

func Invalid(format string, args ...any) error {
	return &Error{Kind: KindInvalid, Message: fmt.Sprintf(format, args...)}
}

func NotFound(format string, args ...any) error {
	return &Error{Kind: KindNotFound, Message: fmt.Sprintf(format, args...)}
}

func Conflict(format string, args ...any) error {
	return &Error{Kind: KindConflict, Message: fmt.Sprintf(format, args...)}
}
