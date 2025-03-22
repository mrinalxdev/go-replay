package db

import (
    "bufio"
    "os"
	"time"
    "sync"
)

// DB represents the database
type DB struct {
    store     map[string]string
    logFile   *os.File
    logWriter *bufio.Writer
    mu        sync.RWMutex
}

// NewDB creates a new database instance
func NewDB(logPath string) (*DB, error) {
    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        return nil, err
    }

    db := &DB{
        store:     make(map[string]string),
        logFile:   file,
        logWriter: bufio.NewWriter(file),
    }

    if err := db.replayLog(); err != nil {
        return nil, err
    }

    return db, nil
}

// Set stores a key-value pair
func (db *DB) Set(key, value string) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    op := Operation{
        Timestamp: time.Now(),
        Command:   "SET",
        Key:       key,
        Value:     value,
    }

    if err := db.writeLog(op); err != nil {
        return err
    }

    db.store[key] = value
    return nil
}

// Get retrieves a value by key
func (db *DB) Get(key string) (string, bool) {
    db.mu.RLock()
    defer db.mu.RUnlock() // This is correct for read operations

    value, exists := db.store[key]
    return value, exists
}

// Delete removes a key-value pair
func (db *DB) Delete(key string) error {
    db.mu.Lock() // Changed from RLock to Lock since we're modifying
    defer db.mu.Unlock()

    op := Operation{
        Timestamp: time.Now(),
        Command:   "DELETE",
        Key:       key,
    }

    if err := db.writeLog(op); err != nil {
        return err
    }

    delete(db.store, key)
    return nil
}

// Close cleans up resources
func (db *DB) Close() error {
    db.mu.Lock()
    defer db.mu.Unlock()

    db.logWriter.Flush()
    return db.logFile.Close()
}
