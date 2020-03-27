//+build mage

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var exname string = ""

func getExecutableName(os string, arch string) string {
	if exname == "" {
		var err error
		exname, err = getExecutableFromPluginJSON()
		if err != nil {
			exname = "set_exe_name_in_plugin_json" // warning in the final name?
		}
	}

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if "windows" == os {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName
}

func getExecutableFromPluginJSON() (string, error) {
	byteValue, err := readFileBytes(path.Join("src", "plugin.json"))
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return "", err
	}
	return result["executable"].(string), nil
}

func readFileBytes(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = jsonFile.Close()
	}()

	return ioutil.ReadAll(jsonFile)
}

func findRunningPIDs(exe string) []int {
	pids := []int{}
	out, err := sh.Output("pgrep", exe[:15]) // full name does not match, only the prefix (on linux anyway)

	if err != nil || out == "" {
		return pids
	}
	for _, txt := range strings.Fields(out) {
		pid, err := strconv.Atoi(txt)
		if err == nil {
			pids = append(pids, pid)
		} else {
			fmt.Printf("Unable to format %s (%s)", txt, err)
		}
	}
	return pids
}

func killAllPIDs(pids []int) error {
	for _, pid := range pids {
		fmt.Printf("Killing process: %d\n", pid)
		err := syscall.Kill(pid, 9)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildBackend(os string, arch string, enableDebug bool) error {
	exeName := getExecutableName(os, arch)

	args := []string{
		"build", "-o", path.Join("dist", exeName), "-tags", "netgo",
	}
	if enableDebug {
		args = append(args, "-gcflags=all=-N -l")
	} else {
		args = append(args, "-ldflags", "-w")
	}
	args = append(args, "./pkg")

	env := map[string]string{
		"GOARCH": arch,
		"GOOS":   os,
	}

	// TODO: Change to sh.RunWithV once available.
	return sh.RunWith(env, "go", args...)
}

// Build is a namespace.
type Build mg.Namespace

// Linux builds the back-end plugin for Linux.
func (Build) Linux() error {
	return buildBackend("linux", "amd64", false)
}

// Windows builds the back-end plugin for Windows.
func (Build) Windows() error {
	return buildBackend("windows", "amd64", false)
}

// Darwin builds the back-end plugin for OSX.
func (Build) Darwin() error {
	return buildBackend("darwin", "amd64", false)
}

// Debug builds the debug version for the current platform
func (Build) Debug() error {
	return buildBackend(runtime.GOOS, runtime.GOARCH, true)
}

// Backend build a production build for all platforms
func (Build) Backend() {
	b := Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin)
}

// BuildAll builds production back-end components.
func BuildAll() {
	b := Build{}
	mg.Deps(b.Backend)
}

// Test runs backend tests.
func Test() error {
	if err := sh.RunV("go", "test", "./pkg/..."); err != nil {
		return nil
	}
	return nil
}

// Coverage runs backend tests and makes a coverage report.
func Coverage() error {
	// Create a coverage file if it does not already exist
	_ = os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm)

	if err := sh.RunV("go", "test", "./pkg/...", "-v", "-cover", "-coverprofile=coverage/backend.out"); err != nil {
		return nil
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.out", "-o", "coverage/backend.html"); err != nil {
		return nil
	}

	return nil
}

// Lint audits the source style
func Lint() error {
	return sh.RunV("golangci-lint", "run", "./...")
}

// Format formats the sources.
func Format() error {
	if err := sh.RunV("gofmt", "-w", "."); err != nil {
		return err
	}

	return nil
}

// Clean cleans build artifacts, by deleting the dist directory.
func Clean() error {
	err := os.RemoveAll("dist")
	if err != nil {
		return err
	}

	err = os.RemoveAll("coverage")
	if err != nil {
		return err
	}

	err = os.RemoveAll("ci")
	if err != nil {
		return err
	}
	return nil
}

func checkLinuxPtraceScope() {
	ptracePath := "/proc/sys/kernel/yama/ptrace_scope"
	byteValue, err := readFileBytes(ptracePath)
	if err != nil {
		fmt.Printf("Unable to read ptrace_scope\n")
	}
	val := strings.TrimSpace(string(byteValue[:]))
	if "0" != val {
		fmt.Printf("WARNING: \n")
		fmt.Printf("ptrace_scope set to value other than 0 (currenlty:%s), this might prevent debugger from connecting\n", val)
		fmt.Printf("try writing \"0\" to %s\n", ptracePath)
		fmt.Printf("Set ptrace_scope to 0? y/N (default N)\n")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if scanner.Text() == "y" || scanner.Text() == "Y" {
				// if err := sh.RunV("echo", "0", "|", "sudo", "tee", ptracePath); err != nil {
				// 	return // Error?
				// }
				fmt.Printf("TODO, run: echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope\n")
			} else {
				fmt.Printf("Did not write\n")
			}
		}
	}
}

// Debugger makes a new debug build and attaches dlv (go-delve).
func Debugger() error {
	// Debug build
	b := Build{}
	mg.Deps(b.Debug)

	// 1. kill any running instance
	exeName := getExecutableName(runtime.GOOS, runtime.GOARCH)

	// Kill any running processs
	_ = killAllPIDs(findRunningPIDs(exeName))
	_ = sh.RunV("pkill", "dlv")

	if runtime.GOOS == "linux" {
		checkLinuxPtraceScope()
	}

	// Wait for grafana to start plugin
	for i := 0; i < 20; i++ {
		pids := findRunningPIDs(exeName)
		if len(pids) > 1 {
			return fmt.Errorf("multiple instances already running")
		}
		if len(pids) > 0 {
			pid := strconv.Itoa(pids[0])
			log.Printf("Running PID: %s", pid)

			// dlv attach ${PLUGIN_PID} --headless --listen=:${PORT} --api-version 2 --log
			if err := sh.RunV("dlv",
				"attach",
				pid,
				"--headless",
				"--listen=:3222",
				"--api-version", "2",
				"--log"); err != nil {
				return err
			}
			fmt.Printf("dlv finished running (%s)\n", pid)
			return nil
		}

		log.Printf("waiting for grafana to start: %s...", exeName)
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("could not find process: %s, perhaps grafana is not running?", exeName)
}

// Default configures the default target.
var Default = BuildAll
