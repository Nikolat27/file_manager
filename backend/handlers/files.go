package handlers

import (
	"encoding/json"
	"errors"
	"file_manager/database/models"
	"file_manager/utils"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"path"
)

func (handler *Handler) UploadUserFile(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	maxUploadSize := utils.GetUserMaxUploadSize(payload.UserPlan)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	fileName := r.FormValue("file_name")
	if fileName == "" {
		fileName = uuid.New().String()
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	folderObjectId, err := getFolderId(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.ValidateFolderId(folderObjectId, userObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	uploadDir := getUserUploadDir(payload.UserId)

	fileAddress, totalUserUploadSize, err := handler.storeUserFile(r, maxUploadSize, payload.UserId, payload.UserPlan, uploadDir)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	expireAt := utils.GetUserExpirationDate(payload.UserPlan)

	// no teamId for user uploaded files
	teamId := primitive.NilObjectID
	if _, err := handler.Models.File.Create(userObjectId, teamId, folderObjectId, fileName, fileAddress, expireAt); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file instance: %w", err))
		return
	}

	updates := bson.M{"total_upload_size": totalUserUploadSize}
	if err := handler.Models.User.Update(userObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("updating user instance: %w", err))
		return
	}

	utils.WriteJSON(w, "file uploaded successfully")
}

// GetFiles -> Returns List
func (handler *Handler) GetFiles(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("checking auth: %w", err))
		return
	}

	pageNumber, pageLimit, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("pagination params: %w", err))
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user id: %w", err))
		return
	}

	filter := bson.M{
		"owner_id": userObjectId,
	}

	file, err := handler.Models.File.GetAll(filter, pageNumber, pageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	data, err := json.MarshalIndent(file, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, data)
}

func (handler *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileObjectId, err := utils.ToObjectID(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": fileObjectId,
	}

	projection := bson.M{
		"address": 1,
	}

	fileInstance, err := handler.Models.File.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.DeleteFileFromDisk(fileInstance.Address); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.Models.File.Delete(fileObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.Models.FileSettings.Delete(fileObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "file deleted successfully")
}

func (handler *Handler) RenameFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileObjectId, err := utils.ToObjectID(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("'name' parameter is missing"))
		return
	}

	updates := bson.M{
		"name": input.Name,
	}

	if err := handler.Models.File.Update(fileObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "file`s name changed successfully")
}

func (handler *Handler) SearchFiles(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		SearchText string `json:"search_text"`
		Page       int64  `json:"page"`
		PageLimit  int64  `json:"page_limit"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageLimit <= 0 || input.PageLimit > 20 {
		input.PageLimit = 6
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"owner_id": userObjectId,
		"name": bson.M{
			"$regex": input.SearchText, "$options": "i", // case insensitive
		},
	}

	// Search through files names
	files, err := handler.Models.File.GetAll(filter, input.Page, input.PageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	data, err := json.MarshalIndent(files, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, data)
}

func (handler *Handler) storeUserFile(r *http.Request, maxUploadSize int64, userId, userPlan, uploadDir string) (string, int64, error) {
	allowedTypes := []string{"image/jpeg", "image/png", "application/zip", "application/pdf"}

	file, err := utils.ReadFile(r, maxUploadSize, allowedTypes)
	if err != nil {
		return "", 0, err
	}

	defer file.File.Close()

	totalUsedStorage, err := handler.IsUserEligibleToUpload(userId, userPlan, file.Size)
	if err != nil {
		return "", 0, err
	}

	fileAddress, err := file.UploadToDisk(uploadDir)
	if err != nil {
		return "", 0, err
	}

	return fileAddress, totalUsedStorage, nil
}

func (handler *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileObjectId, err := utils.ToObjectID(fileIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": fileObjectId,
	}

	projection := bson.M{
		"address": 1,
	}

	fileInstance, err := handler.Models.File.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, err := os.Open(fileInstance.Address)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(file.Name())))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, file); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "download started successfully")
}

func getUserUploadDir(userId string) string {
	return "uploads/user_files/" + userId + "/files/"
}

// GetFile -> Returns One
func (handler *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var providedPassword string
	if r.Method == "POST" {
		var input struct {
			Password string `json:"password"`
		}

		if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		providedPassword = input.Password
	}

	filter := bson.M{
		"short_url": shortUrl,
	}

	projection := bson.M{
		"file_id":         1,
		"approvable":      1,
		"salt":            1,
		"hashed_password": 1,
	}

	fileShareSettings, err := handler.Models.FileSettings.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"_id": fileShareSettings.FileId,
	}

	projection = bson.M{
		"owner_id": 1,
		"address":  1,
	}

	file, err := handler.Models.File.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var requesterId primitive.ObjectID
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err == nil {
		requesterId, _ = utils.ToObjectID(payload.UserId)
	}

	if err := checkPasswordAccess(file.OwnerId, requesterId, providedPassword, fileShareSettings); err != nil {
		utils.WriteError(w, http.StatusNotAcceptable, err)
		return
	}

	if fileShareSettings.Approvable {
		utils.WriteError(w, http.StatusPreconditionRequired, "approval is required")
		return
	}

	utils.WriteJSONData(w, map[string]any{"file_address": file.Address})
	return
}

func checkPasswordAccess(ownerId, requesterId primitive.ObjectID, rawPassword string, fileSettings *models.FileSettings) error {
	if ownerId == requesterId || requesterId == primitive.NilObjectID {
		return nil
	}

	if fileSettings.HashedPassword == "" {
		return nil
	}

	if rawPassword == "" {
		return errors.New("password is required")
	}

	return utils.CheckFilePassword(
		[]byte(fileSettings.HashedPassword),
		[]byte(fileSettings.Salt),
		[]byte(rawPassword),
	)
}
