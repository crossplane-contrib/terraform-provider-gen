package generator

import (
	"testing"

	"github.com/pkg/errors"
)

func TestMultiError(t *testing.T) {
	heading := "testing:"
	me := NewMultiError(heading)
	if len(me.Errors()) > 0 {
		t.Errorf("Observed MultiError.Errors() count=%d", len(me.Errors()))
	}
	if me.Error() != "" {
		t.Errorf("Got output from .Error() with empty list. Note: heading should only print if there are errors.")
	}
	me.Append(errors.New("derp"))
	if len(me.Errors()) != 1 {
		t.Errorf("Observed unexpected MultiError.Errors(). expected=1, actual=%d", len(me.Errors()))
	}
	expected := `testing:
 - derp`
	if me.Error() != expected {
		t.Errorf("Observed unexpected MultiError.Error() ouput. Expected:\n%s\n-----\nActual:\n%s\n", expected, me.Error())
	}
}
