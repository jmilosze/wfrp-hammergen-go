package mockcaptcha

type CaptchaService struct{}

func NewCaptchaService() *CaptchaService {
	return &CaptchaService{}
}

func (e *CaptchaService) Verify(captcha string) bool {
	if captcha == "success" {
		return true
	}
	return false
}
