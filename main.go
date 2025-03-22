package main

import (
	"bufio"
	"fmt"
	"time"
	"os"
	"strings"

	"github.com/mrinalxdev/go-rep/db"
)

func main() {
    database, err := db.NewDB("database.log")
    if err != nil {
        fmt.Printf("Error creating DB: %v\n", err)
        return
    }
    defer database.Close()

    fmt.Println("\nSimple Key-Value DB with Replay System")
    fmt.Println("Commands: set <key> <value>, get <key>, delete <key>, list, log, exit")

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        input := strings.TrimSpace(scanner.Text())
        if input == "" {
            continue
        }

        parts := strings.Split(input, " ")
        command := strings.ToLower(parts[0])

        switch command {
        case "set":
            if len(parts) != 3 {
                fmt.Println("Usage: set <key> <value>")
                continue
            }
            err := database.Set(parts[1], parts[2])
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            } else {
                fmt.Println("OK")
            }

        case "get":
            if len(parts) != 2 {
                fmt.Println("Usage: get <key>")
                continue
            }
            value, exists := database.Get(parts[1])
            if exists {
                fmt.Printf("%s\n", value)
            } else {
                fmt.Println("Not found")
            }

        case "delete":
            if len(parts) != 2 {
                fmt.Println("Usage: delete <key>")
                continue
            }
            err := database.Delete(parts[1])
            if err != nil {
                fmt.Printf("Error: %v\n", err)
            } else {
                fmt.Println("OK")
            }

        case "list":
            database.List() // Use the new public method

        case "log":
            entries, err := database.ListLogEntries()
            if err != nil {
                fmt.Printf("Error reading log: %v\n", err)
                continue
            }
            fmt.Println("Log entries:")
            for _, entry := range entries {
                fmt.Printf("  %s %s %s %s\n",
                    entry.Timestamp.Format(time.RFC3339),
                    entry.Command,
                    entry.Key,
                    entry.Value)
            }

        case "exit":
            fmt.Println("Goodbye!")
            return

        default:
            fmt.Println("Unknown command. Available: set, get, delete, list, log, exit")
        }
    }
}
