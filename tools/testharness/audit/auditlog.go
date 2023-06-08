package audit

import (
	"os"
	"time"

	"github.com/goccy/go-json"
)

const indent = "  "

// audit log for tracking interesting events
type Log interface {
	Append(entry any) error
}

type AuditLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Entry     any       `json:"entry"`
}

type appendOnlyFileAuditLog struct {
	filePath string
}

func (a *appendOnlyFileAuditLog) Append(entry any) error {
	file, err := a.open()
	if err != nil {
		return err
	}
	defer file.Close()

	return a.write(AuditLogEntry{
		Timestamp: time.Now(),
		Entry:     entry,
	}, file)
}

func (*appendOnlyFileAuditLog) write(entry any, file *os.File) error {
	bytes, err := json.MarshalIndent(entry, "", indent)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	file.Write([]byte("\n"))

	if err != nil {
		return err
	}
	return nil
}

// function that opens the file and returns os.File
func (a *appendOnlyFileAuditLog) open() (*os.File, error) {
	return os.OpenFile(a.filePath, os.O_RDWR|os.O_APPEND, 0660)
}

func (a *appendOnlyFileAuditLog) ensureFileExists() error {
	if _, err := os.Stat(a.filePath); os.IsNotExist(err) {
		file, err := os.Create(a.filePath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func NewAppendOnlyFileAuditLog(filePath string) Log {
	auditLog := &appendOnlyFileAuditLog{
		filePath: filePath,
	}
	err := auditLog.ensureFileExists()
	if err != nil {
		panic(err)
	}

	return auditLog
}
