package utilities

import (
	"io"
	"os"
)

func SaveUploadedFile(file io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}

func ReadFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}