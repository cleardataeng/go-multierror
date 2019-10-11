package multierror

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestError_Impl(t *testing.T) {
	var _ error = new(Error)
}

func TestErrorError_custom(t *testing.T) {
	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	fn := func(es []error) string {
		return "foo"
	}

	multi := &Error{Errors: errors, ErrorFormat: fn}
	if multi.Error() != "foo" {
		t.Fatalf("bad: %s", multi.Error())
	}
}

func TestErrorError_default(t *testing.T) {
	expected := `2 errors occurred:
	* foo
	* bar

`

	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := &Error{Errors: errors}
	if multi.Error() != expected {
		t.Fatalf("bad: %s", multi.Error())
	}
}

func TestError_json(t *testing.T) {
	errs := []error{
		errors.New("foo"),
		errors.New("bar"),
	}
	multi := Error{Errors: errs}
	b, err := json.Marshal(&multi)
	if err != nil {
		t.Fatalf("unexpected error; got %#v", err)
	}
	j := `{"errors":["foo","bar"]}`
	if string(b) != j {
		t.Errorf("bad representation; got: %s, want: %s", string(b), j)
	}
	rebuilt := Error{}
	if err = json.Unmarshal(b, &rebuilt); err != nil {
		t.Fatalf("unexpected error; go %#v", err)
	}
	if !reflect.DeepEqual(rebuilt, multi) {
		t.Fatalf("mismatched types; got: %v, want: %v", rebuilt, multi)
	}
}

func TestErrorErrorOrNil(t *testing.T) {
	err := new(Error)
	if err.ErrorOrNil() != nil {
		t.Fatalf("bad: %#v", err.ErrorOrNil())
	}

	err.Errors = []error{errors.New("foo")}
	if v := err.ErrorOrNil(); v == nil {
		t.Fatal("should not be nil")
	} else if !reflect.DeepEqual(v, err) {
		t.Fatalf("bad: %#v", v)
	}
}

func TestErrorWrappedErrors(t *testing.T) {
	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := &Error{Errors: errors}
	if !reflect.DeepEqual(multi.Errors, multi.WrappedErrors()) {
		t.Fatalf("bad: %s", multi.WrappedErrors())
	}
}
