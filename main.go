package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"github.com/fatih/color"
)

const winExt = ".exe"
const buildPath = "dist"

var project string
var pwd string

var currOs string
var currArch string

func init() {
	// Record the environment variables before proceeding
	getEnvironement()

	// Split and store paths for later use
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd, project = path.Split(currentPath)
}

func main() {
	var wg sync.WaitGroup
	arch := []string{"amd64", "386"}
	syst := []string{"darwin", "linux", "windows"}

	clearBuilds()

	color.Green("%s" ,fmt.Sprintf("Starting build in:\n%s%s", pwd, project))

	for _, o := range syst {
		for _, a := range arch {
			wg.Add(1)
			go performBuild(&wg, o, a)
			wg.Wait()
		}
	}
	// reset the environment before exiting
	setEnvironement(currOs, currArch)

	notice := color.GreenString("Done!\nYou will your build under the '%s' folder", buildPath)
	fmt.Println(notice)
}

func clearBuilds() {
	_, err := os.Stat(buildPath)
	if err != nil {
		return
	}

	fmt.Print("Clearing old builds...")

	err = os.RemoveAll(buildPath)
	if err != nil {
		panic(err)
	}

	fmt.Print(" Success.\n")
}

func performBuild(wg *sync.WaitGroup, o, a string) {
	defer wg.Done()

	platform := fmt.Sprintf("%s_%s", o, a)
	folderPath := fmt.Sprintf("%s/%s", buildPath, platform)

	fmt.Println(fmt.Sprintf("Building %s for %s", project, platform))

	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		fmt.Println("Error creating directories: ", err)
		return
	}

	// Set the environment to the currently targeted build
	setEnvironement(o, a)

	err = executeBuild()
	if err != nil {
		fmt.Println("Error running build command", err)
		fmt.Println("Make sure you are running this tool where your main.go is located!")
		return
	}

	// I could use os.Rename, but linking and removing after is safer...
	if o == "windows" {
		filename := fmt.Sprintf("%s%s", project, winExt)
		os.Link(filename, fmt.Sprintf("./%s/%s/%s%s", buildPath, platform, project, winExt))
		os.Remove(filename)
	} else {
		os.Link(project, fmt.Sprintf("./%s/%s/%s", buildPath, platform, project))
		os.Remove(project)
	}
}

func getEnvironement() {
	currOs = os.Getenv("GOOS")
	currArch = os.Getenv("GOARCH")
}

func setEnvironement(system, architecture string) {
	os.Setenv("GOOS", system)
	os.Setenv("GOARCH", architecture)
}

func executeBuild() error {
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
