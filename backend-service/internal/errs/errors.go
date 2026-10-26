package errs

import "errors"

var (
	ErrWrongConfig      = errors.New("invalid db configuration: can not parse connection string")
	ErrEstablishConnect = errors.New("can not establish db connection")
)
