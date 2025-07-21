package handlers

import (
	"file_manager/utils"
	"fmt"
	"net/http"
)

func (handler *Handler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	var input struct {
		FileId string `json:"file_id"`
		Reason string `json:"reason"`
	}

	if err := utils.ParseJSON(r.Body, 10000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	fileObjectId, err := utils.ToObjectID(input.FileId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	isRequired, err := handler.Models.FileSettings.IsApprovalRequired(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !isRequired {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("this file does not require approval"))
		return
	}

	ownerObjectId, err := handler.Models.FileSettings.GetUserId(fileObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userObjectId, err := utils.ToObjectID(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	alreadyRequested, err := handler.Models.Approval.HasAlreadyRequested(fileObjectId, userObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if alreadyRequested {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you have already sent an approval request for this file. Please be patient"))
		return
	}

	if _, err := handler.Models.Approval.Create(fileObjectId, ownerObjectId, userObjectId, input.Reason); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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

	if err := handler.Models.Approval.ValidateOwner(approvalObjectId, userObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := handler.Models.Approval.UpdateStatus(approvalObjectId, input.Status); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, "Approval`s Status changed successfully")
}
