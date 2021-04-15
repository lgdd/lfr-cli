package procutil

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lgdd/deba/pkg/util/fileutil"
)

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

func GetCatalinaPid() (int, error) {
	workingPath, err := fileutil.GetLiferayHomePath()

	if err != nil {
		workingPath, err = fileutil.GetLiferayWorkspacePath()
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
