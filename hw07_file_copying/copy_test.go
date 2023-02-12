package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := os.Mkdir("tests_output", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(m.Run())
}

func TestCopy(t *testing.T) {
	t.Run("Error same path", func(t *testing.T) {
		result := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
		require.ErrorContains(t, result, "file paths are the same")
	})

	t.Run("Error open file", func(t *testing.T) {
		result := Copy("undefined_file.txt", getTestFileName(), 20, 10)
		require.ErrorContains(t, result, "open undefined_file.txt: no such file or directory")
	})

	t.Run("Error get stat file", func(t *testing.T) {
		result := Copy("/dev/urandom", getTestFileName(), 20, 10)
		require.Error(t, result, ErrUnsupportedFile.Error())
	})

	t.Run("Error offset exceed", func(t *testing.T) {
		result := Copy("testdata/input.txt", getTestFileName(), 2000000, 10)
		require.Error(t, result, ErrOffsetExceedsFileSize.Error())
	})

	/*
		t.Run("Error open file", func(t *testing.T) {
			result := Copy("undefined_file.txt", "out2.txt", 0, 0)

			expectedError := errors.New("no such file or directory")
			expectedError = errors2.Wrap(expectedError, "open undefined_file.txt")

			require.ErrorIs(t, result, expectedError)
		})
	*/

	t.Cleanup(func() {
		if err := os.RemoveAll("tests_output"); err != nil {
			fmt.Println(err)
		}
	})
}

func getTestFileName() string {
	return "tests_output/out-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".txt"
}
