package jobs

import (
	"time"

	"github.com/telemoz/backend/internal/services"
)

// StartTripExpirationJob runs a background job to expire trips that have been searching for >3 minutes
func StartTripExpirationJob(tripService services.TripService) {
	ticker := time.NewTicker(30 * time.Second) // Run every 30 seconds

	go func() {
		for range ticker.C {
			err := tripService.ExpireSearchingTrips()
			if err != nil {
				// Log error but continue
				println("Error expiring trips:", err.Error())
			} else {
				println("✅ Trip expiration job completed")
			}
		}
	}()

	println("🕐 Trip expiration background job started (runs every 30s)")
}
