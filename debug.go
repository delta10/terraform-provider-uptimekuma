//go:build debug
// +build debug

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/j0r15/uptime-kuma-terraform-provider/internal/provider"
)

func main() {
	// Create client
	client, err := provider.NewClient("http://localhost:3001", "admin", "cF96H*L9LA3*HiWhx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	// Wait a bit for initial data to load
	fmt.Println("Waiting for initial data to load...")
	time.Sleep(2 * time.Second)
	
	// Get all monitors using existing function
	fmt.Println("=== ALL MONITORS ===")
	monitors, err := client.GetMonitors()
	if err != nil {
		log.Printf("Failed to get monitors: %v", err)
	} else {
		fmt.Printf("Found %d monitors:\n", len(monitors))
		for _, monitor := range monitors {
			fmt.Printf("ID: %d, Name: %s, Type: %s, URL: %s, Active: %t\n", 
				monitor.ID, monitor.Name, monitor.Type, monitor.URL, monitor.Active)
		}
	}
	
	fmt.Println("\n=== NOTIFICATION CLIENT CACHE ===")
	// Get all notifications using existing function
	notifications, err := client.GetNotifications()
	if err != nil {
		log.Printf("Failed to get notifications: %v", err)
	} else {
		fmt.Printf("Found %d notifications in client cache:\n", len(notifications))
		for _, notification := range notifications {
			fmt.Printf("ID: %s, Name: %s, Type: %s, IsDefault: %t, Active: %t\n", 
				notification.ID, notification.Name, notification.Type, 
				notification.IsDefault, notification.Active)
		}
	}
}
