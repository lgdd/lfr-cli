package procutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"

	"github.com/lgdd/liferay-cli/lfr/pkg/util/fileutil"
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

	pidBytes, err := ioutil.ReadAll(pidFile)

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
