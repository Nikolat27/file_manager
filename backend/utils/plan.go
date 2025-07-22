package utils

import (
	"errors"
	"fmt"
	"time"
)

const (
	GBytes = 1024 * 1024 * 1024

	UserFreePlanMaxStorageBytes    int64 = 2 * GBytes
	UserPlusPlanMaxStorageBytes    int64 = 100 * GBytes
	UserPremiumPlanMaxStorageBytes int64 = 1024 * GBytes

	TeamFreePlanMaxStorageBytes    int64 = 10 * GBytes
	TeamPremiumPlanMaxStorageBytes int64 = 1000 * GBytes
)

func GetUserMaxUploadSize(plan string) int64 {
	switch plan {
	case "free":
		return 100 << 20 // 100 MB
	case "plus":
		return 2000 << 20 // 2 GB
	case "premium":
		return 20000 << 10 // 10 GB
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

func GetUserTotalStorage(plan string) (int64, error) {
	switch plan {
	case "free":
		return UserFreePlanMaxStorageBytes, nil
	case "plus":
		return UserPlusPlanMaxStorageBytes, nil
	case "premium":
		return UserPremiumPlanMaxStorageBytes, nil
	case "":
		return 0, errors.New("plan is missing")
	default:
		return 0, fmt.Errorf("invalid plan: %s", plan)
	}
}

func GetTeamMaxUploadSize(plan string) int64 {
	switch plan {
	case "free":
		return 100 << 20 // 100 MB
	case "premium":
		return 2000 << 20 // 2 GB
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

func GetTeamTotalStorage(plan string) (int64, error) {
	switch plan {
	case "free":
		return TeamFreePlanMaxStorageBytes, nil
	case "premium":
		return TeamPremiumPlanMaxStorageBytes, nil
	case "":
		return 0, errors.New("plan is missing")
	default:
		return 0, fmt.Errorf("invalid plan: %s", plan)
	}
}

func GetTeamTotalUsers(plan string) int {
	// Unlimited users for premium plan
	if plan == "premium" {
		return -1
	}

	return 5
}

func GetTeamUploadDir(teamId string) string {
	return "uploads/team_files/files/" + teamId + "/"
}

// ValidateUserPlan -> free, plus, premium
func ValidateUserPlan(plan string) error {
	if plan == "free" || plan == "plus" || plan == "premium" {
		return nil
	} else {
		return fmt.Errorf("plan is invalid: %s. Must be either free, plus or premium", plan)
	}
}

// ValidateTeamPlan -> free, premium
func ValidateTeamPlan(plan string) error {
	if plan == "free" || plan == "premium" {
		return nil
	} else {
		return fmt.Errorf("plan is invalid: %s. Must be either free or premium", plan)
	}
}
