package procutil

import (
	"errors"
	"fmt"
	"os"

	"github.com/lgdd/deba/pkg/util/fileutil"
)

func SetCatalinaPid() error {
	pathSeparator := string(os.PathSeparator)
	workspacePath, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		return errors.Unwrap(err)
	}

	return os.Setenv("CATALINA_PID", fmt.Sprintf("%s%s%s", workspacePath, pathSeparator, ".liferay-pid"))
}
