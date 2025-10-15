package model

type UserRegisterReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserResetPasswordReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
