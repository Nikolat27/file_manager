package handlers

import (
	"crypto/rand"
	"encoding/json"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (handler *Handler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		TeamId string `json:"team_id"`
		Name   string `json:"name"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if input.Name == "" {
		input.Name = rand.Text()
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var teamObjectId primitive.ObjectID
	if input.TeamId != "" {
		teamObjectId, err = utils.ToObjectID(input.TeamId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	if _, err := handler.Models.Folder.Create(userObjectId, teamObjectId, input.Name); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "folder created successfully")
}

func (handler *Handler) GetFolderContents(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	folderId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	folderObjectId, err := utils.ToObjectID(folderId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id":      folderObjectId,
		"owner_id": userObjectId,
	}

	projection := bson.M{
		"_id":  1,
		"name": 1,
	}

	folder, err := handler.Models.Folder.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	page, pageSize, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"folder_id": folderObjectId,
	}

	// getting the files
	files, err := handler.Models.File.GetAll(filter, page, pageSize)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := map[string]any{
		"folder_id":   folderId,
		"folder_name": folder.Name,
		"files":       files,
	}

	utils.WriteJSONData(w, response)
}

func (handler *Handler) ValidateFolderId(folderId, userId primitive.ObjectID) error {
	if folderId == primitive.NilObjectID {
		return nil
	}

	filter := bson.M{
		"_id":      folderId,
		"owner_id": userId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Folder.Get(filter, projection); err != nil {
		return err
	}

	return nil
}

func (handler *Handler) RenameFolder(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	folderId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	folderObjectId, err := utils.ToObjectID(folderId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id":      folderObjectId,
		"owner_id": userObjectId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Folder.Get(filter, projection); err != nil {
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

	updates := bson.M{
		"name":      input.Name,
		"expire_at": time.Now(),
	}

	if err := handler.Models.Folder.Rename(folderObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "folder`s name changed successfully")
}

func (handler *Handler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	folderId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	folderObjectId, err := utils.ToObjectID(folderId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": folderObjectId,
	}

	projection := bson.M{
		"_id":      1,
		"owner_id": 1,
	}

	folderInstance, err := handler.Models.Folder.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if folderInstance.OwnerId != userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "only the folder owner can delete it")
		return
	}

	if err := handler.Models.Folder.Delete(folderObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "folder deleted successfully")
}

func (handler *Handler) GetFoldersList(w http.ResponseWriter, r *http.Request) {
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

	teamId := r.URL.Query().Get("team_id")

	var filter bson.M
	if teamId != "" {
		teamObjectId, err := utils.ToObjectID(teamId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user id: %w", err))
			return
		}

		filter = bson.M{
			"team_id": teamObjectId,
		}
	} else {
		filter = bson.M{
			"owner_id": userObjectId,
			"team_id":  primitive.NilObjectID,
		}
	}

	page, pageSize, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	files, err := handler.Models.Folder.GetAll(filter, page, pageSize)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	data, err := json.MarshalIndent(&files, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, data)
}

func getFolderId(r *http.Request) (primitive.ObjectID, error) {
	folderId := r.FormValue("folder_id")
	if folderId == "" {
		return primitive.NilObjectID, nil
	}

	folderObjectId, err := utils.ToObjectID(folderId)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return folderObjectId, nil
}
