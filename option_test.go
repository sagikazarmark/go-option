package option

import (
	"errors"
	"fmt"
	"testing"
)

func TestSome(t *testing.T) {
	o := Some("hello")

	if !o.HasValue() {
		t.Error("Some should have a value")
	}

	if o.Value() != "hello" {
		t.Error("expected value hello, got:", o.Value())
	}
}

func ExampleSome() {
	o := Some("hello")

	// Note: you are not supposed to call these methods directly.
	// Please take a look at the rest of the functions in the package.
	fmt.Println(o.HasValue())
	fmt.Println(o.Value())

	// Output:
	// true
	// hello
}

func TestIsSome(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		if !IsSome(Some("hello")) {
			t.Error("Some should identify as Some")
		}
	})

	t.Run("None", func(t *testing.T) {
		if IsSome(None[string]()) {
			t.Error("None should not identify as Some")
		}
	})
}

func ExampleIsSome() {
	s := Some("hello")
	n := None[string]()

	fmt.Println(IsSome(s))
	fmt.Println(IsSome(n))

	// Output:
	// true
	// false
}

func TestNone(t *testing.T) {
	o := None[string]()

	if o.HasValue() {
		t.Error("None should not have a value")
	}

	if o.Value() != "" {
		t.Error("None should hold the default value of the type, got:", o.Value())
	}
}

func ExampleNone() {
	o := None[string]()

	// Note: you are not supposed to call these methods directly.
	// Please take a look at the rest of the functions in the package.
	fmt.Println(o.HasValue())
	fmt.Println(o.Value())

	// Output:
	// false
	//
}

func TestIsNone(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		if IsNone(Some("hello")) {
			t.Error("Some should not identify as None")
		}
	})

	t.Run("None", func(t *testing.T) {
		if !IsNone(None[string]()) {
			t.Error("None should identify as None")
		}
	})
}

func ExampleIsNone() {
	s := Some("hello")
	n := None[string]()

	fmt.Println(IsNone(s))
	fmt.Println(IsNone(n))

	// Output:
	// false
	// true
}

func TestUnwrap(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		v := Unwrap(Some("hello"))

		if v != "hello" {
			t.Error("expected Unwrap to return the contained value, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		defer func() {
			v := recover()

			if v == nil {
				t.Error("expected Unwrap to panic on None")
			}
		}()

		Unwrap(None[string]())
	})
}

func ExampleUnwrap() {
	s := Some("hello")

	fmt.Println(Unwrap(s))
	// fmt.Println(Unwrap(None[string]())) // This would panic

	// Output:
	// hello
}

func TestUnwrapOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		v := UnwrapOr(Some("hello"), "world")

		if v != "hello" {
			t.Error("expected UnwrapOr to return the contained value, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		v := UnwrapOr(None[string](), "world")

		if v != "world" {
			t.Error("expected UnwrapOr to return the provided value, got:", v)
		}
	})
}

func ExampleUnwrapOr() {
	s := Some("hello")
	n := None[string]()

	fmt.Println(UnwrapOr(s, "world"))
	fmt.Println(UnwrapOr(n, "world"))

	// Output:
	// hello
	// world
}

func TestUnwrapOrDefault(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		v := UnwrapOrDefault(Some("hello"))

		if v != "hello" {
			t.Error("expected UnwrapOrDefault to return the contained value, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		v := UnwrapOrDefault(None[string]())

		if v != "" {
			t.Error("expected UnwrapOrDefault to return the type default value, got:", v)
		}
	})
}

func ExampleUnwrapOrDefault() {
	s := Some("hello")
	n := None[string]()

	fmt.Println(UnwrapOrDefault(s))
	fmt.Println(UnwrapOrDefault(n))

	// Output:
	// hello
	//
}

func TestUnwrapOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		v := UnwrapOrElse(Some("hello"), func() string { return "world" })

		if v != "hello" {
			t.Error("expected UnwrapOrElse to return the contained value, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		v := UnwrapOrElse(None[string](), func() string { return "world" })

		if v != "world" {
			t.Error("expected UnwrapOrDefault to return the computed value, got:", v)
		}
	})
}

func ExampleUnwrapOrElse() {
	s := Some("hello")
	n := None[string]()

	fmt.Println(UnwrapOrElse(s, func() string { return "world" }))
	fmt.Println(UnwrapOrElse(n, func() string { return "world" }))

	// Output:
	// hello
	// world
}

func TestMap(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")

		v := Map(o, func(v string) int { return len(v) })

		if !Equals(v, Some(5)) {
			t.Error("expected Map to return Some(5), got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := Map(o, func(v string) int { return len(v) })

		if !IsNone(v) {
			t.Error("expected Map to return None, got:", v)
		}
	})
}

func TestTryMap(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := Some("hello")

			v, err := TryMap(o, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if !Equals(v, Some(5)) {
				t.Error("expected TryMap to return Some(5), got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := Some("hello")

			e := errors.New("error")

			v, err := TryMap(o, func(v string) (int, error) { return len(v), e })
			if err == nil {
				t.Fatal("expected error")
			}

			if err != e {
				t.Error("expected TryMap to return error, got:", err)
			}

			if !IsNone(v) {
				t.Error("expected TryMap to return None, got:", v)
			}
		})
	})

	t.Run("None", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := None[string]()

			v, err := TryMap(o, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if !IsNone(v) {
				t.Error("expected TryMap to return None, got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := None[string]()

			e := errors.New("error")

			v, err := TryMap(o, func(v string) (int, error) { return len(v), e })
			if err != nil {
				t.Fatal(err)
			}

			if !IsNone(v) {
				t.Error("expected TryMap to return None, got:", v)
			}
		})
	})
}

func TestMapOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")

		v := MapOr(o, 10, func(v string) int { return len(v) })

		if v != 5 {
			t.Error("expected MapOr to return 5, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := MapOr(o, 10, func(v string) int { return len(v) })

		if v != 10 {
			t.Error("expected MapOr to return 10, got:", v)
		}
	})
}

func TestTryMapOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := Some("hello")

			v, err := TryMapOr(o, 10, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if v != 5 {
				t.Error("expected TryMapOr to return 5, got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := Some("hello")

			e := errors.New("error")

			v, err := TryMapOr(o, 10, func(v string) (int, error) { return len(v), e })
			if err == nil {
				t.Fatal("expected error")
			}

			if err != e {
				t.Error("expected TryMapOr to return error, got:", err)
			}

			if v != 5 {
				t.Error("expected TryMapOr to return 5, got:", v)
			}
		})
	})

	t.Run("None", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := None[string]()

			v, err := TryMapOr(o, 10, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if v != 10 {
				t.Error("expected TryMapOr to return 10, got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := None[string]()

			e := errors.New("error")

			v, err := TryMapOr(o, 10, func(v string) (int, error) { return len(v), e })
			if err != nil {
				t.Fatal(err)
			}

			if v != 10 {
				t.Error("expected TryMapOr to return 10, got:", v)
			}
		})
	})
}

func TestMapOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")

		v := MapOrElse(o, func() int { return 10 }, func(v string) int { return len(v) })

		if v != 5 {
			t.Error("expected MapOrElse to return 5, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := MapOrElse(o, func() int { return 10 }, func(v string) int { return len(v) })

		if v != 10 {
			t.Error("expected MapOrElse to return 10, got:", v)
		}
	})
}

func TestTryMapOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := Some("hello")

			v, err := TryMapOrElse(o, func() int { return 10 }, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if v != 5 {
				t.Error("expected TryMapOrElse to return 5, got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := Some("hello")

			e := errors.New("error")

			v, err := TryMapOrElse(o, func() int { return 10 }, func(v string) (int, error) { return len(v), e })
			if err == nil {
				t.Fatal("expected error")
			}

			if err != e {
				t.Error("expected TryMapOrElse to return error, got:", err)
			}

			if v != 5 {
				t.Error("expected TryMapOrElse to return 5, got:", v)
			}
		})
	})

	t.Run("None", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			o := None[string]()

			v, err := TryMapOrElse(o, func() int { return 10 }, func(v string) (int, error) { return len(v), nil })
			if err != nil {
				t.Fatal(err)
			}

			if v != 10 {
				t.Error("expected TryMapOrElse to return 10, got:", v)
			}
		})

		t.Run("Error", func(t *testing.T) {
			o := None[string]()

			e := errors.New("error")

			v, err := TryMapOrElse(o, func() int { return 10 }, func(v string) (int, error) { return len(v), e })
			if err != nil {
				t.Fatal(err)
			}

			if v != 10 {
				t.Error("expected TryMapOrElse to return 10, got:", v)
			}
		})
	})
}

func TestAnd(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")
		o2 := Some("world")

		v := And(o, o2)

		if !Equals(v, o2) {
			t.Error("expected And to return Some(\"world\"), got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()
		o2 := Some("world")

		v := And(o, o2)

		if !IsNone(v) {
			t.Error("expected And to return None, got:", v)
		}
	})
}

func TestAndThen(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")

		v := AndThen(o, func(v string) Option[string] {
			r := []rune(v)
			for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return Some(string(r))
		})

		if !Equals(v, Some("olleh")) {
			t.Error("expected AndThen to return Some(\"olleh\"), got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := AndThen(o, func(v string) Option[string] {
			r := []rune(v)
			for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
				r[i], r[j] = r[j], r[i]
			}
			return Some(string(r))
		})

		if !IsNone(v) {
			t.Error("expected AndThen to return None, got:", v)
		}
	})
}

func TestOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")
		o2 := Some("world")

		v := Or(o, o2)

		if !Equals(v, o) {
			t.Error("expected Or to return Some(\"hello\"), got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()
		o2 := Some("world")

		v := Or(o, o2)

		if !Equals(v, o2) {
			t.Error("expected Or to return Some(\"world\"), got:", v)
		}
	})
}

func TestOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")

		v := OrElse(o, func() Option[string] { return Some("world") })

		if !Equals(v, o) {
			t.Error("expected OrElse to return Some(\"hello\"), got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := OrElse(o, func() Option[string] { return Some("world") })

		if !Equals(v, Some("world")) {
			t.Error("expected OrElse to return Some(\"world\"), got:", v)
		}
	})
}

func TestXor(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		o := Some("hello")
		o2 := Some("world")

		v := Xor(o, o2)

		if !IsNone(v) {
			t.Error("expected Xor to return None, got:", v)
		}
	})

	t.Run("None", func(t *testing.T) {
		t.Run("left", func(t *testing.T) {
			o := None[string]()
			o2 := Some("world")

			v := Xor(o, o2)

			if !Equals(v, o2) {
				t.Error("expected Xor to return Some(\"world\"), got:", v)
			}
		})

		t.Run("right", func(t *testing.T) {
			o := Some("hello")
			o2 := None[string]()

			v := Xor(o, o2)

			if !Equals(v, o) {
				t.Error("expected Xor to return Some(\"hello\"), got:", v)
			}
		})
	})
}

func TestFilter(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			o := Some("hello")

			v := Filter(o, func(v string) bool { return true })

			if !Equals(v, o) {
				t.Error("expected Filter to return Some(\"hello\"), got:", v)
			}
		})

		t.Run("false", func(t *testing.T) {
			o := Some("hello")

			v := Filter(o, func(v string) bool { return false })

			if !IsNone(v) {
				t.Error("expected Filter to return None, got:", v)
			}
		})
	})

	t.Run("None", func(t *testing.T) {
		o := None[string]()

		v := Filter(o, func(v string) bool { return true })

		if !IsNone(v) {
			t.Error("expected Filter to return None, got:", v)
		}
	})
}

func TestEquals(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		t.Run("True", func(t *testing.T) {
			o1 := Some("hello")
			o2 := Some("hello")

			if !Equals(o1, o2) {
				t.Error("two Somes holding the same value are expected to be equal")
			}
		})

		t.Run("False", func(t *testing.T) {
			o1 := Some("hello")
			o2 := Some("world")

			if Equals(o1, o2) {
				t.Error("two Somes holding different values are not expected to be equal")
			}
		})
	})

	t.Run("None", func(t *testing.T) {
		o1 := None[string]()
		o2 := None[string]()

		if !Equals(o1, o2) {
			t.Error("two Nones are expected to be equal")
		}
	})

	t.Run("SomeAndNone", func(t *testing.T) {
		o1 := Some("hello")
		o2 := None[string]()

		if Equals(o1, o2) {
			t.Error("a Some and a None should never be equal")
		}
	})
}
