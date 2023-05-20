package web

import (
	"fmt"
)

const (
	ErrParseFormMsg    = "Neplatný formulář."
	ErrParamMsg        = "Neplatný parametr."
	ErrItemNotFoundMsg = "Položka nenalezena."
	ErrCartEmptyMsg    = "Košík je prázdný."
	ErrTOCMsg          = "Obchodní podmínky nebyly odsouhlaseny."
	ErrMsg             = "Chyba na straně serveru."
	ErrCheckoutMsg     = "Dodací formulář je neplatný."
)

type Error struct {
	Err     error
	Message string
	Code    int
}

func NewWebError(err error, message string, code int) *Error {
	return &Error{Err: err, Message: message, Code: code}
}

func (e *Error) Error() string {
	return fmt.Sprintf(e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}
