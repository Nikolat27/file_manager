package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"file_manager/utils"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

const (
	MaxUploadSize = 100 << 20
)

func (handler *Handler) CreateFile(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	newFileName := r.FormValue("file_name")
	approvableStr := r.FormValue("approvable")
	rawPassword := r.FormValue("password")
	maxDownloadsStr := r.FormValue("max_downloads")
	expireAtStr := r.FormValue("expire_at")

	var approvable bool
	if approvableStr != "" {
		approvable, err = strconv.ParseBool(approvableStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	var maxDownloads uint64
	if maxDownloadsStr != "" {
		maxDownloads, err = strconv.ParseUint(maxDownloadsStr, 0, 64)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	var expireAt time.Time
	if expireAtStr != "" {
		expireAt, err = time.Parse("2006-01-02T15:04", expireAtStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	if newFileName == "" {
		newFileName = uuid.New().String()
	}

	address, err := utils.UploadFileToDisk("file", payload.UserId, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	shortFileUrl := uuid.New().String()

	salt, err := utils.GenerateSalt()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var hashedPassword [32]byte
	var encodedPasswordHash string

	if rawPassword != "" {
		hashedPassword = utils.Hash256([]byte(rawPassword), salt)
		encodedPasswordHash = hex.EncodeToString(hashedPassword[:])
	}

	encodedSalt := hex.EncodeToString(salt)

	fileId, err := handler.Models.File.Create(userId, newFileName, address, shortFileUrl, expireAt)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file instance: %w", err))
		return
	}

	if err := handler.Models.FileSettings.Create(fileId, encodedSalt, encodedPasswordHash, maxDownloads, false, approvable); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file share setting instance: %w", err))
		return
	}

	resp := "file created successfully. Short Url: " + shortFileUrl

	utils.WriteJSON(w, resp)
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

	fileId, err := handler.Models.File.GetIdByShortUrl(fileShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if the fileShareSetting requires a password; alert if needed.
	requirePassword, err := handler.Models.FileSettings.IsPasswordRequired(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if requirePassword && input.RawPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("password required (send password via POST method)"))
		return
	}

	fileShareSetting, err := handler.Models.FileSettings.GetOne(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// If password is required and provided, validate it
	if requirePassword && input.RawPassword != "" {
		if err := checkFilePassword([]byte(fileShareSetting.HashedPassword), []byte(fileShareSetting.Salt), []byte(input.RawPassword)); err != nil {
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
		Name string `json:"new_name"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if input.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("'new_name' parameter is missing"))
		return
	}

	if err := handler.Models.File.Rename(fileObjectId, []byte(input.Name)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "file`s name changed successfully")
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

func checkFilePassword(hashedPassword, salt, rawPassword []byte) error {
	decodedHashPassword, err := hex.DecodeString(string(hashedPassword))
	if err != nil {
		return err
	}

	decodedSalt, err := hex.DecodeString(string(salt))
	if err != nil {
		return err
	}

	if !utils.ValidateHash(rawPassword, decodedHashPassword, decodedSalt) {
		return errors.New("password is invalid")
	}

	return nil
}
