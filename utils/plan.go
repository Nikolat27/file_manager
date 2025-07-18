package utils

import (
	"fmt"
	"time"
)

// Available Plans: free, plus, premium

func GetUserMaxUploadSize(plan string) int64 {
	switch plan {
	case "free":
		return 100 << 20 // 100 MB
	case "plus":
		return 2000 << 20 // 2 GB
	case "premium":
		return 20000 << 20 // 20 GB
	default:
		return 100 << 20 // default 100 MB
	}
}

func GetUserExpirationDate(plan string) time.Time {
	switch plan {
	case "free":
		return time.Now().Add(7 * time.Hour * 24) // 7 Days
	case "plus":
		return time.Now().Add(30 * time.Hour * 24) // 30 Days
	case "premium":
		return time.Now().Add(180 * time.Hour * 24) // 180 Days
	default:
		return time.Now().Add(7 * time.Hour * 24) // 7 Days
	}
}

func GetTeamMaxUploadSize(plan string) int64 {
	switch plan {
	case "free":
		return 100 << 20 // 100 MB
	case "premium":
		return 10 << 20 // 10 GB
	default:
		return 100 << 20 // default 100 MB
	}
}

func GetTeamExpirationDate(plan string) time.Time {
	switch plan {
	case "free":
		return time.Now().Add(14 * time.Hour * 24) // 14 Days
	case "premium":
		return time.Now().Add(120 * time.Hour * 24) // 120 Days
	default:
		return time.Now().Add(7 * time.Hour * 24) // 7 Days
	}
}

func ValidatePlan(plan string) error {
	if plan == "free" || plan == "plus" || plan == "premium" {
		return nil
	} else {
		return fmt.Errorf("plan is invalid: %s", plan)
	}
}
