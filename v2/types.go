package crapiconfigurator

type Config struct {
	GlobalConfig
	TargetURL string `json:"target_url"`
}

type GlobalConfig struct {
	Hostname string `json:"hostname" yaml:"hostname"`
	Email    string `json:"email" yaml:"email"`
	Password string `json:"password" yaml:"password"`
	LoginURL string `json:"login_url" yaml:"login_url"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
