package model

type UserInfoResp struct {
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	UserCreateTime int64  `json:"user_create_time"`
}
