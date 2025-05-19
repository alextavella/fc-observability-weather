package exception

import "errors"

var (
	ErrInvalidZipcode    = errors.New("invalid zipcode")
	ErrCanNotFindZipcode = errors.New("can not find zipcode")
)
