package handlers

import (
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

const (
	GB                               = 1024 * 1024 * 1024
	FreePlanMaxStorageBytes    int64 = 2 * GB
	PlusPlanMaxStorageBytes    int64 = 100 * GB
	PremiumPlanMaxStorageBytes int64 = 1024 * GB
)

func (handler *Handler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		Plan string `json:"plan"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidatePlan(input.Plan); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	updates := bson.M{
		"plan": input.Plan,
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.Models.User.Update(userObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "user`s plan updated successfully")
}

func (handler *Handler) IsUserEligibleToUpload(userId, userPlan string, fileSize int64) (int64, error) {
	totalStorage, err := getUserTotalStorage(userPlan)
	if err != nil {
		return 0, err
	}

	if fileSize > totalStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your plan's total storage limit (%d bytes)", fileSize, totalStorage)
	}

	usedStorage, err := handler.getUsedStorage(userId)
	if err != nil {
		return 0, fmt.Errorf("failed to check used storage: %w", err)
	}

	remainedStorage := totalStorage - usedStorage
	if fileSize > remainedStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your remaining storage (%d bytes)", fileSize, remainedStorage)
	}

	newTotalStorage := fileSize + usedStorage
	return newTotalStorage, nil
}

func (handler *Handler) getUsedStorage(userId string) (int64, error) {
	if userId == "" {
		return 0, errors.New("user id is missing")
	}

	userObjectId, err := utils.ToObjectID(userId)
	if err != nil {
		return 0, err
	}

	return handler.Models.User.GetUsedStorage(userObjectId)
}

func getUserTotalStorage(plan string) (int64, error) {
	switch plan {
	case "free":
		return FreePlanMaxStorageBytes, nil
	case "plus":
		return PlusPlanMaxStorageBytes, nil
	case "premium":
		return PremiumPlanMaxStorageBytes, nil
	case "":
		return 0, errors.New("plan is missing")
	default:
		return 0, fmt.Errorf("invalid plan: %s", plan)
	}
}
