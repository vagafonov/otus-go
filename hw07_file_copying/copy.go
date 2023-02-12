package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSamePath              = errors.New("file paths are the same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrSamePath
	}

	inputFile, openErr := os.Open(fromPath)
	if openErr != nil {
		return openErr
	}

	fStat, statErr := inputFile.Stat()
	if statErr != nil {
		return statErr
	}

	if fStat.Size() == 0 {
		return ErrOffsetExceedsFileSize
	}

	if limit < fStat.Size() && offset > fStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fStat.Size()
	}

	_, err := inputFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	outFile, _ := os.Create(toPath)
	var count int64 = 1
	bytes := 1
	bar := pb.StartNew(int(limit) / bytes)

	for {
		if count > limit {
			break
		}

		read, err := io.CopyN(outFile, inputFile, int64(bytes))
		count += read
		if err != nil {
			break
		}
		bar.Increment()
	}
	bar.Finish()

	return errors.New("Done")
}
