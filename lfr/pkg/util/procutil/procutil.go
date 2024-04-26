package procutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ps "github.com/mitchellh/go-ps"

	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
)

// Set the catalina pid as an environment variable
func SetCatalinaPid() error {
	workingPath, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		workingPath, err = fileutil.GetLiferayHomePath()
		if err != nil {
			return errors.Unwrap(err)
		}
	}

	return os.Setenv("CATALINA_PID", filepath.Join(workingPath, ".liferay-pid"))
}

// Get the catalina pid from an environment variable or file
func GetCatalinaPid() (int, error) {
	workingPath, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		workingPath, err = fileutil.GetLiferayHomePath()
		if err != nil {
			return 0, err
		}
	}

	pidString := os.Getenv("CATALINA_PID")

	if pidString != "" {
		pid, err := strconv.Atoi(pidString)

		if err != nil {
			return 0, err
		}

		return pid, nil
	}

	pidPath := filepath.Join(workingPath, ".liferay-pid")
	pidFile, err := os.Open(pidPath)

	if err != nil {
		return 0, err
	}

	defer pidFile.Close()

	pidBytes, err := io.ReadAll(pidFile)

	if err != nil {
		return 0, err
	}

	pidString = string(pidBytes)
	pidString = strings.ReplaceAll(pidString, "\n", "")
	pid, err := strconv.Atoi(pidString)

	if err != nil {
		return 0, err
	}

	if pid == 0 {
		return 0, errors.New("pid not found")
	}

	return pid, nil
}

// Checks if the Liferay bundle is running by checking its pid
func IsCatalinaRunning() (bool, int, error) {
	pid, err := GetCatalinaPid()

	if err != nil {
		return false, 0, errors.New("couldn't find the PID of a Liferay Tomcat bundle running")
	}

	proc, err := ps.FindProcess(pid)

	if err != nil {
		return false, 0, nil
	}

	defer func() {
		if err := recover(); err != nil {
			// avoiding next method calls to throw panics
			// to let the caller handle properly tuple returned
			fmt.Print("")
		}
	}()

	ppid := proc.PPid()

	_, err = os.FindProcess(ppid)

	if err != nil {
		return false, 0, err
	}

	return true, pid, nil
}

func GetCurrentJavaVersion() (string, string, error) {
	_, javaVersionCmdErr, err := Exec("java", "-version")

	if err != nil {
		return "", "", err
	}

	versionRegex := regexp.MustCompile(`\d+.\d+.\d+`)
	javaVersionResult := javaVersionCmdErr.String()

	fullVersion := versionRegex.FindString(javaVersionResult)
	majorVersion := getMajorJavaVersion(fullVersion)

	return majorVersion, fullVersion, nil
}

func getMajorJavaVersion(fullVersion string) string {
	versionNumbers := strings.Split(fullVersion, ".")

	if versionNumbers[0] == "1" {
		return versionNumbers[1]
	}

	return versionNumbers[0]
}

func Exec(command string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	var cmdOut, cmdErr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()

	return cmdOut, cmdErr, err
}

func ExecStd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}
