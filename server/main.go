package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const project_name = "aqueduct"

var root = fmt.Sprintf("/etc/%s", project_name)
var config_root = fmt.Sprintf("%s/config.conf", root)
var template_root = fmt.Sprintf("%s/public/templates", root)
var assets_root = fmt.Sprintf("%s/public/assets", root)

func main() {
	setup()

	var path string = GetConfig("path")
	var port string = GetConfig("port")
	var root, err = filepath.Abs("/etc/aqueduct")
	if err != nil {
		LogFatal(err.Error())
	}

	fmt.Println(root)
	if dir, err := os.Stat(path); err != nil {
		LogInfo(fmt.Sprintf("PATH=%s", path))
		LogFatal("O valor do paramêtro PATH definido no arquivo de configuração não existe ou o programa não tem acesso de leitura ao caminho.")
	} else {
		if !dir.IsDir() {
			LogFatal("Caminho de diretório invalido. Precisa ser uma pasta.")
		}
	}

	router := gin.Default()

	// Handles assets
	// Note: etc/(project-name)/public/templates/*
	router.LoadHTMLGlob(fmt.Sprintf("%s/*", template_root))
	// Note: etc/(project-name/public/assets/*
	router.Static("/assets", assets_root)
	router.Static("/files", path)
	router.GET("/", ServeTemplate)

	// Handles CRUD
	router.POST("/upload", HandleUpload)
	router.POST("/delete", HandleDelete)
	router.GET("/download/:folder", HandleDownload)
	// Note: For folder download

	// Enable user sessions to delete data uploaded in session
	router.POST("/cookie", HandleCookie)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		LogFatal(err.Error())
	}
}
