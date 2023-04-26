package log

import(
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
)

type FileHook struct {
	fileName string
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	if _, ok := entry.Data["message"]; !ok {
		entry.Data["message"] = entry.Message
	}
	if _, err := os.Stat(hook.fileName); os.IsNotExist(err) {
		_, err := os.Create(hook.fileName)
		if err != nil {
			return err
		}
	} else {
		f, _ := os.OpenFile(hook.fileName, os.O_RDWR|os.O_APPEND, 0660);

		if message, ok := entry.Data["message"]; ok {
			f.Write([]byte(fmt.Sprintf("%v\n", message)))
		}
	}

	return nil
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}