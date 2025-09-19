package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	file, err := os.Open(".runway")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var commands []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Fatalln(err)
	}

	errs := make(chan error, len(commands))

	for _, command := range commands {
		fields := strings.Fields(command)
		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

		go func(cmd *exec.Cmd) {
			err := cmd.Start()
			if err != nil {
				errs <- err
				return
			}
			errs <- cmd.Wait()
		}(cmd)
	}

	for range commands {
		err := <-errs
		if err != nil {
			fmt.Println(err)
		}
	}
}
