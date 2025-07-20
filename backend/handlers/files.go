package handlers

import (
	"encoding/json"
	"errors"
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

	userId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user id: %w", err))
		return
	}

	file, err := handler.Models.File.GetAll(userId, pageNumber, pageLimit)
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

// GetFile -> Returns One
func (handler *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	fileShortUrl, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var input struct {
		RawPassword string `json:"password"`
	}

	// Optionally decode the JSON input for password (only POST requests)
	if r.Method == "POST" {
		if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	fileId, err := handler.Models.FileSettings.GetFileIdByUrl(fileShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if the fileShareSetting requires a password. alert if needed.
	requirePassword, err := handler.Models.FileSettings.IsPasswordRequired(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if requirePassword && input.RawPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("password required (send password via POST method)"))
		return
	}

	fileShareSetting, err := handler.Models.FileSettings.Get(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// If password is required and provided, validate it
	if requirePassword && input.RawPassword != "" {
		if err := utils.CheckFilePassword([]byte(fileShareSetting.HashedPassword), []byte(fileShareSetting.Salt), []byte(input.RawPassword)); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	file, err := handler.Models.File.GetOne(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// If fileShareSetting`s accessing requires an approval, verify user`s approval status
	if fileShareSetting.Approvable {
		if err := checkUserApprovalStatus(r, handler, file.Id, file.OwnerId); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	resp, err := json.MarshalIndent(&file, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, resp)
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

	address, err := handler.Models.File.GetDiskAddressById(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.DeleteFileFromDisk(address); err != nil {
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

	if err := handler.Models.File.Rename(fileObjectId, []byte(input.Name)); err != nil {
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

	// Search through files names
	files, err := handler.Models.File.Search(userObjectId, input.SearchText, input.Page, input.PageLimit)
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
	allowedTypes := []string{"image/jpeg", "image/png", "application/zip"}

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

func getUserUploadDir(userId string) string {
	return "uploads/user_files/" + userId + "/files/"
}

func checkUserApprovalStatus(r *http.Request, handler *Handler, fileId, ownerId primitive.ObjectID) error {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		return errors.New("this file needs approval, you must be logged in to send your approval")
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		return err
	}

	// The owner does not need approval
	if userObjectId == ownerId {
		return nil
	}

	status, err := handler.Models.Approval.CheckStatus(fileId, userObjectId)
	if err != nil {
		return err
	}

	if status == "rejected" {
		return errors.New("your approval has been rejected by the file owner")
	}

	if status == "pending" {
		return errors.New("your approval is in pending status. Please be patient")
	}

	return nil
}

func (handler *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileShortUrl, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileObjectId, err := utils.ToObjectID(fileShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileAddress, err := handler.Models.File.GetDiskAddressById(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, err := os.Open(string(fileAddress))
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
