package types

type VideoInfo struct {
	ID            uint             `json:"id"`
	Author        UserInfoResponse `json:"author"`
	PlayURL       string           `json:"play_url"`
	CoverURL      string           `json:"cover_url"`
	FavoriteCount int64            `json:"favorite_count"`
	CommentCount  int64            `json:"comment_count"`
	IsFavourite   bool             `json:"is_favourite"`
	Title         string           `json:"title"`
}
type FeedResponse struct {
	NextTime  int64       `json:"next_time"`
	VideoList []VideoInfo `json:"video_list"`
}
type VideoListResponse struct {
	VideoList []VideoInfo `json:"video_list"`
}
