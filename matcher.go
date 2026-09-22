package assert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// Errors returned when comparisons go bad.
var (
	errInvalidTypeComparison = errors.New("invalid type for comparison")
	errIncompatibleTypes     = errors.New("incompatible types for comparison")
	errInvalidType           = errors.New("invalid type")
)

// These are the basic types, distilled from the variety of more specific types.
type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	stringKind
	uintKind
	sliceKind
)

// Matcher holds the current state of the assertion.
type Matcher struct {
	t         *testing.T
	actual    any
	hasActual bool
	match     bool
}

// With creates a new matcher using the current test reporter.
func With(t *testing.T) *Matcher {
	m := new(Matcher)
	m.t = t
	return m
}

// That specifies the actual value under test.
func (m *Matcher) That(actual any) *Matcher {
	m.requireTestingT()
	m.actual = actual
	m.hasActual = true
	return m
}

// ThatPanics expects the test function to panic. This is the only matcher
// that doesn't return (*Matcher) since a panic is a terminal condition; the
// Matcher should not be reused after calling it.
//
// Any panic satisfies the assertion — ThatPanics cannot tell an expected panic
// from one raised for an unrelated reason, such as a nil dereference inside
// the function. The recovered value is logged via t.Log to make it visible
// with `go test -v` or when the test otherwise fails.
func (m *Matcher) ThatPanics(actual func()) {
	m.requireTestingT()
	m.t.Helper()
	defer func() {
		// The deferred closure is its own call frame, so it needs t.Helper()
		// of its own. The one on ThatPanics above does not cover it.
		m.t.Helper()
		if r := recover(); r == nil {
			m.t.Error("Did not panic.")
			m.match = false
		} else {
			m.t.Logf("recovered from panic: %v", r)
		}
	}()
	m.match = true
	actual()
}

// IsNil passes when the actual value is an untyped nil or a nillable value
// (chan, func, map, pointer, slice, etc.) that is nil. A non-nil nillable
// value fails. A non-nillable type like int or struct will fail as a misuse
// because those types can never be nil.
func (m *Matcher) IsNil() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if m.match = v.IsNil(); !m.match {
			m.t.Error("is not nil")
		}
	case reflect.Invalid:
		// This case happens when a literal nil was passed in.
		m.match = true
	default:
		// Simply prevent a whole class of errors.
		m.match = false
		m.t.Error("is not a nillable type")
	}
	return m
}

// IsNotNil passes when the actual value is a non-nil nillable value (chan,
// func, map, pointer, slice, etc.). A nil nillable value or an untyped nil
// fails. A non-nillable type like int or struct will fail as a misuse
// because those types can never be nil.
func (m *Matcher) IsNotNil() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if m.match = !v.IsNil(); !m.match {
			m.t.Error("is nil")
		}
	case reflect.Invalid:
		// This case happens when a literal nil was passed in.
		m.match = false
		m.t.Error("is a literal nil")
	default:
		// Simply prevent a whole class of errors.
		m.match = false
		m.t.Error("is not a nillable type")
	}
	return m
}

// IsEmpty matches an empty string, slice, map, or array (one with length zero).
// Types that cannot be empty, such as int, struct, and bool, fail.
func (m *Matcher) IsEmpty() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if m.match = v.Len() == 0; !m.match {
			m.t.Error("is not empty")
		}
	default:
		// reflect.Invalid (a literal nil) lands here too; neither it nor
		// types like int, bool, or struct have a length.
		m.match = false
		m.t.Error("is not a length-bearing type")
	}
	return m
}

// IsNotEmpty matches a non-empty string, slice, map, or array (one with a
// length greater than zero). Types that cannot be empty, such as int, struct,
// and bool, fail.
func (m *Matcher) IsNotEmpty() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if m.match = v.Len() > 0; !m.match {
			m.t.Error("is empty")
		}
	default:
		// reflect.Invalid (a literal nil) lands here too; neither it nor
		// types like int, bool, or struct have a length.
		m.match = false
		m.t.Error("is not a length-bearing type")
	}
	return m
}

// IsOk verifies that the actual value is a nil error. It will pass for an
// untyped nil (no error occurred) and fails for a non-nil error. The value
// must be an error (or nil). Any other type is a misuse that fails.
func (m *Matcher) IsOk() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	switch m.actual.(type) {
	case nil:
		// nil represents the absence of an error.
		m.match = true
	case error:
		// A non-nil error value is present, so the result is not ok.
		m.match = false
		if isNilValue(reflect.ValueOf(m.actual)) {
			// A typed nil (for example, (*MyError)(nil)), is non-nil at the
			// interface level but carries no value. We'll call out this use
			// rather than a simple "is not ok: <nil>" to prevent confusion.
			m.t.Errorf("is not ok (typed nil %T — assign the error to a concrete type before asserting IsOk)", m.actual)
		} else {
			m.t.Errorf("is not ok: %v", m.actual)
		}
	default:
		// IsOk only accepts error (or nil) values.
		m.match = false
		m.t.Error("is not an error")
	}
	return m
}

// IsTrue verifies actual value captured in `That()` is `true`.
func (m *Matcher) IsTrue() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	if m.match = v.IsValid() && v.Kind() == reflect.Bool && v.Bool(); !m.match {
		m.t.Error("is not true")
	}
	return m
}

// IsFalse verifies the actual value captured in `That()` is `false`
func (m *Matcher) IsFalse() *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	v := reflect.ValueOf(m.actual)
	if m.match = v.IsValid() && v.Kind() == reflect.Bool && !v.Bool(); !m.match {
		m.t.Error("is not false")
	}
	return m
}

// IsEqualTo verifies that the actual value captured in `That()` is equal to
// the expected value. Signed integers are compared by magnitude among
// themselves, as are unsigned integers, so equal values of differing widths
// (int64 vs. int, or uint64 vs. uint) match. Comparisons across the
// signed/unsigned boundary are not normalized. An int and a uint fall back to
// reflect.DeepEqual and so only match when their boxed types are identical.
// Floats are likewise not normalized. A float32 and a float64 fall back to
// reflect.DeepEqual. Every other type (structs, arrays, maps, pointers, etc.)
// also falls back to reflect.DeepEqual.
func (m *Matcher) IsEqualTo(expected any) *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	m.match = false
	av := reflect.ValueOf(m.actual)
	ev := reflect.ValueOf(expected)

	// There's an edge case where both values are nil. We would argue that the
	// correct way is to check for nil with `IsNil()` first. However, it's not
	// worth failing the test since it's a fairly common use-case.
	if m.actual == nil && expected == nil {
		m.match = true
		return m
	}

	// Both values must be valid.
	if av.IsValid() && ev.IsValid() {
		ak, aErr := basicKind(av)
		ek, eErr := basicKind(ev)
		switch {
		case aErr == nil && eErr == nil && ak == ek && ak == intKind:
			m.match = av.Int() == ev.Int()
		case aErr == nil && eErr == nil && ak == ek && ak == uintKind:
			m.match = av.Uint() == ev.Uint()
		default:
			m.match = reflect.DeepEqual(m.actual, expected)
		}
	}

	if !m.match {
		m.t.Errorf("expected:<[%s]> but was <[%s]>", stringValue(ev), stringValue(av))
	}

	return m
}

// IsGreaterThan matches if the actual value is greater than the expected
// value. Integers, unsigned integers, and floats are compared by magnitude.
// Strings are compared lexically (byte order). Mismatched or unordered types
// (bool, complex, slice, etc.) fail with an error.
func (m *Matcher) IsGreaterThan(expected any) *Matcher {
	m.validateMatcherState()
	m.t.Helper()
	k, err := typeCheck(m.actual, expected)
	if err != nil {
		m.match = false
		m.t.Error(err)
	} else {
		av := reflect.ValueOf(m.actual)
		ev := reflect.ValueOf(expected)
		switch k {
		case floatKind:
			m.match = av.Float() > ev.Float()
		case intKind:
			m.match = av.Int() > ev.Int()
		case uintKind:
			m.match = av.Uint() > ev.Uint()
		case stringKind:
			m.match = av.String() > ev.String()
		default:
			m.match = false
			m.t.Errorf("%s: %T", errInvalidType, m.actual)
			return m
		}

		if !m.match {
			m.t.Errorf("expected: Greater Than <[%s]> but was <[%s]>", stringValue(ev), stringValue(av))
		}
	}

	return m
}

func typeCheck(actual any, expected any) (kind, error) {
	if reflect.TypeOf(actual) == nil {
		return invalidKind, errors.New("actual value was nil")
	}

	if reflect.TypeOf(expected) == nil {
		return invalidKind, errors.New("expected value was nil")
	}

	av := reflect.ValueOf(actual)
	ev := reflect.ValueOf(expected)

	ak, err := basicKind(av)
	if err != nil {
		return invalidKind, errors.New("Actual " + err.Error())
	}

	ek, err := basicKind(ev)
	if err != nil {
		return invalidKind, errors.New("Expected " + err.Error())
	}

	if ak != ek {
		return invalidKind, errIncompatibleTypes
	}

	return ak, nil
}

// stringValue uses reflection to convert a `reflect.Value` to a string for
// use in error messages.
func stringValue(rv reflect.Value) string {
	switch rv.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Complex64:
		c := rv.Complex()
		return "(" + strconv.FormatFloat(real(c), 'g', -1, 32) + "," + strconv.FormatFloat(imag(c), 'g', -1, 32) + ")"
	case reflect.Complex128:
		c := rv.Complex()
		return "(" + strconv.FormatFloat(real(c), 'g', -1, 64) + "," + strconv.FormatFloat(imag(c), 'g', -1, 64) + ")"
	case reflect.String:
		return rv.String()
	case reflect.Invalid:
		// An invalid Value comes from reflect.ValueOf(nil), i.e. a literal nil
		// passed as the expected or actual value.
		return "nil"
	default:
		// Pointers, slices, maps, structs, etc. don't take part in the
		// comparisons above but can still reach an error message; format them
		// generically rather than panicking.
		return fmt.Sprintf("%v", rv)
	}
}

// requireTestingT panics with a useful message if the Matcher was not
// initialized via With(*testing.T). It is a pre-condition guard called by
// That, ThatPanics, and validateMatcherState.
func (m *Matcher) requireTestingT() {
	if m.t == nil {
		panic("Use With(*testing.T) to initialize matcher")
	}
}

// validateMatcherState guards against an uninitialized actual value, in
// addition to the checks performed by requireTestingT.
func (m *Matcher) validateMatcherState() {
	m.requireTestingT()
	if !m.hasActual {
		panic("Use That(any) to specify the actual value")
	}
}

// isNilValue reports whether v holds a nillable kind that is currently nil,
// i.e. a typed nil such as (*MyError)(nil) boxed into an interface.
func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}

// basicKind simplifies the type down to the particular class to which it belongs.
func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	case reflect.Slice:
		return sliceKind, nil
	}
	return invalidKind, errInvalidTypeComparison
}
