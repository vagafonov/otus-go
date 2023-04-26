package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here

	// test with slash and withous

	t.Run("additional tests", func(t *testing.T) {
		result, err := ReadDir("./testdata/reader/")
		require.NoError(t, err)

		expected := make(Environment, 5)
		expected["TABEND"] = EnvValue{
			Value:      "with end tab",
			NeedRemove: false,
		}

		expected["TABBEGIN"] = EnvValue{
			Value:      "	with begin tab",
			NeedRemove: false,
		}

		require.Equal(t, expected, result)
	})

	t.Run("file name without slash", func(t *testing.T) {
		_, err := ReadDir("./testdata/reader")
		require.NoError(t, err)
	})
}
