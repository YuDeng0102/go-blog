package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"
	"time"
)

type UserApi struct{}

// Register 注册
func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		global.Log.Info("savedEmail not match", zap.Any("savedEmail", savedEmail), zap.Any("req.email", req.Email))
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verificationCode")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		global.Log.Info("verification code not match", zap.Any("savedCode", savedCode), zap.Any("req.VerificationCode", req.VerificationCode))
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expiretime")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}
	user, err := service.ServiceGroupApp.UserService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	userApi.TokenNext(c, user)

}

func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "qq":
	default:
		userApi.EmailLogin(c)
	}
}

func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := service.ServiceGroupApp.UserService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login by email:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}
		userApi.TokenNext(c, user)
	}
}

func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	if user.Freeze {
		response.FailWithMessage("The user has been frozen", c)
		return
	}

	baseClaim := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	//创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaim)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	//创建刷新令牌
	refreshClaim := j.CreateRefreshClaims(baseClaim)
	refreshToken, err := j.CreateRefreshToken(refreshClaim)
	if err != nil {
		global.Log.Error("Failed to get refresh token:", zap.Error(err))
		response.FailWithMessage("Failed to get refresh token", c)
		return
	}

	// 是否开启了多地点登录拦截
	if !global.Config.System.UseMultipoint {
		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaim.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}

	if jwtStr, err := service.ServiceGroupApp.JwtService.GetRedisJWT(user.UUID); errors.Is(err, redis.Nil) {
		if err := service.ServiceGroupApp.JwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
			return
		}

		//设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaim.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	} else if err != nil {
		global.Log.Error("Failed to get login status:", zap.Error(err))
		response.FailWithMessage("Failed to get login status", c)
	} else {
		//Redis 中已经存在该用户的JWT，将旧的JWT加入黑名单（限制异地多账户登录）

		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := service.ServiceGroupApp.JwtService.JoinInBlacklist(blacklist); err != nil {
			global.Log.Error("Failed to invalidata jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", c)
			return
		}

		//设置新的JWT到service
		if err := service.ServiceGroupApp.JwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaim.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	}

}

func (UserApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		global.Log.Info("savedEmail not match", zap.Any("savedEmail", savedEmail), zap.Any("req.email", req.Email))
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}
	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verificationCode")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		global.Log.Info("verification code not match", zap.Any("savedCode", savedCode), zap.Any("req.VerificationCode", req.VerificationCode))
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expiretime")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	if err := service.ServiceGroupApp.UserService.ForgotPassword(req); err != nil {
		global.Log.Error("ForgotPassword failed", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully retrieved", c)
}

func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userCard, err := service.ServiceGroupApp.UserService.UserCard(req.UUID)
	if err != nil {
		global.Log.Error("UserCard failed", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(userCard, c)
}

// Logout 登出
func (userApi *UserApi) Logout(c *gin.Context) {
	service.ServiceGroupApp.UserService.Logout(c)
	response.OkWithMessage("Successful logout", c)
}

// UserResetPassword 修改密码
func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = service.ServiceGroupApp.UserService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("Failed to modify, orginal password does not match the current account", c)
		return
	}
	response.OkWithMessage("Successfully changed password, please log in again", c)
	service.ServiceGroupApp.UserService.Logout(c)
}

// UserInfo 获取个人信息
func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := service.ServiceGroupApp.UserService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(user, c)
}

// UserChangeInfo 修改个人信息
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = service.ServiceGroupApp.UserService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("Successfully changed user information", c)
}

// UserWeather 获取天气
func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	//ip := "112.2.249.25"
	weather, err := service.ServiceGroupApp.UserService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据，登录和注册人数
func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := service.ServiceGroupApp.UserService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to user chart", c)
		return
	}
	response.OkWithData(data, c)
}

//// UserList 获取用户列表
//func (userApi *UserApi) UserList(c *gin.Context) {
//	var pageInfo request.UserList
//	err := c.ShouldBindQuery(&pageInfo)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//
//	list, total, err := userService.UserList(pageInfo)
//	if err != nil {
//		global.Log.Error("Failed to get user list:", zap.Error(err))
//		response.FailWithMessage("Failed to get user list", c)
//		return
//	}
//	response.OkWithData(response.PageResult{
//		List:  list,
//		Total: total,
//	}, c)
//}
