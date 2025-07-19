package utils

import (
	"errors"
	"fmt"
)

const (
	GBytes      = 1024 * 1024 * 1024
	TeamFreePlanMaxStorageBytes    int64 = 2 * GBytes
	TeamPremiumPlanMaxStorageBytes int64 = 1024 * GBytes
)

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

func GetTeamUploadDir(teamId string) string {
	return "uploads/team_files/files/" + teamId + "/"
}
