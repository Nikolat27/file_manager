package handlers

import (
	"encoding/hex"
	"encoding/json"
	"file_manager/utils"
	"fmt"
	"github.com/google/uuid"
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

	newFileName := r.FormValue("file")
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

	_, err = handler.Models.File.CreateFileInstance(userId, newFileName, address, shortFileUrl, encodedSalt, encodedPasswordHash, approvable, maxDownloads, expireAt)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR creating file instance: %s", err).Error(), http.StatusBadRequest)
		return
	}

	resp := "file created successfully. Short Url: " + shortFileUrl
	w.Write([]byte(resp))
}

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

	file, err := handler.Models.File.GetFilesInstances(userId, pageNumber, pageLimit)

	data, err := json.MarshalIndent(file, "", "	")
	w.Write(data)
}

func (handler *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ReadFileId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileObjectId, err := utils.ConvertStringToObjectID(fileId)

	address, err := handler.Models.File.GetFileAddress(fileObjectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.DeleteFileFromDisk(address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.Models.File.DeleteFileInstance(fileObjectId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("file deleted successfully"))
}

func (handler *Handler) RenameFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := utils.ReadFileId(r)
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

	if err := utils.ReadJson(r, 1000, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "'new_name' parameter is missing", http.StatusBadRequest)
		return
	}

	if err := handler.Models.File.RenameFileInstance(fileObjectId, []byte(input.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("file`s name changed successfully"))
}

func (handler *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	url, err := utils.ReadFileShortUrl(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		RawPassword string `json:"password"`
	}

	if r.Method == "POST" {
		if err := utils.ReadJson(r, 1000, &input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	requirePassword, err := handler.Models.File.RequirePassword(url)
	if requirePassword && input.RawPassword == "" {
		http.Error(w, "password is required (The Password must be sent via POST method)", http.StatusBadRequest)
		return
	}

	file, err := handler.Models.File.GetFileInstance(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.RawPassword != "" && requirePassword {
		decodedHashPassword, err := hex.DecodeString(file.HashedPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		decodedSalt, err := hex.DecodeString(file.Salt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !utils.ValidateHash([]byte(input.RawPassword), decodedHashPassword, decodedSalt) {
			http.Error(w, "password is incorrect", http.StatusBadRequest)
			return
		}
	}

	staticFileUrl, err := utils.GetStaticFilesUrl(file.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file.Address = staticFileUrl
	
	resp, err := json.MarshalIndent(file, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
