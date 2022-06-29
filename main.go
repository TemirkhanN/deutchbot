package main

import (
	"DeutchBot/cmd"
	"os"
)

func main() {
	cmd.StartBot(os.Getenv("TG_BOT_TOKEN"))

	/*
		projectDir, _ := os.Getwd()
		cmd.RegenerateTasks(projectDir)
	*/
}
