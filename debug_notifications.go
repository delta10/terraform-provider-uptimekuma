//go:build test

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/j0r15/terraform-provider-uptimekuma/internal/provider"
)

func main() {
	fmt.Println("Testing notification retrieval...")

	// Create client
	client, err := provider.NewClient("http://localhost:3001", "admin", "cF96H*L9LA3*HiWhx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Connected to Uptime Kuma")
	
	// Wait a bit for the notificationList event to be received
	fmt.Println("⏳ Waiting for initial events...")
	time.Sleep(3 * time.Second)

	// Test getting all notifications
	fmt.Println("📋 Getting all notifications...")
	notifications, err := client.GetNotifications()
	if err != nil {
		log.Printf("❌ Failed to get notifications: %v", err)
	} else {
		fmt.Printf("✅ Found %d notifications:\n", len(notifications))
		for _, notif := range notifications {
			fmt.Printf("  - ID: %d, Name: %s, Type: %s, Active: %t\n", notif.ID, notif.Name, notif.Type, notif.Active)
		}
	}

	// Test getting specific notification by ID 1
	fmt.Println("🔍 Getting notification ID 1...")
	notif, err := client.GetNotification(1)
	if err != nil {
		log.Printf("❌ Failed to get notification ID 1: %v", err)
	} else {
		fmt.Printf("✅ Retrieved notification: ID=%d, Name=%s, Type=%s\n", notif.ID, notif.Name, notif.Type)
	}
}
