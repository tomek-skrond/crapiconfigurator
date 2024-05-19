package crapiconfigurator

type Config struct {
	TargetURL string `json:"target_url"`
	LoginURL  string `json:"login_url"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
