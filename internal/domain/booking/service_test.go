package booking

import (
	"testing"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
)

func TestIsNoiseCompatible(t *testing.T) {
	tests := []struct {
		name      string
		userNoise user.NoiseLevel
		locLimit  location.NoiseLimit
		expected  bool
	}{
		// Light Limit
		{"Light user on Light limit", user.NoiseLight, location.LimitLight, true},
		{"Medium user on Light limit", user.NoiseMedium, location.LimitLight, false},
		{"Hard user on Light limit", user.NoiseHard, location.LimitLight, false},
		{"None user on Light limit", user.NoiseNone, location.LimitLight, false},

		// Medium Limit
		{"Light user on Medium limit", user.NoiseLight, location.LimitMedium, true},
		{"Medium user on Medium limit", user.NoiseMedium, location.LimitMedium, true},
		{"Hard user on Medium limit", user.NoiseHard, location.LimitMedium, false},
		{"None user on Medium limit", user.NoiseNone, location.LimitMedium, false},

		// Hard Limit
		{"Light user on Hard limit", user.NoiseLight, location.LimitHard, false},
		{"Medium user on Hard limit", user.NoiseMedium, location.LimitHard, true},
		{"Hard user on Hard limit", user.NoiseHard, location.LimitHard, true},
		{"None user on Hard limit", user.NoiseNone, location.LimitHard, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNoiseCompatible(tt.userNoise, tt.locLimit)
			if got != tt.expected {
				t.Errorf("IsNoiseCompatible(%s, %s) = %v; want %v", tt.userNoise, tt.locLimit, got, tt.expected)
			}
		})
	}
}
