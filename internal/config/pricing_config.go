package config

// PricingConfig holds pricing configuration for different service types
type PricingConfig struct {
	BaseFare        float64
	PerKmRate       float64
	MinimumFare     float64
	SurgeMultiplier float64 // For future surge pricing
}

// GetPricingForService returns pricing config for a service type
func GetPricingForService(serviceType string) PricingConfig {
	// Default pricing (taxi)
	config := PricingConfig{
		BaseFare:        2.0, // $2 base fare
		PerKmRate:       4.0, // $4 per kilometer
		MinimumFare:     5.0, // $5 minimum
		SurgeMultiplier: 1.0, // No surge by default
	}

	// Customize pricing by service type
	switch serviceType {
	case "delivery":
		config.BaseFare = 1.5
		config.PerKmRate = 3.5
		config.MinimumFare = 4.0
	case "bus":
		config.BaseFare = 1.0
		config.PerKmRate = 2.0
		config.MinimumFare = 3.0
	case "taxi":
		// Use defaults above
	}

	return config
}
