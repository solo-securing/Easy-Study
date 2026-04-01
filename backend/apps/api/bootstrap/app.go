package bootstrap

import "api/internal/config"

type Application struct {
	Config *config.AppConfig
}

func App(configPath string) *Application {
	var app = &Application{}

	app.Config = config.Load(configPath)

	return app
}
