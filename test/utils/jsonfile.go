package testutils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type jsonfile struct {
	path string
}

func newJsonFile(path string) *jsonfile {
	filepath := path
	if !strings.HasSuffix(filepath, ".json") {
		filepath = filepath + ".json"
	}

	return &jsonfile{path: filepath}
}

func (f *jsonfile) read() ([]byte, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	return bytes, nil
}

func (f *jsonfile) unmarshal(v interface{}) error {
	bytes, err := f.read()
	if err != nil {
		return err
	}
	json.Unmarshal(bytes, &v)
	return nil
}

func ReaderFromJsonFile(filename string) (io.Reader, error) {
	file := newJsonFile(filename)
	bytes, err := file.read()
	return strings.NewReader(string(bytes)), err
}

// Gets an new instance of T by reading and unmarsalling the json file
func NewFromJsonFile[T any](filename string) (T, error) {
	file := newJsonFile(filename)
	var v T
	err := file.unmarshal(&v)

	return v, err
}
