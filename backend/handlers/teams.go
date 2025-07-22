package handlers

import (
	"crypto/rand"
	"encoding/json"
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

// GetTeams -> Returns List
func (handler *Handler) GetTeams(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)

	filter := bson.M{
		"users": userObjectId,
	}

	teams, err := handler.Models.Team.GetAll(filter, bson.M{}, 1, 6)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	data, err := json.MarshalIndent(teams, "", "\t")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, data)
}

// GetTeam -> Returns One
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

	filter := bson.M{
		"_id": teamObjectId,
	}

	teamInstance, err := handler.Models.Team.Get(filter, bson.M{})
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

	filter := bson.M{
		"_id": teamObjectId,
	}

	projection := bson.M{
		"owner_id": 1,
	}

	teamInstance, err := handler.Models.Team.Get(filter, projection)
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

	filter := bson.M{
		"_id": teamObjectId,
	}

	projection := bson.M{
		"owner_id": 1,
	}

	teamInstance, err := handler.Models.Team.Get(filter, projection)
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
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

	filter := bson.M{
		"_id": teamObjectId,
	}

	projection := bson.M{
		"admins": 1,
		"users":  1,
	}

	teamInstance, err := handler.Models.Team.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !slices.Contains(teamInstance.Admins, adminObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "only the team admins can add users")
		return
	}

	userObjectId, err := utils.ToObjectID(input.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"_id": userObjectId,
	}

	projection = bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.User.Get(filter, projection); err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if slices.Contains(teamInstance.Users, userObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "This user is already in the team")
		return
	}

	// adding the new user
	currentUsers := teamInstance.Users
	if err := canAddUser(len(currentUsers), teamInstance.Plan); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// users updated
	currentUsers = append(currentUsers, userObjectId)

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

	filter := bson.M{
		"_id": teamObjectId,
	}

	projection := bson.M{
		"plan":         1,
		"storage_used": 1,
	}

	teamInstance, err := handler.Models.Team.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileName := r.FormValue("name")
	if fileName == "" {
		fileName = rand.Text()
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	folderObjectId, err := getFolderId(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if folderObjectId != primitive.NilObjectID {
		if err := handler.ValidateFolderId(folderObjectId, userObjectId); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	maxUploadSize := utils.GetTeamMaxUploadSize(teamInstance.Plan)

	uploadDir := utils.GetTeamUploadDir(teamIdStr)

	fileAddress, totalUploadSize, err := handler.storeTeamFile(r, maxUploadSize, teamInstance.StorageUsed, teamInstance.Plan, uploadDir)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	expireAt := utils.GetTeamExpirationDate(teamInstance.Plan)

	if _, err := handler.Models.File.Create(userObjectId, teamObjectId, folderObjectId, fileName, fileAddress, expireAt); err != nil {
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

func (handler *Handler) storeTeamFile(r *http.Request, maxUploadSize, totalUsedStorage int64, plan, uploadDir string) (string, int64, error) {
	allowedTypes := []string{"image/jpeg", "image/png", "application/zip"}

	file, err := utils.ReadFile(r, maxUploadSize, allowedTypes)
	if err != nil {
		return "", 0, err
	}

	defer file.File.Close()

	newTotalStorage, err := handler.isTeamEligibleToUpload(plan, totalUsedStorage, file.Size)
	if err != nil {
		return "", 0, err
	}

	fileAddress, err := file.UploadToDisk(uploadDir)
	if err != nil {
		return "", 0, err
	}

	return fileAddress, file.Size + newTotalStorage, nil
}

func (handler *Handler) isTeamEligibleToUpload(plan string, usedStorage, fileSize int64) (int64, error) {
	totalStorage, err := utils.GetTeamTotalStorage(plan)
	if err != nil {
		return 0, err
	}

	if fileSize > totalStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your plan's total storage limit (%d bytes)", fileSize, totalStorage)
	}

	remainedStorage := totalStorage - usedStorage
	if fileSize > remainedStorage {
		return 0, fmt.Errorf("your file size (%d bytes) exceeds your remaining storage (%d bytes)", fileSize, remainedStorage)
	}

	newTotalStorage := fileSize + usedStorage
	return newTotalStorage, nil
}

func canAddUser(currentUsersAmount int, teamPlan string) error {
	totalAllowedUsers := utils.GetTeamTotalUsers(teamPlan)
	if totalAllowedUsers == -1 {
		return nil
	}

	if currentUsersAmount+1 < totalAllowedUsers {
		return nil
	}

	return fmt.Errorf("you have exceed you total user adding limitation. For free plan is :%d", totalAllowedUsers)
}
