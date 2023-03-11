package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const configFileName string = "dbsize.yaml"

func main() {
	log.Print("start DbSize")

	config, err := NewConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration: %#v", config)

	webUI := NewWebUI(config)
	err = webUI.Start()
	if err != nil {
		log.Fatal(err)
	}

	app := fmt.Sprintf("--app=http://localhost:%v", config.WebUI.Port)
	user_data_dir := fmt.Sprintf("--user-data-dir=%v/DbSize/Chrome", os.Getenv("localappdata"))
	cmd := exec.Command(config.WebUI.ChromePath, app, user_data_dir)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exitSignal
		cmd.Process.Kill()
		os.Exit(0)
	}()

	cmd.Wait()
	webUI.Stop()
}
