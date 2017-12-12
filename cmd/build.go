// Copyright © 2017 Igor Franca <lee12rock@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
)

// Process Struct defines the behavior from the given process
// To be used as part of a analysis system.
type Process struct {
	status    int
	lang      string
	originURL string
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build <args>",
	Short: "Runs the build command, followed by your tipical 'run' like command :D, e.g. node server.js",
	Long:  `'Build' command will analyze your application to support some 'insiders' info, like the Git origin URL`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("help").Value.String() == "" {
			printUsage()
		}
		Build(args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("help", "", "A help flag for the command")
	buildCmd.PersistentFlags().Bool("help", false, "A help flag for the command")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printUsage() {
	fmt.Printf(`
-> build <args> [flags]: Runs the main command! e.g node server.js
Samples: 
	- morfeu build node server.js`)
}

// Build ... Function that is the main functionality,
// just running the build command from the contestant application
func Build(args []string) {
	var out Process
	var control sync.WaitGroup

	fmt.Printf("Running... %v\n", args[0:])

	control.Add(1)

	//See doc above the function def :D
	go getLangVersion(&out, &control, args)

	control.Add(1)

	//See doc above the function def :D
	go getExitCode(&out, &control, args)

	control.Add(1)

	//See doc above the function def :D
	go getGitOrigin(&out, &control)

	control.Wait()

	fmt.Printf("OUTPUT => %v\n", out)
}

// getLangVersion is used to get from a map of supported languages the command that will output
// the language version
// @TODO implements support to frameworks (like Rails, Ionic, etc...)
func getLangVersion(out *Process, control *sync.WaitGroup, args []string) {
	defer control.Done()

	version, err := exec.Command(langVersion(args[0])).CombinedOutput()

	validateError(err)

	out.lang = cleanReturn(os.Args[0], version)
}

// getExitCode runs the command passed to the system and determines if its return a valid output
func getExitCode(out *Process, control *sync.WaitGroup, args []string) {
	defer control.Done()

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Start()
	cmd.Wait()

	exitCode, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)

	out.status = exitCode.ExitStatus()
}

// getGitOrigin will catch up the origin URL of the Git project and send to Shawee API :D
func getGitOrigin(out *Process, control *sync.WaitGroup) {
	defer control.Done()

	git, err := exec.LookPath("/usr/bin/git")

	validateError(err)

	args := []string{"remote", "-v"}

	origin, err := exec.Command(git, args...).CombinedOutput()

	validateError(err)

	out.originURL = cleanOriginString(string(origin))
}

// validateError is a utility function to dont repeat code.
func validateError(err error) {
	if err != nil {
		fmt.Printf("Ocorreu um erro! => %v\n", err)
		os.Exit(0)
	}
	return
}

// cleanReturn is a utility function that cleans and converts the output from the process
// to the struct that is storing it
func cleanReturn(args string, version []byte) string {
	stringVersion := string(version)

	args = strings.Title(args)

	return args + " " + stringVersion
}

func cleanOriginString(origin string) string {
	inter := strings.SplitAfter(origin, ".git")[0]

	inter = strings.Replace(inter, "origin", "", -1)

	inter = strings.Replace(inter, ".git", "", -1)

	return inter
}
