package db

import (
    "bufio"
    "fmt"
    // "os"
    "strings"
    "time"
)

func (db *DB) writeLog(op Operation) error {
    logEntry := fmt.Sprintf("%s|%s|%s|%s\n",
        op.Timestamp.Format(time.RFC3339),
        op.Command,
        op.Key,
        op.Value)

    if _, err := db.logWriter.WriteString(logEntry); err != nil {
        return err
    }
    return db.logWriter.Flush()
}

func (db *DB) replayLog() error {
    if _, err := db.logFile.Seek(0, 0); err != nil {
        return err
    }

    fmt.Println("Replaying log entries:")
    scanner := bufio.NewScanner(db.logFile)
    for scanner.Scan() {
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

        fmt.Printf("  %s %s %s %s\n", timestamp, command, key, value)
        switch command {
        case "SET":
            db.store[key] = value
        case "DELETE":
            delete(db.store, key)
        }
    }

    return scanner.Err()
}

// ListLogEntries returns all log entries for display
func (db *DB) ListLogEntries() ([]Operation, error) {
    if _, err := db.logFile.Seek(0, 0); err != nil {
        return nil, err
    }

    var entries []Operation
    scanner := bufio.NewScanner(db.logFile)
    for scanner.Scan() {
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

        ts, err := time.Parse(time.RFC3339, timestamp)
        if err != nil {
            continue
        }

        entries = append(entries, Operation{
            Timestamp: ts,
            Command:   command,
            Key:       key,
            Value:     value,
        })
    }

    return entries, scanner.Err()
}
