package bootstrap

type Application struct {
	Config *AppConfig
}

func New(configPath string) *Application {
	cfg := Load(configPath)
	return &Application{
		Config: cfg,
	}
}
