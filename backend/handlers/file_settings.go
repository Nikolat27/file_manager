package handlers

import (
	"encoding/hex"
	"errors"
	"file_manager/utils"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

	filter := bson.M{
		"_id": fileObjectId,
	}

	projection := bson.M{
		"_id":      1,
		"owner_id": 1,
	}

	file, err := handler.Models.File.Get(filter, projection)
	if file.OwnerId != userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "only the file owner can create settings(shortUrl) for it")
		return
	}

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"file_id": fileObjectId,
	}

	projection = bson.M{
		"file_id": 1,
	}

	_, err = handler.Models.FileSettings.Get(filter, projection)
	// If it exists, return error
	if !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "setting with this filter does not exist")
		return
	}

	hashedPassword, salt, err := getPasswordAndSalt(r)
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

	if err := handler.Models.FileSettings.Create(fileObjectId, userObjectId, fileShortUrl.String(), salt, hashedPassword, maxDownloads,
		viewOnly, approvable, expireAt); err != nil {

		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file share setting instance: %w", err))
		return
	}

	data := map[string]string{
		"short_url": fileShortUrl.String(),
	}

	utils.WriteJSONData(w, data)
}

func (handler *Handler) GetFilesSettings(w http.ResponseWriter, r *http.Request) {
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
		"user_id": userObjectId,
	}

	files, err := handler.Models.FileSettings.GetAll(filter)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := map[string]any{
		"sharedUrls": files,
	}

	utils.WriteJSONData(w, response)
}

func (handler *Handler) DeleteFileSettings(w http.ResponseWriter, r *http.Request) {
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

	settingId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	settingObjectId, err := utils.ToObjectID(settingId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id":     settingObjectId,
		"user_id": userObjectId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.FileSettings.Get(filter, projection); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"_id": settingObjectId,
	}

	if err := handler.Models.FileSettings.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "setting deleted successfully")
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
	if maxDownloadsStr == "" || maxDownloadsStr == "-1" {
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
		return time.Now().Add(7 * time.Hour * 24), nil
	}

	if plan == "" {
		return time.Time{}, errors.New("user`s plan is missing")
	}

	if plan == "free" {
		return time.Now().Add(7 * time.Hour * 24), nil
	}

	return time.Parse(time.RFC3339, expireAtStr)
}
