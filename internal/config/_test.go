package config

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	path := "./testdata/test.env"
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		file.Close()
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("TEST_CONFIG_ENTRY=filevalue \n")
	os.Clearenv()

	m.Run()
}
