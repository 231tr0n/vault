package errorwrap

import (
	"fmt"
	"runtime"
)

// ErrWrap takes an error, wraps it with the stack trace and returns it.
func ErrWrap(err error) error {
	if err == nil {
		return nil
	}

	counter, _, _, ok := runtime.Caller(1)
	if ok {
		temp := runtime.FuncForPC(counter)

		return fmt.Errorf("%s\n\t%w", temp.Name(), err)
	}

	return fmt.Errorf("%s\n\t%w", "Unknown", err)
}
