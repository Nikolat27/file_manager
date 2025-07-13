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
	rawPassword := r.FormValue("raw_password")
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

	address, err := utils.UploadFile("file", payload.UserId, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := utils.ConvertStringToObjectID(payload.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hashedPassword [32]byte
	if rawPassword != "" {
		hashedPassword = utils.Hash256([]byte(rawPassword), salt)
	}

	encodedSalt := hex.EncodeToString(salt)
	encodedPasswordHash := hex.EncodeToString(hashedPassword[:])

	_, err = handler.Models.File.CreateFileInstance(userId, newFileName, address, encodedSalt, encodedPasswordHash, approvable, maxDownloads, expireAt)
	if err != nil {
		http.Error(w, fmt.Errorf("ERROR creating file instance: %s", err).Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("file created successfully"))
}

func (handler *Handler) GetUserFiles(w http.ResponseWriter, r *http.Request) {
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
	
	file, err := handler.Models.File.GetUsersFileInstance(userId, pageNumber, pageLimit)

	data, err := json.MarshalIndent(file, "", "	")
	w.Write(data)
}
