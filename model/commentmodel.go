package model

type Comment struct {
	ID          int    `json:"id" form:"id"`
	Username    string `json:"username" form:"username"`
	Nickname    string `json:"nickname" form:"nickname"`
	ProductID   int    `json:"product_id" form:"product_id"`
	ProductName string `json:"product_name" form:"product_name"`
	Comment     string `json:"comment" form:"comment"`
	Time        string `json:"time" form:"time"`
	ParentID    int    `json:"parent_id" form:"parent_id" `
}
type NewComment struct {
	Comment  Comment
	SonCount int `json:"son_count"`
}
