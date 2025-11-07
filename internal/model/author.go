package model

type Author struct {
	ID             string  `json:"id,omitempty"`
	FullName       *string `json:"full_name"`
	NickName       string  `json:"nick_name,omitempty"`
	Specialization string  `json:"specialization,omitempty"`
}
