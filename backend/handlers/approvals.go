package handlers

import (
	"errors"
	"file_manager/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (handler *Handler) GetApprovalsList(w http.ResponseWriter, r *http.Request) {
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
		"sender_id": userObjectId,
	}

	approvals, err := handler.Models.Approval.GetAll(filter, bson.M{}, 1, 10)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			response := map[string]any{
				"approvals": nil,
			}

			utils.WriteJSONData(w, response)
			return
		}

		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response := map[string]any{
		"approvals": approvals,
	}

	utils.WriteJSONData(w, response)
}

func (handler *Handler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		ShortUrl string `json:"short_url"`
		Reason   string `json:"reason"`
	}

	if err := utils.ParseJSON(r.Body, 10000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"short_url": input.ShortUrl,
	}

	projection := bson.M{
		"user_id":    1,
		"name":       1,
		"approvable": 1,
		"file_id":    1,
	}

	fileSettings, err := handler.Models.FileSettings.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !fileSettings.Approvable {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("this file does not require approval"))
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"file_id":   fileSettings.FileId,
		"sender_id": userObjectId,
	}

	projection = bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Approval.Get(filter, projection); err == nil {
		utils.WriteError(w, http.StatusBadRequest, "approval request has already sent")
		return
	} else {
		// If approval exists, Return error
		if !errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	filter = bson.M{
		"_id": fileSettings.FileId,
	}

	projection = bson.M{
		"name": 1,
	}

	file, err := handler.Models.File.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := handler.Models.Approval.Create(fileSettings.FileId, fileSettings.UserId, userObjectId, file.Name, input.Reason); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, "Your Approval request has been sent successfully")
}

func (handler *Handler) UpdateApproval(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		ApprovalId string `json:"approval_id"`
		Status     string `json:"status"`
	}

	if err := utils.ParseJSON(r.Body, 10000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	approvalObjectId, err := utils.ToObjectID(input.ApprovalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id": approvalObjectId,
	}

	projection := bson.M{
		"owner_id": 1,
	}

	approvalInstance, err := handler.Models.Approval.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if approvalInstance.OwnerId != userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "this user is not the approval`s owner")
		return
	}

	updates := bson.M{
		"status":      input.Status,
		"reviewed_at": time.Now(),
	}

	if err := handler.Models.Approval.Update(approvalObjectId, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "Approval`s Status changed successfully")
}

func (handler *Handler) CheckApproval(w http.ResponseWriter, r *http.Request) {
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

	shortUrl, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"short_url": shortUrl,
	}

	projection := bson.M{
		"_id":     1,
		"file_id": 1,
	}

	settingsInstance, err := handler.Models.FileSettings.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"_id": settingsInstance.FileId,
	}

	projection = bson.M{
		"owner_id": 1,
	}

	file, err := handler.Models.File.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := checkApprovalAccess(settingsInstance.FileId, file.OwnerId, userObjectId, handler); err != nil {
		var approvalError *utils.ApprovalError

		if errors.As(err, &approvalError) {
			utils.WriteError(w, http.StatusPreconditionRequired, approvalError.Type)
			return
		}

		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "approved")
}

func (handler *Handler) DeleteApproval(w http.ResponseWriter, r *http.Request) {
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

	approvalId, err := utils.ParseIdParam(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	approvalObjectId, err := utils.ToObjectID(approvalId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter := bson.M{
		"_id":       approvalObjectId,
		"sender_id": userObjectId,
	}

	projection := bson.M{
		"_id": 1,
	}

	approvalInstance, err := handler.Models.Approval.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "this approval either does not exist or you are not it`s owner")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	filter = bson.M{
		"_id": approvalInstance.Id,
	}

	if err := handler.Models.Approval.DeleteOne(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "approval deleted successfully")
}

func checkApprovalAccess(fileId, ownerId, requesterId primitive.ObjectID, handler *Handler) error {
	if ownerId == requesterId {
		return nil
	}

	filter := bson.M{
		"file_id":   fileId,
		"sender_id": requesterId,
	}

	projection := bson.M{
		"status": 1,
	}

	approvalInstance, err := handler.Models.Approval.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &utils.ApprovalError{
				Type:    "not_requested",
				Message: "approval is required but not requested yet",
			}
		}
		return err
	}

	switch approvalInstance.Status {
	case "approved":
		return nil
	case "rejected":
		return &utils.ApprovalError{Type: "rejected", Message: "your approval has been rejected"}
	case "pending":
		return &utils.ApprovalError{Type: "pending", Message: "your approval is still pending"}
	default:
		return &utils.ApprovalError{Type: "error", Message: "your approval status is invalid"}
	}

}
