package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	errors2 "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Error same path", func(t *testing.T) {
		result := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
		require.ErrorContains(t, result, "file paths are the same")
	})

	t.Run("Error open file", func(t *testing.T) {
		result := Copy("undefined_file.txt", "out.txt", 20, 10)
		require.ErrorContains(t, result, "open undefined_file.txt: no such file or directory")
	})

	t.Run("Error get stat file", func(t *testing.T) {
		result := Copy("/dev/urandom", "out.txt", 20, 10)
		require.Error(t, result, ErrUnsupportedFile.Error())
	})

	t.Run("Error offset exceed", func(t *testing.T) {
		result := Copy("testdata/input.txt", "out.txt", 2000000, 10)
		require.Error(t, result, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("Error open file", func(t *testing.T) {
		result := Copy("undefined_file.txt", "out2.txt", 0, 0o0)

		expectedError := errors.New("no such file or directory")
		expectedError = errors2.Wrap(expectedError, "open undefined_file.txt")

		require.ErrorIs(t, result, expectedError)
	})

	t.Cleanup(func() {
		if err := os.Remove("out.txt"); err != nil {
			fmt.Println(err)
		}
	})
}
