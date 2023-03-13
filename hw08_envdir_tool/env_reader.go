package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))
	for _, f := range files {
		if strings.Contains(f.Name(), "=") {
			continue
		}

		file, err := os.Open(strings.TrimSuffix(dir, "/") + "/" + f.Name())
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		envValue := EnvValue{
			Value:      strings.TrimRight(strings.ReplaceAll(scanner.Text(), "\x00", "\n"), " \t"),
			NeedRemove: false,
		}

		if fInfo.Size() == 0 {
			envValue.NeedRemove = true
		}

		env[f.Name()] = envValue
	}

	return env, nil
}
