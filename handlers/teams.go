package handlers

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"slices"
	"time"
)

const (
	MegaBytes int64 = 1 << 20
)

func (handler *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if err := utils.ValidateUserPlan(payload.UserPlan); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.UserPlan != "premium" {
		utils.WriteError(w, http.StatusBadRequest, "Team creation is only available for 'premium' plan users")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*MegaBytes)
	name := r.FormValue("name")
	if name == "" {
		name = rand.Text()
	}

	description := r.FormValue("description")

	teamId := primitive.NewObjectID()
	avatarAddress, err := handler.uploadAvatar(r, 5<<20, teamId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := handler.Models.Team.Create(teamId, userObjectId, name, description, avatarAddress); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, fmt.Sprintf("team created successfully. Id: %s", teamId.Hex()))
}

func (handler *Handler) uploadAvatar(r *http.Request, maxUploadSize int64, teamId string) (string, error) {
	allowedTypes := []string{"image/jpeg", "image/png", "image/jpg"}

	file, err := utils.ReadFile(r, maxUploadSize, allowedTypes)
	if err != nil {
		return "", err
	}

	uploadDir := "uploads/team_files/avatars/" + teamId + "/"
	return file.UploadToDisk(uploadDir)
}

func (handler *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)

	teamIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamObjectId, err := utils.ToObjectID(teamIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamInstance, err := handler.Models.Team.Get(teamObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !slices.Contains(teamInstance.Users, userObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "you are not member of this team")
		return
	}

	data, err := json.MarshalIndent(teamInstance, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, data)
}

func (handler *Handler) UpdateTeamPlan(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	ownerObjectId, err := utils.ToObjectID(payload.UserId)

	teamIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamObjectId, err := utils.ToObjectID(teamIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamInstance, err := handler.Models.Team.Get(teamObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if teamInstance.OwnerId != ownerObjectId {
		utils.WriteError(w, http.StatusBadRequest, "only the team owner can update its plan")
		return
	}

	var input struct {
		Plan string `json:"plan"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.ValidateTeamPlan(input.Plan); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	updates := bson.M{
		"plan":       input.Plan,
		"updated_at": time.Now(),
	}

	if err := handler.Models.Team.Update(teamObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "team`s plan updated successfully")
}

func (handler *Handler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)

	teamIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamObjectId, err := utils.ToObjectID(teamIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamInstance, err := handler.Models.Team.Get(teamObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if teamInstance.OwnerId != userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "only the team owner can delete it")
		return
	}

	if err := handler.Models.Team.Delete(teamObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "team deleted successfully")
}

func (handler *Handler) AddUserToTeam(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		UserId string `json:"user_id"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	adminObjectId, err := utils.ToObjectID(payload.UserId)

	teamIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamObjectId, err := utils.ToObjectID(teamIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamInstance, err := handler.Models.Team.Get(teamObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !slices.Contains(teamInstance.Admins, adminObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "only the team admins can add users")
		return
	}

	newUserObjectId, err := utils.ToObjectID(input.UserId)

	if err := handler.Models.User.CheckExist(newUserObjectId); err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	// adding the new user
	currentUsers := teamInstance.Users
	currentUsers = append(currentUsers, newUserObjectId)

	updates := bson.M{
		"users":      currentUsers,
		"updated_at": time.Now(),
	}

	if err := handler.Models.Team.Update(teamObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "team users updated successfully")
}

func (handler *Handler) UploadTeamFile(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	teamIdStr, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamObjectId, err := utils.ToObjectID(teamIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	teamInstance, err := handler.Models.Team.Get(teamObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileName := r.FormValue("name")
	if fileName == "" {
		fileName = rand.Text()
	}

	maxUploadSize := utils.GetTeamMaxUploadSize(teamInstance.Plan)

	uploadDir := utils.GetTeamUploadDir(teamIdStr)

	fileAddress, totalUploadSize, err := handler.storeTeamFile(r, maxUploadSize, teamIdStr, teamInstance.Plan, uploadDir)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	expireAt := utils.GetTeamExpirationDate(teamInstance.Plan)

	if _, err := handler.Models.File.Create(userObjectId, teamObjectId, fileName, fileAddress, expireAt); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("creating file instance: %w", err))
		return
	}

	updates := bson.M{
		"storage_used": totalUploadSize,
		"updated_at":   time.Now(),
	}

	if err := handler.Models.Team.Update(teamObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("updating team instance: %w", err))
		return
	}

	utils.WriteJSON(w, "file uploaded successfully")
}

func (handler *Handler) storeTeamFile(r *http.Request, maxUploadSize int64, id, plan, uploadDir string) (string, int64, error) {
	allowedTypes := []string{"image/jpeg", "image/png", "application/zip"}

	file, err := utils.ReadFile(r, maxUploadSize, allowedTypes)
	if err != nil {
		return "", 0, err
	}

	defer file.File.Close()

	totalUsedStorage, err := handler.isTeamEligibleToUpload(id, plan, file.Size)
	if err != nil {
		return "", 0, err
	}

	fileAddress, err := file.UploadToDisk(uploadDir)
	if err != nil {
		return "", 0, err
	}

	return fileAddress, file.Size + totalUsedStorage, nil
}

func (handler *Handler) isTeamEligibleToUpload(id, plan string, fileSize int64) (int64, error) {
	totalStorage, err := utils.GetTeamTotalStorage(plan)
	if err != nil {
		return 0, err
	}

	if fileSize > totalStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your plan's total storage limit (%d bytes)", fileSize, totalStorage)
	}

	usedStorage, err := handler.getTeamUsedStorage(id)
	if err != nil {
		return 0, fmt.Errorf("failed to check used storage: %w", err)
	}

	remainedStorage := totalStorage - usedStorage
	if fileSize > remainedStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your remaining storage (%d bytes)", fileSize, remainedStorage)
	}

	newTotalStorage := fileSize + usedStorage
	return newTotalStorage, nil
}

func (handler *Handler) getTeamUsedStorage(id string) (int64, error) {
	if id == "" {
		return 0, errors.New("user id is missing")
	}

	userObjectId, err := utils.ToObjectID(id)
	if err != nil {
		return 0, err
	}

	return handler.Models.Team.GetUsedStorage(userObjectId)
}
