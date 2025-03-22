package main

import (
	"fmt"

	"github.com/mrinalxdev/go-rep/db"
)

func main() {
    database, err := db.NewDB("database.log")
    if err != nil {
        fmt.Printf("Error creating DB: %v\n", err)
        return
    }
    defer database.Close()

    // Example usage
    database.Set("user1", "Alice")
    database.Set("user2", "Bob")
    database.Delete("user1")

    if value, exists := database.Get("user2"); exists {
        fmt.Printf("user2: %s\n", value)
    }
    if _, exists := database.Get("user1"); !exists {
        fmt.Println("user1: not found")
    }
}
