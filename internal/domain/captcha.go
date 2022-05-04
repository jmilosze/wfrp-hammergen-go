package domain

type CaptchaService interface {
	Verify(captcha string) bool
}
