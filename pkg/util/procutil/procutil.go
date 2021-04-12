package procutil

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/lgdd/deba/pkg/util/fileutil"
)

func SetCatalinaPid() error {
	workspacePath, err := fileutil.GetLiferayWorkspacePath()

	if err != nil {
		return errors.Unwrap(err)
	}

	return os.Setenv("CATALINA_PID", filepath.Join(workspacePath, ".liferay-pid"))
}
