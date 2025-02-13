package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

const project_name = "aqueduct"

var root = fmt.Sprintf("/etc/%s", project_name)
var template_root = fmt.Sprintf("%s/public/templates", root)
var assets_root = fmt.Sprintf("%s/public/assets", root)

func main() {
	dev_mode := flag.Bool("dev", false, "Sets the program to development mode, disabling help messages and proceding without env")

	flag.Parse()

	// Listen for sigterm signals | finishes program if finds one
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig

		os.Exit(1)
	}()

	if os.Getenv("GIN_MODE") != "release" {
		if !*dev_mode {
			fmt.Printf("-- Local do arquivo de configurações: %s\n", root)
			fmt.Printf("-- Para parar o programa (caso esteja funcionando em segundo plano): \n")
			fmt.Println()
			fmt.Printf("	$ sudo systemctl stop %s.service \n", project_name)
			fmt.Printf("	$ sudo systemctl disable %s.service \n", project_name)

			os.Exit(0)
		}
	}

	var path string = GetConfig("path")
	var port string = GetConfig("port")

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
	router.GET("/delete", HandleDelete)
	router.GET("/download", HandleDownload)
	// Note: For folder download

	// Enable user sessions to delete data uploaded in session
	router.POST("/cookie", HandleCookie)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		LogFatal(err.Error())
	}
}

func GetConfig(key string) (value string) {
	if val := os.Getenv(key); val != "" {
		return val
	}

	LogFatal(
		fmt.Sprintf(
			"Variavel de configuração (%s) necessária não encontrada.",
			key,
		),
	)

	return ""
}
