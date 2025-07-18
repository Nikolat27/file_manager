package handlers

import (
	"crypto/rand"
	"file_manager/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (handler *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if err := utils.ValidatePlan(payload.UserPlan); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.UserPlan != "premium" {
		utils.WriteError(w, http.StatusBadRequest, "Team creation is only available for 'premium' plan users")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
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

	if _, err := handler.Models.Teams.Create(teamId, userObjectId, name, description, avatarAddress); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "team created successfully")
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
