package db

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

func (db *DB) writeLog(op Operation) error {
	logEntry := fmt.Sprintf("%s | %s | %s | %s \n", op.Timestamp.Format(time.RFC3339), op.Command, op.Key, op.Value)

	if _, err := db.logWriter.WriteString(logEntry); err != nil {
		return err
	}
	return db.logWriter.Flush()
}

func (db *DB) replayLog() error {
	if _, err := db.logFile.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(db.logFile)
	for scanner.Scan(){
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		timestamp, command, key := parts[0], parts[1], parts[2]
		value := ""
		if len(parts) > 3 {
			value = parts[3]
		}

		if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
			continue
		}

		switch command {
		case "SET":
			db.store[key] = value
		case "DELETE":
			delete(db.store, key)
		}
	}

	return scanner.Err()
}
