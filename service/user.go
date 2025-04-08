package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
)

type UserService struct{}

func (userService *UserService) Register(u database.User) (database.User, error) {
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}

	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "./images/avatar.png"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email
	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}
	return u, nil
}

func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("invalid password")
		}
		return user, nil
	}
	return database.User{}, err
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	err := global.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return errors.New("invalid email")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(uuid string) (response.UserCard, error) {
	var user database.User
	err := global.DB.Where("uuid=?", uuid).First(&user).Error
	if err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userService *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	global.Redis.Del(uuid.String())
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	err := global.DB.Where("password=?", utils.BcryptHash(req.Password)).First(&user)
	if err != nil || user.ID != req.UserID {
		return errors.New("invalid password")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}
func (userService *UserService) UserWeather(ip string) (string, error) {
	// 从redis中获取天气数据，如果没有数据，则调用高德api进行查询
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil {
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}

		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		// 将天气数据存入redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	// 构建查询条件
	where := global.DB.Where(fmt.Sprintf("date_sub(curdate(), interval %d day) <= created_at", req.Date))

	var res response.UserChart

	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		res.DateList = append(res.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	// 获取注册数据
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)

	for _, date := range res.DateList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}

	return res, nil
}
