package models

type User struct {
	Id          int    `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"-"`
	Role        string `json:"user_role,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
