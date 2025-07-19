package handlers

import (
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
)

const (
	GBytes                               = 1024 * 1024 * 1024
	UserFreePlanMaxStorageBytes    int64 = 2 * GBytes
	UserPlusPlanMaxStorageBytes    int64 = 100 * GBytes
	UserPremiumPlanMaxStorageBytes int64 = 1024 * GBytes
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

	if err := utils.ValidateUserPlan(input.Plan); err != nil {
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

func (handler *Handler) UploadUserAvatar(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var maxAvatarSize int64 = 5 << 20 // 5 MB
	if err := r.ParseMultipartForm(maxAvatarSize); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := handler.Models.User.GetById(userObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	allowedTypes := []string{"image/jpg", "image/jpeg", "image/png"}

	file, err := utils.ReadFile(r, maxAvatarSize, allowedTypes)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	uploadDir := getUserAvatarUploadDir(payload.UserId)

	fileAddress, err := file.UploadToDisk(uploadDir)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	updates := bson.M{
		"avatar_url": fileAddress,
	}

	if user.AvatarUrl != "" {
		if err := os.Remove(user.AvatarUrl); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	if err := handler.Models.User.Update(userObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "user`s avatar uploaded successfully")
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

func getUserAvatarUploadDir(userId string) string {
	return "uploads/user_files/" + userId + "/avatar/"
}
