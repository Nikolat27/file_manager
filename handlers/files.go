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
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var maxDownloads uint64
	if maxDownloadsStr != "" {
		maxDownloads, err = strconv.ParseUint(maxDownloadsStr, 0, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var expireAt time.Time
	if expireAtStr != "" {
		expireAt, err = time.Parse("2006-01-02T15:04", expireAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if newFileName == "" {
		newFileName = uuid.New().String()
	}

	address, err := utils.UploadFileToDisk("file", payload.UserId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := utils.ConvertStringToObjectID(payload.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortFileUrl := uuid.New().String()

	salt, err := utils.GenerateSalt()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, fmt.Errorf("ERROR creating file instance: %s", err).Error(), http.StatusBadRequest)
		return
	}

	if err := handler.Models.FileSettings.Create(fileId, encodedSalt, encodedPasswordHash, maxDownloads, false, approvable); err != nil {
		http.Error(w, fmt.Errorf("ERROR creating file share setting instance: %s", err).Error(), http.StatusBadRequest)
		return
	}

	resp := "file created successfully. Short Url: " + shortFileUrl
	w.Write([]byte(resp))
}

// GetFiles -> Returns List
func (handler *Handler) GetFiles(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR checking auth: %s", err).Error(), http.StatusUnauthorized)
		return
	}

	pageNumber, pageLimit, err := utils.GetPaginationParams(r)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR checking auth: %s", err).Error(), http.StatusUnauthorized)
		return
	}

	userId, err := utils.ConvertStringToObjectID(payload.UserId)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR checking auth: %s", err).Error(), http.StatusUnauthorized)
		return
	}

	file, err := handler.Models.File.GetAll(userId, pageNumber, pageLimit)

	data, err := json.MarshalIndent(file, "", "\t")
	w.Write(data)
}

// GetFile -> Returns One
func (handler *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	fileShortUrl, err := utils.ParseIdParam(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		RawPassword string `json:"password"`
	}

	// Optionally decode the JSON input for password (only POST requests)
	if r.Method == "POST" {
		if err := utils.ParseJsonData(r.Body, 1000, &input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	
	fileId, err := handler.Models.File.GetIdByShortUrl(fileShortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the fileShareSetting requires a password; alert if needed.
	requirePassword, err := handler.Models.FileSettings.IsPasswordRequired(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requirePassword && input.RawPassword == "" {
		http.Error(w, "Password required (send password via POST method)", http.StatusBadRequest)
		return
	}

	fileShareSetting, err := handler.Models.FileSettings.GetOne(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If password is required and provided, validate it
	if requirePassword && input.RawPassword != "" {
		if err := checkFilePassword([]byte(fileShareSetting.HashedPassword), []byte(fileShareSetting.Salt), []byte(input.RawPassword)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	file, err := handler.Models.File.GetOne(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If fileShareSetting`s accessing requires an approval, verify user`s approval status
	if fileShareSetting.Approvable {
		if err := checkUserApprovalStatus(r, handler, file.Id, file.OwnerId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	resp, err := json.MarshalIndent(&file, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (handler *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileObjectId, err := utils.ConvertStringToObjectID(fileId)

	address, err := handler.Models.File.GetDiskAddressById(fileObjectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.DeleteFileFromDisk(address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.Models.File.Delete(fileObjectId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("file deleted successfully"))
}

func (handler *Handler) RenameFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileObjectId, err := utils.ConvertStringToObjectID(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Name string `json:"new_name"`
	}

	if err := utils.ParseJsonData(r.Body, 1000, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "'new_name' parameter is missing", http.StatusBadRequest)
		return
	}

	if err := handler.Models.File.Rename(fileObjectId, []byte(input.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("file`s name changed successfully"))
}

func checkUserApprovalStatus(r *http.Request, handler *Handler, fileId, ownerId primitive.ObjectID) error {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		return errors.New("this file needs Approval, You must Logged in to send your approval")
	}

	userObjectId, err := utils.ConvertStringToObjectID(payload.UserId)
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
