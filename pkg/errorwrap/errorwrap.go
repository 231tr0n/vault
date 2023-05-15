package errorwrap

import (
	"fmt"
	"runtime"
)

func ErrWrap(err error) error {
	if err != nil {
		counter, _, _, ok := runtime.Caller(1)
		if ok {
			temp := runtime.FuncForPC(counter)
			return fmt.Errorf("%s\n\t%w", temp.Name(), err)
		}
		return err
	}
	return nil
}
