package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func validateRequirements() {
	execs := map[string][]string{
		"FF-MPEG":         {"which", "ffmpeg"},
		"Deezer Spleeter": {"which", "spleeter"},
	}
	for commandName, cmds := range execs {
		cmd := exec.Command(cmds[0], cmds[1:]...)
		_, err := cmd.CombinedOutput()
		if err != nil {
			panic(fmt.Sprintf("%s exec file is not found", commandName))
		}
	}

	// make important directories
	dir, _ := os.Getwd()
	CreateDirIfNotExist(fmt.Sprintf("%s/%s", dir, "media/output"))
	CreateDirIfNotExist(fmt.Sprintf("%s/%s", dir, "media/upload"))
}

func split(file, output string) error {
	if !fileExists(file) {
		return fmt.Errorf("File [%s] is not a valid file", file)
	}

	if !isDir(output) {
		return fmt.Errorf("Directory [%s] is not exists", output)
	}

	command := []string{"spleeter", "separate", "-i", file, "-o", output}
	cmd := exec.Command(command[0], command[1:]...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isDir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
