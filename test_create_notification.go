//go:build test

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/j0r15/terraform-provider-uptimekuma/internal/provider"
)

func main() {
	fmt.Println("Testing notification creation...")

	// Create client
	client, err := provider.NewClient("http://localhost:3001", "admin", "cF96H*L9LA3*HiWhx")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Connected to Uptime Kuma")

	// Wait a bit for the initial events
	fmt.Println("⏳ Waiting for initial events...")
	time.Sleep(2 * time.Second)

	// Test creating a notification
	fmt.Println("📧 Creating test notification...")
	testNotif := &provider.Notification{
		Name:          "Test Notification",
		Type:          "discord",
		IsDefault:     false,
		ApplyExisting: false,
		Config: map[string]interface{}{
			"discordWebhookUrl": "https://discord.com/api/webhooks/test",
		},
	}

	createdNotif, err := client.CreateNotification(testNotif)
	if err != nil {
		log.Printf("❌ Failed to create notification: %v", err)
	} else {
		fmt.Printf("✅ Created notification: ID=%d, Name=%s, Type=%s\n", createdNotif.ID, createdNotif.Name, createdNotif.Type)
	}

	// List all notifications again
	fmt.Println("📋 Getting all notifications after creation...")
	notifications, err := client.GetNotifications()
	if err != nil {
		log.Printf("❌ Failed to get notifications: %v", err)
	} else {
		fmt.Printf("✅ Found %d notifications:\n", len(notifications))
		for _, notif := range notifications {
			fmt.Printf("  - ID: %d, Name: %s, Type: %s, Active: %t\n", notif.ID, notif.Name, notif.Type, notif.Active)
		}
	}
}
