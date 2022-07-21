package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHas(t *testing.T) {
	t.Run("string array", func(t *testing.T) {
		array := []string{"hello", "world", "welcome"}

		require.True(t, Has(array, "hello"))
		require.True(t, Has(array, "world"))
		require.False(t, Has(array, "golang"))
	})

	t.Run("int array", func(t *testing.T) {
		array := []int{0, 1, 2, 3, 4}

		require.True(t, Has(array, 0))
		require.True(t, Has(array, 3))
		require.False(t, Has(array, 12))
	})
}

func TestFind(t *testing.T) {
	type TestingFind struct {
		foo string
		bar string
	}

	array := make([]*TestingFind, 3)
	array[0] = &TestingFind{
		foo: "foo1",
		bar: "bar1",
	}
	array[1] = &TestingFind{
		foo: "foo2",
		bar: "bar2",
	}
	array[2] = &TestingFind{
		foo: "foo3",
		bar: "bar3",
	}

	el1, ok := Find(array, func(e *TestingFind) bool {
		return e.foo == "foo1"
	})
	require.True(t, ok)
	require.True(t, el1.foo == "foo1")

	el2, ok := Find(array, func(e *TestingFind) bool {
		return e.foo == "foo2"
	})
	require.True(t, ok)
	require.True(t, el2.foo == "foo2")

	el3, ok := Find(array, func(e *TestingFind) bool {
		return e.foo == "golang"
	})
	require.False(t, ok)
	require.Empty(t, el3)
}
