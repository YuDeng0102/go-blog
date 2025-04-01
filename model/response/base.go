package response

type Captcha struct {
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
}
