package internal

import (
	"flag"
	"os"
	"sync"

	env "github.com/joho/godotenv"
)

type Enviroment int

const (
	EnvironmentDevelopment Enviroment = iota
	EnviromentProduction
)

type Info struct {
	once sync.Once

	ProgramName string
	Port        string
	Environment Enviroment

	ProgramFolder  string
	TemplateFolder string
	AssetsFolder   string
	SharedFolder   string
	Logsfolder     string
}

var defaultInfo *Info = &Info{
	ProgramName: "aqueduct",
	Port:        "80",
	Environment: EnvironmentDevelopment,

	ProgramFolder:  "/etc/aqueduct",
	TemplateFolder: "/etc/aqueduct/public/templates",
	AssetsFolder:   "/etc/aqueduct/public/assets",
	SharedFolder:   "",
	Logsfolder:     "/var/log/aqueduct",
}

// NewInfo returns the same info everytime
func NewInfo() *Info {
	defaultInfo.once.Do(func() {
		var environment string

		flag.StringVar(&defaultInfo.SharedFolder, "public-folder", "", "sets the folder where files will be stored and shared")
		flag.StringVar(&defaultInfo.Port, "port", defaultInfo.Port, "sets the port the server will listen to. Default is 80")
		flag.StringVar(&defaultInfo.Logsfolder, "logs", defaultInfo.Logsfolder, "Sets the folder where logs will be stored")
		flag.StringVar(&environment, "environment", "", "Sets the environment the server will be using: development or production")
		flag.Parse()

		if environment == "production" {
			defaultInfo.Environment = EnviromentProduction
		} else {
			defaultInfo.Environment = EnvironmentDevelopment
		}

		err := env.Load(defaultInfo.ProgramFolder + "/config.env")
		if err == nil {
			SharedFolder, ok := os.LookupEnv("path")
			if ok {
				defaultInfo.SharedFolder = SharedFolder
			}
		}
	})

	return defaultInfo
}

// Make function to change certain properties
// For example Changin ProgramFIle will change all other properties related to it in a single go
