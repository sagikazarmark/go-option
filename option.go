// Package option provides tools for working with optional values.
//
// It is heavily inspired by the option module in Rust implementing the same functionality:
// https://doc.rust-lang.org/std/option/index.html
package option

// Option represents an optional value.
// It either contains a value or it does not.
//
// An Option that contains a value is often called Some,
// while an Option without a value is called None.
// The terminology comes from Rust's option module:
// https://doc.rust-lang.org/std/option/index.html
//
// Option describes a low-level interface used by the high-level API implemented by this package.
// The methods defined in Option are not supposed to be called directly.
type Option[T any] interface {
	// HasValue returns true if the Option contains a value.
	HasValue() bool

	// Value returns the value (or its default) stored in the Option.
	Value() T
}

// Some returns a new Option that contains a value.
func Some[T any](value T) Option[T] {
	return some[T]{
		value: value,
	}
}

// IsSome returns true if o contains a value.
func IsSome[T any](o Option[T]) bool {
	return o.HasValue()
}

type some[T any] struct {
	value T
}

func (s some[T]) HasValue() bool {
	return true
}

func (s some[T]) Value() T {
	return s.value
}

// None returns a new Option that does not contain a value.
func None[T any]() Option[T] {
	return none[T]{}
}

// IsNone returns true if o does not contain a value.
func IsNone[T any](o Option[T]) bool {
	return !o.HasValue()
}

type none[T any] struct{}

func (none[T]) HasValue() bool {
	return false
}

func (none[T]) Value() T {
	var value T

	return value
}

// Unwrap returns the contained value or panics.
func Unwrap[T any](o Option[T]) T {
	if IsNone(o) {
		panic("option does not contain any value")
	}

	return o.Value()
}

// UnwrapOr returns the contained value (if any) or returns the provided default value.
func UnwrapOr[T any](o Option[T], d T) T {
	if IsNone(o) {
		return d
	}

	return o.Value()
}

// UnwrapOrDefault returns the contained value (if any) or returns the default value of the type.
func UnwrapOrDefault[T any](o Option[T]) T {
	return o.Value()
}

// UnwrapOrElse returns the contained value (if any) or computes it from the provided default function.
func UnwrapOrElse[T any](o Option[T], d func() T) T {
	if IsNone(o) {
		return d()
	}

	return o.Value()
}

// Map applies the provided function to the contained value (if any) or returns a None.
func Map[T any, U any](o Option[T], f func(v T) U) Option[U] {
	if IsNone(o) {
		return None[U]()
	}

	return Some(f(o.Value()))
}

// TryMap applies the provided function to the contained value (if any) or returns a None.
// If the function returns an error, it propagates back (with a None).
func TryMap[T any, U any](o Option[T], f func(v T) (U, error)) (Option[U], error) {
	if IsNone(o) {
		return None[U](), nil
	}

	v, err := f(o.Value())
	if err != nil {
		return None[U](), err
	}

	return Some(v), nil
}

// MapOr applies the provided function to the contained value (if any) or returns the provided default value.
func MapOr[T any, U any](o Option[T], d U, f func(v T) U) U {
	if IsNone(o) {
		return d
	}

	return f(o.Value())
}

// TryMapOr applies the provided function to the contained value (if any) or returns the provided default value.
// If the function returns an error, it propagates back.
func TryMapOr[T any, U any](o Option[T], d U, f func(v T) (U, error)) (U, error) {
	if IsNone(o) {
		return d, nil
	}

	return f(o.Value())
}

// MapOrElse applies the provided function to the contained value (if any) or computes it from the provided default function.
func MapOrElse[T any, U any](o Option[T], d func() U, f func(v T) U) U {
	if IsNone(o) {
		return d()
	}

	return f(o.Value())
}

// TryMapOrElse applies the provided function to the contained value (if any) or computes it from the provided default function.
func TryMapOrElse[T any, U any](o Option[T], d func() U, f func(v T) (U, error)) (U, error) {
	if IsNone(o) {
		return d(), nil
	}

	return f(o.Value())
}

// And returns o2 if o contains a value.
func And[T any](o Option[T], o2 Option[T]) Option[T] {
	if IsNone(o) {
		return None[T]()
	}

	return o2
}

// AndThen applies the provided function to the contained value (if any) and returns the new value or returns a None.
func AndThen[T any](o Option[T], f func(v T) Option[T]) Option[T] {
	if IsNone(o) {
		return None[T]()
	}

	return f(o.Value())
}

// Or returns o if it contains a value, returns o2 otherwise.
func Or[T any](o Option[T], o2 Option[T]) Option[T] {
	if IsNone(o) {
		return o2
	}

	return o
}

// OrElse returns o if it contains a value or returns the result of calling the provided function.
func OrElse[T any](o Option[T], f func() Option[T]) Option[T] {
	if IsNone(o) {
		return f()
	}

	return o
}

// Xor returns o or o2 if exactly one of them contains a value, otherwise returns a None.
func Xor[T any](o Option[T], o2 Option[T]) Option[T] {
	if IsSome(o) && IsNone(o2) {
		return o
	}

	if IsNone(o) && IsSome(o2) {
		return o2
	}

	return None[T]()
}

// Filter returns o if it contains a value and the provided predicate applied to the contained value returns true.
func Filter[T any](o Option[T], f func(T) bool) Option[T] {
	if IsNone(o) {
		return None[T]()
	}

	if !f(o.Value()) {
		return None[T]()
	}

	return o
}

// Equals checks if two values are equal to each other according to the following:
// - Two Nones are always equal
// - Two Somes are equal if their values are equal
func Equals[T comparable](o1 Option[T], o2 Option[T]) bool {
	if IsSome(o1) != IsSome(o2) {
		return false
	}

	// Technically this check is not necessary: a None should always return the type default value.
	// This check ensures to return a correct result even for a non-compliant implementation.
	if IsNone(o1) && IsNone(o2) {
		return true
	}

	return o1.Value() == o2.Value()
}
