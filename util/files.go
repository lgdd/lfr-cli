package util

import (
	"bytes"
	"github.com/markbates/pkger"
	"html/template"
	"io"
	"io/ioutil"
	"os"
)

func CopyPkgedFile(sourcePath, destPath string) error {
	source, err := pkger.Open(sourcePath)
	if err != nil {
		return err
	}

	defer source.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
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
