package types

type CommentInfo struct {
	ID              uint             `json:"id"`
	User            UserInfoResponse `json:"user"`
	Content         string           `json:"content"`
	CreatedDateBase string           `json:"created_db"`
}
type CommentListResponse struct {
	CommentList []CommentInfo `json:"comment_list"`
}
type UserListResponse struct {
	UserList []UserInfoResponse `json:"user_list"`
}
