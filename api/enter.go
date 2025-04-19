package api

type ApiGroup struct {
	BaseApi
	UserApi
	ImageApi
	ArticleApi
	CommentApi
	AdvertisementApi
	FriendLinkApi
	FeedbackApi
	WebsiteApi
	ConfigApi
}

var ApiGroupApp = new(ApiGroup)
