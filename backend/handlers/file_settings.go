package handlers

import (
	"encoding/hex"
	"errors"
	"file_manager/utils"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func (handler *Handler) CreateFileSettings(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	fileId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if fileId == "" {
		utils.WriteError(w, http.StatusBadRequest, "file id is missing")
		return
	}

	fileObjectId, err := utils.ToObjectID(fileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileExist, err := handler.Models.File.IsExist(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !fileExist {
		utils.WriteError(w, http.StatusBadRequest, "file with this id does not exist")
		return
	}

	settingExist, err := handler.Models.FileSettings.IsExist(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if settingExist {
		utils.WriteError(w, http.StatusBadRequest, "this file has settings already")
		return
	}

	password, salt, err := getPasswordAndSalt(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	approvable, err := getApproval(r, payload.UserPlan)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	viewOnly, err := getViewOnly(r, payload.UserPlan)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	maxDownloads, err := getMaxDownloads(r, payload.UserPlan)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	expireAt, err := getExpireAt(r, payload.UserPlan)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileShortUrl, err := uuid.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed generating file short url: %s", err))
		return
	}

	if err = handler.Models.FileSettings.Create(fileObjectId, fileShortUrl.String(), salt, password, maxDownloads,
		viewOnly, approvable, expireAt); err != nil {

		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file share setting instance: %w", err))
		return
	}

	data := map[string]string{
		"short_url": fileShortUrl.String(),
	}

	utils.WriteJSONData(w, data)
}

func getPasswordAndSalt(r *http.Request) (string, string, error) {
	rawPassword := r.FormValue("password")
	if rawPassword == "" {
		return "", "", nil
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		return "", "", err
	}

	hashedPassword := utils.Hash256([]byte(rawPassword), salt)

	encodedPasswordHash := hex.EncodeToString(hashedPassword[:])
	encodedSalt := hex.EncodeToString(salt)

	return encodedPasswordHash, encodedSalt, nil
}

func getApproval(r *http.Request, plan string) (bool, error) {
	approvalStr := r.FormValue("approvable")
	if approvalStr == "" || approvalStr == "false" {
		return false, nil
	}

	if plan == "" {
		return false, errors.New("user`s plan is missing")
	}

	if plan == "free" {
		return false, errors.New("users with 'free' plan cant make approval required URLs")
	}

	return strconv.ParseBool(approvalStr)
}

func getMaxDownloads(r *http.Request, plan string) (int64, error) {
	maxDownloadsStr := r.FormValue("max_downloads")
	if maxDownloadsStr == "" {
		return -1, nil // -1 means unlimited downloads
	}

	if plan == "" {
		return -1, errors.New("user`s plan is missing")
	}

	if plan == "free" {
		return -1, errors.New("users with 'free' plan can`t use 'max_downloads' feature")
	}

	return strconv.ParseInt(maxDownloadsStr, 0, 64)
}

func getViewOnly(r *http.Request, plan string) (bool, error) {
	viewOnlyStr := r.FormValue("view_only")
	if viewOnlyStr == "" || viewOnlyStr == "false" {
		return false, nil
	}

	if plan == "" {
		return false, errors.New("user`s plan is missing")
	}

	if plan == "free" {
		return false, errors.New("users with 'free' plan can`t change the view mode (by default its read-write)")
	}

	return strconv.ParseBool(viewOnlyStr)
}

func getExpireAt(r *http.Request, plan string) (time.Time, error) {
	expireAtStr := r.FormValue("expiration_at")
	if expireAtStr == "" {
		return time.Time{}, nil
	}

	if plan == "" {
		return time.Time{}, errors.New("user`s plan is missing")
	}

	if plan == "free" {
		return time.Now().Add(7 * time.Hour * 24), nil
	}

	return time.Parse(expireAtStr, time.DateTime)
}
