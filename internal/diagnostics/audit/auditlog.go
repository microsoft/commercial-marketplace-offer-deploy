package audit

import (
	"os"
	"time"

	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
)

const indent = "  "

// audit log for tracking interesting events
type Log interface {
	// appends entry asynchronously
	Append(entry any)
}

type AuditLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Entry     any       `json:"entry"`
}

type appendOnlyFileAuditLog struct {
	filePath string
}

func (a *appendOnlyFileAuditLog) Append(entry any) {
	go func() {
		err := a.write(AuditLogEntry{
			Timestamp: time.Now(),
			Entry:     entry,
		})
		if err != nil {
			log.Errorf("error writing to audit log: %v", err)
		}
	}()
}

func (a *appendOnlyFileAuditLog) write(entry any) error {
	entries, err := a.read()
	if err != nil {
		return err
	}
	entries = append(entries, AuditLogEntry{
		Timestamp: time.Now(),
		Entry:     entry,
	})
	bytes, err := json.MarshalIndent(entries, "", indent)
	if err != nil {
		return err
	}

	err = os.WriteFile(a.filePath, bytes, 0660)
	if err != nil {
		return err
	}

	return nil
}

func (a *appendOnlyFileAuditLog) read() ([]AuditLogEntry, error) {
	entries := []AuditLogEntry{}
	bytes, err := os.ReadFile(a.filePath)

	if err != nil {
		return entries, err
	}

	err = json.Unmarshal(bytes, &entries)
	if err != nil {
		return entries, err
	}
	return entries, nil
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

func NewAppendOnlyFileAuditLog(filePath string) (Log, error) {
	auditLog := &appendOnlyFileAuditLog{
		filePath: filePath,
	}
	err := auditLog.ensureFileExists()
	if err != nil {
		return nil, err
	}

	return auditLog, nil
}
