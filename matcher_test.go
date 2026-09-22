package assert

import (
	"errors"
	"reflect"
	"testing"
)

func TestWith(t *testing.T) {
	assert := With(new(testing.T))

	if assert == nil {
		t.Error("With returned nil.")
	}
}

func TestMatcher_That(t *testing.T) {
	assert := With(new(testing.T)).That(nil)

	if assert == nil {
		t.Error("That returned nil.")
	}
}

func TestMatcher_That_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("That did not panic")
		}
	}()

	assert := new(Matcher)
	assert.That(nil)
	t.Error("That did not panic.")
}

func TestMatcher_Method_PanicsWithoutThat(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("matcher method did not panic without That")
		}
	}()

	With(new(testing.T)).IsTrue()
	t.Error("matcher method did not panic without That.")
}

func TestMatcher_IsNil_ShouldPassWithNil(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsNil()

	if assert == nil {
		t.Error("IsNil returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsNil matcher failed.")
	}
}

func TestMatcher_IsNil_ShouldPassWithNilPointer(t *testing.T) {
	assert := With(new(testing.T)).That((*string)(nil)).IsNil()

	if assert == nil {
		t.Error("IsNil returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsNil matcher failed.")
	}
}

func TestMatcher_IsNil_ShouldFailWithInt(t *testing.T) {
	assert := With(new(testing.T)).That(0).IsNil()

	if assert == nil {
		t.Error("IsNotNil returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsNotNil matcher failed.")
	}
}

func TestMatcher_IsNil_ShouldFailWithNonNilPointer(t *testing.T) {
	assert := With(new(testing.T)).That(new(int)).IsNil()

	if assert.match == true {
		t.Error("IsNil should fail for a non-nil pointer")
	}
}

func TestMatcher_IsNotNil_ShouldFailWithNil(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsNotNil()

	if assert == nil {
		t.Error("IsNotNil returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsNotNil matcher failed.")
	}
}

func TestMatcher_IsNotNil_ShouldFailWithInt(t *testing.T) {
	assert := With(new(testing.T)).That(0).IsNotNil()

	if assert == nil {
		t.Error("IsNotNil returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsNotNil matcher failed.")
	}
}

func TestMatcher_IsNotNil_ShouldPassWithSlice(t *testing.T) {
	assert := With(new(testing.T)).That(make([]byte, 0)).IsNotNil()

	if assert == nil {
		t.Error("IsNotNil returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsNotNil matcher failed.")
	}
}

func TestMatcher_IsNotNil_ShouldPassWithObject(t *testing.T) {
	assert := With(new(testing.T)).That(new(Matcher)).IsNotNil()

	if assert == nil {
		t.Error("IsNotNil returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsNotNil matcher failed.")
	}
}

func TestMatcher_IsNotNil_ShouldFailWithNilPointer(t *testing.T) {
	assert := With(new(testing.T)).That((*string)(nil)).IsNotNil()

	if assert.match == true {
		t.Error("IsNotNil should fail for a nil pointer")
	}
}

func TestMatcher_IsEmpty_ShouldPassWithEmptyString(t *testing.T) {
	assert := With(new(testing.T)).That("").IsEmpty()

	if assert == nil {
		t.Error("IsEmpty returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEmpty matcher failed.")
	}
}

func TestMatcher_IsEmpty_ShouldFailWithString(t *testing.T) {
	assert := With(new(testing.T)).That("abc").IsEmpty()

	if assert == nil {
		t.Error("IsEmpty returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEmpty matcher failed.")
	}
}

func TestMatcher_IsEmpty_ShouldFailWithNonLengthBearingType(t *testing.T) {
	assert := With(new(testing.T)).That(42).IsEmpty()

	if assert.match == true {
		t.Error("IsEmpty should fail for a type with no length")
	}
}

func TestMatcher_IsNotEmpty_ShouldFailWithNonLengthBearingType(t *testing.T) {
	assert := With(new(testing.T)).That(42).IsNotEmpty()

	if assert.match == true {
		t.Error("IsNotEmpty should fail for a type with no length")
	}
}

func TestMatcher_IsNotEmpty_ShouldPassWithString(t *testing.T) {
	assert := With(new(testing.T)).That("abc").IsNotEmpty()

	if assert == nil {
		t.Error("IsNotEmpty returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsNotEmpty matcher failed.")
	}
}

func TestMatcher_IsNotEmpty_ShouldFailWithEmptyString(t *testing.T) {
	assert := With(new(testing.T)).That("").IsNotEmpty()

	if assert == nil {
		t.Error("IsNotEmpty returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsNotEmpty matcher failed.")
	}
}

func TestMatcher_IsEmpty_ShouldPassWithEmptyContainers(t *testing.T) {
	cases := map[string]any{
		"nil slice":   []int(nil),
		"empty slice": []int{},
		"nil map":     map[string]int(nil),
		"empty map":   map[string]int{},
		"empty array": [0]int{},
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			assert := With(new(testing.T)).That(value).IsEmpty()

			if assert == nil {
				t.Error("IsEmpty returned nil")
				return
			}

			if assert.match == false {
				t.Error("IsEmpty matcher failed.")
			}
		})
	}
}

func TestMatcher_IsEmpty_ShouldFailWithNonEmptyContainers(t *testing.T) {
	cases := map[string]any{
		"slice": []int{1},
		"map":   map[string]int{"a": 1},
		"array": [1]int{1},
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			assert := With(new(testing.T)).That(value).IsEmpty()

			if assert == nil {
				t.Error("IsEmpty returned nil")
				return
			}

			if assert.match == true {
				t.Error("IsEmpty matcher failed.")
			}
		})
	}
}

func TestMatcher_IsNotEmpty_ShouldPassWithNonEmptyContainers(t *testing.T) {
	cases := map[string]any{
		"slice": []int{1},
		"map":   map[string]int{"a": 1},
		"array": [1]int{1},
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			assert := With(new(testing.T)).That(value).IsNotEmpty()

			if assert == nil {
				t.Error("IsNotEmpty returned nil")
				return
			}

			if assert.match == false {
				t.Error("IsNotEmpty matcher failed.")
			}
		})
	}
}

func TestMatcher_IsNotEmpty_ShouldFailWithEmptyContainers(t *testing.T) {
	cases := map[string]any{
		"nil slice":   []int(nil),
		"empty slice": []int{},
		"nil map":     map[string]int(nil),
		"empty map":   map[string]int{},
		"empty array": [0]int{},
	}
	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			assert := With(new(testing.T)).That(value).IsNotEmpty()

			if assert == nil {
				t.Error("IsNotEmpty returned nil")
				return
			}

			if assert.match == true {
				t.Error("IsNotEmpty matcher failed.")
			}
		})
	}
}

func TestMatcher_IsOk_ShouldFailWithError(t *testing.T) {
	err := errors.New("test")
	assert := With(new(testing.T)).That(err).IsOk()

	if assert == nil {
		t.Error("IsOk returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsOk matcher failed.")
	}
}

func TestMatcher_IsOk_ShouldFailWithNonError(t *testing.T) {
	assert := With(new(testing.T)).That("not an error").IsOk()

	if assert.match == true {
		t.Error("IsOk should reject a non-error value")
	}
}

type customError struct{}

func (*customError) Error() string { return "custom error" }

func TestMatcher_IsOk_ShouldFailWithTypedNilError(t *testing.T) {
	// A typed nil that implements error is non-nil at the interface level, so
	// IsOk must report it as a failure rather than treating it as nil.
	var err error = (*customError)(nil)
	assert := With(new(testing.T)).That(err).IsOk()

	if assert.match == true {
		t.Error("IsOk should fail for a typed-nil error")
	}
}

type valueError struct{}

func (valueError) Error() string { return "value error" }

func TestMatcher_IsOk_ShouldFailWithValueError(t *testing.T) {
	// A value-type error (non-nillable kind) is never a typed nil, so it takes
	// the plain failure path.
	var err error = valueError{}
	assert := With(new(testing.T)).That(err).IsOk()

	if assert.match == true {
		t.Error("IsOk should fail for a value-type error")
	}
}

func TestMatcher_IsOk_ShouldPassWithNil(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsOk()

	if assert == nil {
		t.Error("IsOk returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsOk matcher failed.")
	}
}

func TestMatcher_IsTrue_ShouldPassWithTrue(t *testing.T) {
	assert := With(new(testing.T)).That(true).IsTrue()

	if assert == nil {
		t.Error("IsTrue returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsTrue matcher failed.")
	}
}

func TestMatcher_IsTrue_ShouldFailWithFalse(t *testing.T) {
	assert := With(new(testing.T)).That(false).IsTrue()

	if assert == nil {
		t.Error("IsTrue returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsTrue matcher failed.")
	}
}

func TestMatcher_IsFalse_ShouldPassWithFalse(t *testing.T) {
	assert := With(new(testing.T)).That(false).IsFalse()

	if assert == nil {
		t.Error("IsFalse returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsFalse matcher failed.")
	}
}

func TestMatcher_IsFalse_ShouldFailWithTrue(t *testing.T) {
	assert := With(new(testing.T)).That(true).IsFalse()

	if assert == nil {
		t.Error("IsFalse returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsFalse matcher failed.")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithInvalidActualType(t *testing.T) {
	assert := With(new(testing.T)).That(new(testing.T)).IsEqualTo("abc")

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithInvalidExpectedType(t *testing.T) {
	assert := With(new(testing.T)).That("abc").IsEqualTo(new(testing.T))

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithUnequalValues(t *testing.T) {
	assert := With(new(testing.T)).That("abc").IsEqualTo("def")

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithNil(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsEqualTo(nil)

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithBool(t *testing.T) {
	assert := With(new(testing.T)).That(true).IsEqualTo(true)

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithComplex(t *testing.T) {
	assert := With(new(testing.T)).That(complex(1.0, 1.0)).IsEqualTo(complex(1.0, 1.0))

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithFloat(t *testing.T) {
	assert := With(new(testing.T)).That(3.14159).IsEqualTo(3.14159)

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithInt(t *testing.T) {
	assert := With(new(testing.T)).That(int64(-128)).IsEqualTo(-128)

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithString(t *testing.T) {
	assert := With(new(testing.T)).That("The quick brown fox jumps over the lazy dog").IsEqualTo("The quick brown fox jumps over the lazy dog")

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithUint(t *testing.T) {
	assert := With(new(testing.T)).That(uint(1073741824)).IsEqualTo(uint(1073741824))

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithArray(t *testing.T) {
	actual := []byte{'t', 'e', 's', 't'}
	expect := []byte{'t', 'e', 's', 't'}

	assert := With(new(testing.T)).That(actual).IsEqualTo(expect)
	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo array matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithStruct(t *testing.T) {
	type point struct {
		X, Y int
	}

	assert := With(new(testing.T)).That(point{1, 2}).IsEqualTo(point{1, 2})
	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo struct matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithUnequalStruct(t *testing.T) {
	type point struct {
		X, Y int
	}

	assert := With(new(testing.T)).That(point{1, 2}).IsEqualTo(point{3, 4})
	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo struct matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldPassWithMap(t *testing.T) {
	actual := map[string]int{"a": 1, "b": 2}
	expect := map[string]int{"a": 1, "b": 2}

	assert := With(new(testing.T)).That(actual).IsEqualTo(expect)
	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsEqualTo map matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithUnequalMap(t *testing.T) {
	actual := map[string]int{"a": 1, "b": 2}
	expect := map[string]int{"a": 1, "b": 3}

	assert := With(new(testing.T)).That(actual).IsEqualTo(expect)
	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo map matcher failed")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithDifferentTypes(t *testing.T) {
	assert := With(new(testing.T)).That(true).IsEqualTo(1.0)

	if assert == nil {
		t.Error("IsEqualTo returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsEqualTo matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldPassWithFloat(t *testing.T) {
	assert := With(new(testing.T)).That(3.14159).IsGreaterThan(3.14158)

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithFloat(t *testing.T) {
	assert := With(new(testing.T)).That(3.14158).IsGreaterThan(3.14159)

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldPassWithInt(t *testing.T) {
	assert := With(new(testing.T)).That(1073741824).IsGreaterThan(-1073741824)

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithInt(t *testing.T) {
	assert := With(new(testing.T)).That(-1073741824).IsGreaterThan(1073741824)

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldPassWithUint(t *testing.T) {
	assert := With(new(testing.T)).That(uint(1073741824)).IsGreaterThan(uint(1073741823))

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == false {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithUint(t *testing.T) {
	assert := With(new(testing.T)).That(uint(1073741823)).IsGreaterThan(uint(1073741824))

	if assert == nil {
		t.Error("IsGreaterThan returned nil")
		return
	}

	if assert.match == true {
		t.Error("IsGreaterThan matcher failed")
	}
}

func TestMatcher_IsGreaterThan_ShouldPassWithString(t *testing.T) {
	assert := With(new(testing.T)).That("banana").IsGreaterThan("apple")

	if assert.match == false {
		t.Error("IsGreaterThan should pass when the string is lexically greater")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithString(t *testing.T) {
	assert := With(new(testing.T)).That("apple").IsGreaterThan("banana")

	if assert.match == true {
		t.Error("IsGreaterThan should fail when the string is lexically smaller")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithUnorderedType(t *testing.T) {
	assert := With(new(testing.T)).That(true).IsGreaterThan(false)

	if assert.match == true {
		t.Error("IsGreaterThan should fail for an unordered type")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithComplex(t *testing.T) {
	assert := With(new(testing.T)).That(complex(1, 2)).IsGreaterThan(complex(3, 4))

	if assert.match == true {
		t.Error("IsGreaterThan should fail for an unordered complex type")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithSlice(t *testing.T) {
	assert := With(new(testing.T)).That([]int{1}).IsGreaterThan([]int{2})

	if assert.match == true {
		t.Error("IsGreaterThan should fail for an unordered slice type")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithMismatchedTypes(t *testing.T) {
	assert := With(new(testing.T)).That(1).IsGreaterThan("a")

	if assert.match == true {
		t.Error("IsGreaterThan should fail for mismatched types")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithNonComparableType(t *testing.T) {
	assert := With(new(testing.T)).That(new(int)).IsGreaterThan(new(int))

	if assert.match == true {
		t.Error("IsGreaterThan should fail for a non-comparable type")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithNonComparableExpected(t *testing.T) {
	assert := With(new(testing.T)).That(1).IsGreaterThan(new(int))

	if assert.match == true {
		t.Error("IsGreaterThan should fail for a non-comparable expected type")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithNil(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsGreaterThan(1)

	if assert.match == true {
		t.Error("IsGreaterThan should fail for a nil value")
	}
}

func TestMatcher_IsGreaterThan_ShouldFailWithNilExpected(t *testing.T) {
	assert := With(new(testing.T)).That(1).IsGreaterThan(nil)

	if assert.match == true {
		t.Error("IsGreaterThan should fail for a nil expected value")
	}
}

func TestMatcher_ThatPanics_ShouldPassWithPanic(t *testing.T) {
	assert := With(new(testing.T))
	p := func() {
		panic("Panic! at the Disco")
	}

	assert.ThatPanics(p)

	if assert.match == false {
		t.Error("ThatPanics matcher failed.")
	}
}

func TestMatcher_ThatPanics_ShouldFailWithoutPanic(t *testing.T) {
	assert := With(new(testing.T))

	p := func() {
		// Do nothing
	}

	assert.ThatPanics(p)

	if assert.match == true {
		t.Error("ThatPanics matcher failed.")
	}
}

// A typed nil value (e.g. a nil pointer or slice) compared against the literal
// nil must report a normal failure, not panic while formatting the message.
func TestMatcher_IsEqualTo_ShouldFailWithTypedNilPointer(t *testing.T) {
	assert := With(new(testing.T)).That((*int)(nil)).IsEqualTo(nil)

	if assert.match == true {
		t.Error("typed nil pointer should not equal literal nil")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithTypedNilSlice(t *testing.T) {
	assert := With(new(testing.T)).That([]int(nil)).IsEqualTo(nil)

	if assert.match == true {
		t.Error("typed nil slice should not equal literal nil")
	}
}

func TestMatcher_IsEqualTo_ShouldFailWithNilActual(t *testing.T) {
	assert := With(new(testing.T)).That(nil).IsEqualTo(42)

	if assert.match == true {
		t.Error("literal nil should not equal 42")
	}
}

// stringValue must produce a string for every reflect.Kind it can be handed
// while building an error message, including the literal-nil (Invalid) case
// and the generic default path for pointers, slices, and the like.
func TestStringValue(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"bool", true, "true"},
		{"int", int64(-128), "-128"},
		{"uint", uint(42), "42"},
		{"float32", float32(1.5), "1.5"},
		{"float64", 3.14, "3.14"},
		{"complex64", complex64(complex(1, 2)), "(1,2)"},
		{"complex128", complex(1, 2), "(1,2)"},
		{"string", "hello", "hello"},
		{"nil literal", nil, "nil"},
		{"nil pointer", (*int)(nil), "<nil>"},
		{"slice", []int{1, 2}, "[1 2]"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := stringValue(reflect.ValueOf(c.in)); got != c.want {
				t.Errorf("stringValue(%#v) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
