package handlers

import (
	"crypto/rand"
	"encoding/json"
	"file_manager/utils"
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
		Name string `json:"name"`
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

	if _, err := handler.Models.Folder.Create(userObjectId, input.Name); err != nil {
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

	if err := handler.Models.Folder.Validate(folderObjectId, userObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	page, pageSize, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// getting the files
	files, err := handler.Models.File.GetByFolderId(folderObjectId, page, pageSize)
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

func (handler *Handler) ValidateFolderId(folderId, userId primitive.ObjectID) error {
	return handler.Models.Folder.Validate(folderId, userId)
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

	if err := handler.Models.Folder.Validate(folderObjectId, userObjectId); err != nil {
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

	if err := handler.Models.Folder.Validate(folderObjectId, userObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.Models.Folder.Delete(folderObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "folder deleted successfully")
}

func (handler *Handler) GetListFolders(w http.ResponseWriter, r *http.Request) {
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

	page, pageSize, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	files, err := handler.Models.Folder.GetAll(userObjectId, page, pageSize)
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
