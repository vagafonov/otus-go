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

	_, err := inputFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	reader := io.LimitReader(inputFile, limit)
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(reader)
	_, err = io.Copy(outFile, barReader)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}
