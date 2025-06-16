package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/Noeeekr/aqueduct/internal"
	"github.com/gin-gonic/gin"
)

type Server struct {
}

func Start() error {
	var instance *internal.Instance = internal.NewInstance()
	SharedFolder := instance.Info.SharedFolder
	port := instance.Info.Port

	if SharedFolder == "" {
		return errors.New(" O caminho da pasta de compartilhamento não foi definido")
	}

	// Check if shared folder exists
	if dir, err := os.Stat(SharedFolder); err != nil {
		return errors.New(" O caminho da pasta de compartilhamento não existe ou não foi encontrado: " + err.Error())
	} else {
		if !dir.IsDir() {
			return errors.New(" Caminho de diretório invalido. Precisa ser uma pasta. ")
		}
	}

	if instance.Info.Environment == internal.EnviromentProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	handlers := NewHandlers()

	// Handles assets
	// Note: etc/(project-name)/public/templates/*
	router.LoadHTMLGlob(fmt.Sprintf("%s/*", instance.Info.TemplateFolder))
	// Note: etc/(project-name/public/assets/*
	router.Static("/assets", instance.Info.AssetsFolder)
	router.Static("/files", SharedFolder)
	router.GET("/", handlers.ServeTemplate)

	// Handles CRUD
	router.POST("/upload", handlers.HandleUpload)
	router.GET("/delete", handlers.HandleDelete)
	router.GET("/download", handlers.HandleDownload)
	// Note: For folder download

	// Enable user sessions to delete data uploaded in session
	router.POST("/cookie", handlers.HandleCookie)

	fmt.Println("Starting server on port " + port)
	return router.Run(fmt.Sprintf(":%s", port))
}
