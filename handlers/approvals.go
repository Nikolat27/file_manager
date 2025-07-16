package handlers

import (
	"file_manager/utils"
	"net/http"
)

func (handler *Handler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var input struct {
		FileId  string `json:"file_id"`
		OwnerId string `json:"owner_id"`
		Reason  string `json:"reason"`
	}

	if err := utils.ParseJsonData(r.Body, 10000, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	fileObjectId, err := utils.ConvertStringToObjectID(input.FileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isRequired, err := handler.Models.File.CheckFileRequiresApproval(fileObjectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isRequired {
		http.Error(w, "this file does not require approval", http.StatusBadRequest)
		return
	}

	ownerObjectId, err := utils.ConvertStringToObjectID(input.OwnerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userObjectId, err := utils.ConvertStringToObjectID(payload.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := handler.Models.Approval.CreateApprovalInstance(fileObjectId, ownerObjectId, userObjectId, input.Reason); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Your Approval request has been sent successfully"))
}

func (handler *Handler) ChangeApprovalStatus(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r, handler.PasetoMaker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var input struct {
		ApprovalId string `json:"approval_id"`
		Status     string `json:"status"`
	}

	if err := utils.ParseJsonData(r.Body, 10000, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userObjectId, err := utils.ConvertStringToObjectID(payload.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approvalObjectId, err := utils.ConvertStringToObjectID(input.ApprovalId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.Models.Approval.ValidateApprovalOwner(approvalObjectId, userObjectId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.Models.Approval.UpdateApprovalStatus(approvalObjectId, input.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Approval`s Status changed successfully"))
}
