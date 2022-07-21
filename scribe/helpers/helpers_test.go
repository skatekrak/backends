package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetIfNotNil(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		var b *bool

		newBool := SetIfNotNil(b, true)
		require.True(t, newBool)

		newBool = SetIfNotNil(b, false)
		require.False(t, newBool)

		b2 := true
		newBool2 := SetIfNotNil(&b2, false)
		require.True(t, newBool2)

		newBool2 = SetIfNotNil(&b2, true)
		require.True(t, newBool2)
	})

	t.Run("string", func(t *testing.T) {
		var str *string

		newStr := SetIfNotNil(str, "hello")
		require.True(t, newStr == "hello")

		newStr = SetIfNotNil(str, "world")
		require.True(t, newStr == "world")

		str2 := "golang"
		newStr2 := SetIfNotNil(&str2, "hello")
		require.True(t, newStr2 == "golang")
		require.False(t, newStr2 == "hello")
	})

	t.Run("interface", func(t *testing.T) {
		type TestingInterface struct {
			foo string
			bar string
		}

		var e *TestingInterface

		newEl := SetIfNotNil(e, TestingInterface{foo: "foo", bar: "bar"})
		require.True(t, newEl.foo == "foo")

		e2 := &TestingInterface{
			foo: "hello",
			bar: "golang",
		}

		newEl2 := SetIfNotNil(e2, TestingInterface{foo: "foo", bar: "bar"})
		require.True(t, newEl2.foo == "hello")
		require.False(t, newEl2.foo == "foo")
	})
}
