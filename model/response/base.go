package response

type Captcha struct {
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}
