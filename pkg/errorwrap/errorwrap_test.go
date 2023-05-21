package errorwrap_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/231tr0n/vault/pkg/errorwrap"
)

func failTestCase(t *testing.T, i, o, w any) {
	t.Helper()
	t.Error("Input:", i, "|", "Output:", o, "|", "Want:", w)
}

func TestErrWrap(t *testing.T) {
	s := "github.com/231tr0n/vault/pkg/errorwrap_test.TestErrWrap\n\ttest"
	out := fmt.Sprint(errorwrap.ErrWrap(errors.New("test")))
	t.Log(out)
	if s != out {
		failTestCase(t, "test", out, s)
	}
}
