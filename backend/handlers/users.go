package handlers

import (
	"encoding/hex"
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"net/http"
	"os"
)

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": userObjectId,
	}

	user, err := handler.Models.User.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := map[string]interface{}{
		"avatar_url": user.AvatarUrl,
	}

	utils.WriteJSONData(w, response)
}

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
	totalStorage, err := utils.GetUserTotalStorage(userPlan)
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
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": userObjectId,
	}

	projection := bson.M{
		"avatar_url": 1,
	}

	user, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	allowedTypes := []string{"image/jpg", "image/jpeg", "image/png", "image/webp"}

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

	filter := bson.M{
		"_id": userObjectId,
	}

	projection := bson.M{
		"total_upload_size": 1,
	}

	userInstance, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		return 0, err
	}

	return userInstance.TotalUploadSize, nil
}

func (handler *Handler) DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Password string `json:"password"`
	}

	filter := bson.M{
		"_id": userObjectId,
	}

	projection := bson.M{
		"hashed_password": 1,
		"salt":            1,
		"avatar_url":      1,
	}

	user, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if input.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "you must enter your password")
		return
	}

	decodedHashPassword, err := hex.DecodeString(user.HashedPassword)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	decodedSalt, err := hex.DecodeString(user.Salt)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !utils.ValidateHash([]byte(input.Password), decodedHashPassword, decodedSalt) {
		utils.WriteError(w, http.StatusBadRequest, "password is incorrect")
		return
	}

	if err := handler.Models.User.Delete(userObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.removeUserFilesAndAvatar(userObjectId, user.AvatarUrl); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "user deleted successfully")
}

func (handler *Handler) removeUserFilesAndAvatar(userId primitive.ObjectID, userAvatarUrl string) error {
	if err := os.Remove(userAvatarUrl); err != nil {
		slog.Error("removing user avatar", "error", err)
	}

	filter := bson.M{
		"owner_id": userId,
	}

	files, err := handler.Models.File.GetAll(filter, 1, 6)
	if err != nil {
		slog.Error("retrieving user file addresses", "error", err)
		return err
	}

	for _, file := range files {
		if err := os.Remove(file.Address); err != nil {
			slog.Error(fmt.Sprintf("removing user file: %s", file.Address), "error", err)
			continue
		}
	}

	return nil
}

func getUserAvatarUploadDir(userId string) string {
	return "uploads/user_files/" + userId + "/avatar/"
}
