package fileutil

import (
	"bytes"
	"github.com/lgdd/deba/pkg/project"
	"github.com/lgdd/deba/pkg/util/printutil"
	"github.com/markbates/pkger"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
)

func CopyFromAssets(sourcePath, destPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	source, err := pkger.Open(sourcePath)
	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		printutil.Error(err.Error())
		os.Exit(1)
	}
}

func UpdateWithData(pomPath string, metadata *project.Metadata) error {
	pomContent, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return err
	}

	tpl, err := template.New(pomPath).Parse(string(pomContent))
	if err != nil {
		return err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, metadata)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pomPath, result.Bytes(), 0664)
	if err != nil {
		return err
	}

	return nil
}
