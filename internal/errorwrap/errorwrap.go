package errorwrap

import (
	"fmt"
	"runtime"
)

// Wrap takes an error, wraps it with the stack trace and returns it.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	counter, _, _, ok := runtime.Caller(1)
	if ok {
		temp := runtime.FuncForPC(counter)

		return fmt.Errorf("error: %s\n%w", temp.Name(), err)
	}

	return fmt.Errorf("error: %s\n%w", "Unknown", err)
}
