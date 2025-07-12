package handlers

import (
	"file_manager/utils"
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

	handler.Models.File.CreateFileInstance(userId, newFileName, address, hashedPassword[:], salt, approvable, maxDownloads, expireAt)

}
