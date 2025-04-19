package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	ImageService
	ArticleService
	CommentService
	AdvertisementService
	FriendLinkService
	FeedbackService
	WebsiteService
	HotSearchService
	CalendarService
	ConfigService
}

var ServiceGroupApp = new(ServiceGroup)
