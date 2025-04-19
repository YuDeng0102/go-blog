package router

type RouterGroup struct {
	baseRouter
	UserRouter
	ImageRouter
	ArticleRouter
	CommentRouter
	AdvertisementRouter
	FriendLinkRouter
	FeedbackRouter
	ConfigRouter
	WebsiteRouter
}

var RouterGroupApp = new(RouterGroup)
