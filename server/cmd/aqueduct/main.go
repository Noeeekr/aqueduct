package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Noeeekr/aqueduct/internal"
	"github.com/Noeeekr/aqueduct/internal/server"
)

func main() {

	var instance *internal.Instance = internal.NewInstance()

	LogsFolder, err := os.Open(instance.Info.Logsfolder)
	if err != nil {
		LogsFolder = os.Stdout
	}
	instance.Logger.SetOutput(LogsFolder)

	// Listen for sigterm signals | finishes program if finds one
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig

		os.Exit(1)
	}()

	if instance.Info.Environment == internal.EnvironmentDevelopment {
		fmt.Println("Arquivo de configurações: " + instance.Info.ProgramFolder)
		fmt.Println("Parar o programa, caso esteja funcionando em segundo plano: ")
		fmt.Println("	$ sudo systemctl stop aqueduct")
		fmt.Println("Iniciar o programa em segundo plano: ")
		fmt.Println("	$ sudo systemctl start aqueduct")
		fmt.Println("Iniciar em segundo plano sempre que iniciar o computador")
		fmt.Println("	$ sudo systemctl enable aqueduct")
		os.Exit(0)
	}

	err = server.Start()
	if err != nil {
		instance.Logger.Error(err.Error())
		return
	}
}
