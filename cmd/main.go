package main

import (
	"github.com/yervsil/auth_service/internal/configs"
	"github.com/yervsil/auth_service/internal/app"
	_ "github.com/yervsil/auth_service/docs" 
)

// @title           Auth service API
// @version         1.0
// @description     Your API description
// @host            localhost:8000
// @BasePath        /
func main() {
	cfg, err := configs.Init()
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}

