package errors

import "fmt"

var ErrOpenFile = fmt.Errorf("error opening file")
var ErrReadFile = fmt.Errorf("error reading file")
var ErrFriday = fmt.Errorf("error - friday's are for fun, not files!")
