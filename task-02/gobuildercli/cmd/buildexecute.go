/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	buildDir    string
	copyDir     string
	exeName     string
	exclueTests bool
)

func handleCopy(dest, src string) error {
	matches, err := filepath.Glob(src)
	if err != nil {
		return err
	}
	for _, match := range matches {
		fiSrc, err := os.Stat(match)
		if err != nil {
			return nil
		}
		destPath := filepath.Join(dest, fiSrc.Name())
		fiDest, _ := os.Stat(destPath)
		// If both the files are same, continue
		if os.SameFile(fiSrc, fiDest) {
			continue
		}
		if fiSrc.IsDir() {
			CopyDir(match, destPath)
		} else {
			os.MkdirAll(dest, os.ModePerm)
			CopyFile(match, destPath)
		}
	}
	return nil
}

// buildexecuteCmd represents the buildexecute command
var buildexecuteCmd = &cobra.Command{
	Use:   "buildexecute",
	Short: "buildexecute means the command has been executed",
	RunE: func(cmd *cobra.Command, args []string) error {
		// handle copy
		if copyDir != "" {
			err := handleCopy(buildDir, copyDir)
			if err != nil {
				return nil
			}
		}

		if exeName != "" {
			shArgs := []string{"build"}
			if exeName != "" {
				shArgs = append(shArgs, "-o")
				shArgs = append(shArgs, exeName)
			}
			shCmd := exec.Command("go", shArgs...)
			if err := shCmd.Start(); err != nil {
				return err
			}

			if err := shCmd.Wait(); err != nil {
				return err
			}
		}
		return nil

	},
}

func init() {
	// TODO: how to do single dash?
	buildexecuteCmd.Flags().StringVarP(&buildDir, "builddir", "d", ".", "Copy destination")
	buildexecuteCmd.Flags().StringVarP(&copyDir, "copydir", "c", "", "Copy source")
	buildexecuteCmd.Flags().StringVarP(&exeName, "exe", "o", "", "Executable name")
	buildexecuteCmd.Flags().BoolVarP(&exclueTests, "exclude-tests", "e", false, "Exclude test files")
	rootCmd.AddCommand(buildexecuteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildexecuteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

// SRC:  https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	if exclueTests {
		if strings.HasSuffix(src, "_test.go") {
			return
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
