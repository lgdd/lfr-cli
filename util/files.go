package util

import (
	"bytes"
	"github.com/markbates/pkger"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

func CopyPkgedFile(sourcePath, destPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	source, err := pkger.Open(sourcePath)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}
}

func UpdateFile(pomPath string, metadata *Project) error {
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
