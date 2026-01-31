package services

import (
	"os"
	"path/filepath"
	"strings"
)

type LogService interface {
	GetAuditLog() ([]string, error)
	GetAuthLog() ([]string, error)
}

type logService struct {
	logDir string
}

func NewLogService(logDir string) LogService {
	return &logService{logDir: logDir}
}

func (s *logService) GetAuditLog() ([]string, error) {
	return s.readLogFile("audit.log")
}

func (s *logService) GetAuthLog() ([]string, error) {
	return s.readLogFile("auth.log")
}

func (s *logService) readLogFile(filename string) ([]string, error) {
	path := filepath.Join(s.logDir, filename)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []string{}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Reverse the lines so newest logs are first
	var result []string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}
