package dto

type ApplyLeaveRequest struct {
	LeaveType string `json:"leave_type" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
	StartDate Date   `json:"start_date" binding:"required"`
	EndDate   Date   `json:"end_date" binding:"required"`
}

type ApproveLeaveRequest struct {
	Approved bool   `json:"approved" binding:"required"`
	Remarks  string `json:"remarks,omitempty"`
}



