package main

import (
	"fmt"
	"log"
	"os/exec"
)

func RunProgram(project string, does_auto_close bool) {
	log.Println("Running Program...")

	file := "main"
	str_cmd := fmt.Sprintf("pushd %s && go build -o %s && popd", project, file)
	cmd := exec.Command("bash", "-c", str_cmd)
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Go Build Failed: %s --- %s\n --- %s\n", project, err.Error(), string(stdout))
		return
	}
	log.Printf("Go Build Success: %s", cmd.String())

	eop := "&& "
	switch does_auto_close {
	case true:
		seconds := 3
		eop += fmt.Sprintf("sleep %d", seconds)
	case false:
		eop += "read -p 'press a key to exit...'"
	}
	str_cmd = fmt.Sprintf("gnome-terminal --geometry=136x43 -- bash -c \"%s/%s %s\"", project, file, eop)
	cmd = exec.Command("bash", "-c", str_cmd)
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to run program: %s --- %s\n", project+"/main", err.Error())
		return
	}
	log.Printf("Successfully ran program: %s\n --- %s\n", cmd.String(), string(stdout))
}
