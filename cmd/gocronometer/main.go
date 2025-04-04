package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/cttautomations/gocronometer"
)

func main() {
	username := os.Getenv("CRONOMETER_USERNAME")
	password := os.Getenv("CRONOMETER_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Missing CRONOMETER_USERNAME or CRONOMETER_PASSWORD environment variables")
	}

	location, _ := time.LoadLocation("America/New_York")

	// Get current local time
	now := time.Now().In(location)

	// Force today's date (00:00 EST)
	todayDate := now.Format("2006-01-02")
	today, _ := time.ParseInLocation("2006-01-02", todayDate, location)

	client := gocronometer.NewClient(nil)
	if err := client.Login(context.Background(), username, password); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	ctx := context.Background()

	// Parsed exports
	servings, _ := client.ExportServingsParsedWithLocation(ctx, today, today, location)
	biometrics, _ := client.ExportBiometricRecordsParsedWithLocation(ctx, today, today, location)
	exercises, _ := client.ExportExercisesParsedWithLocation(ctx, today, today, location)
	notes, _ := client.ExportNotes(ctx, today, today)

	// Raw daily summary CSV
	nutrition, _ := client.ExportDailyNutrition(ctx, today, today)

	payload := map[string]interface{}{
		"date":       todayDate,
		"timestamp":  now.Format(time.RFC3339),
		"servings":   servings,
		"biometrics": biometrics,
		"exercises":  exercises,
		"notes":      notes,
		"nutrition":  nutrition,
	}

	if err := json.NewEncoder(os.Stdout).Encode(payload); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}
}
