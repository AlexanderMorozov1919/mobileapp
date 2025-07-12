package models

type DoctorLoginRequest struct {
	Login    string `json:"login" binding:"required" example:"+79123456789"`
	Password string `json:"password" binding:"required"`
}

type DoctorAuthResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Token string `json:"token"`
}
