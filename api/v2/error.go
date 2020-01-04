package v2

import "fmt"

type Error struct {
	s string
}

func (e *Error) Error() string {
	return fmt.Sprintf(e.s)
}

func NewError(text string) error {
	return &Error{
		s: text,
	}
}
