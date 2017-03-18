package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/fatih/color"
)

const (
	windowsExtension = ".exe"
	buildPath        = "dist"
)

var (
	project string
	pwd     string

	currentOS             string
	currentArchchitecture string

	architectures = []string{"amd64", "386"}
	systems       = []string{"darwin", "linux", "windows"}

	// user specified system to target
	target string
)

func init() {
	// Record the environment variables before proceeding
	getFromEnvironement()

	// Split and store paths for later use
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd, project = path.Split(currentPath)
}

func main() {
	flag.StringVar(&target, "for", "", "builder -for linux")
	flag.Parse()

	// only pass in current target
	if target != "" && isSupported(target) {
		systems = []string{target}
	}

	clearBuilds()

	color.Green("%s", fmt.Sprintf("Starting build in:\n%s%s", pwd, project))

	var wg sync.WaitGroup
	for _, targetSystem := range systems {
		for _, targetArch := range architectures {
			wg.Add(1)
			go performBuild(&wg, targetSystem, targetArch)
			wg.Wait()
		}
	}
	// reset the environment before exiting
	setEnvironement(currentOS, currentArchchitecture)

	notice := color.GreenString("Done!\nYou will your build under the '%s' folder", buildPath)
	fmt.Println(notice)
}

func isSupported(target string) bool {
	for _, sys := range systems {
		if target == sys {
			return true
		}
	}
	return false
}

// clearBuilds removes the old builds before starting a new one
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

	err = executeGoBuild()
	if err != nil {
		fmt.Println("Error running build command", err)
		fmt.Println("Make sure you are running this tool where your main.go is located!")
		return
	}

	// I could use os.Rename, but linking and removing after is safer...
	if o == "windows" {
		filename := fmt.Sprintf("%s%s", project, windowsExtension)
		os.Link(filename, fmt.Sprintf("./%s/%s/%s%s", buildPath, platform, project, windowsExtension))
		os.Remove(filename)
	} else {
		os.Link(project, fmt.Sprintf("./%s/%s/%s", buildPath, platform, project))
		os.Remove(project)
	}
}

func getFromEnvironement() {
	currentOS = os.Getenv("GOOS")
	currentArchchitecture = os.Getenv("GOARCH")
}

func setEnvironement(system, architecture string) {
	os.Setenv("GOOS", system)
	os.Setenv("GOARCH", architecture)
}

func executeGoBuild() error {
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	return err
}
