/*****************************************************************************/
/* create.go                                                                 */
/*****************************************************************************/
/*                           This file is part of:                           */
/*                                KAIJU ENGINE                               */
/*                          https://kaijuengine.org                          */
/*****************************************************************************/
/* MIT License                                                               */
/*                                                                           */
/* Copyright (c) 2023-present Kaiju Engine contributors (CONTRIBUTORS.md).   */
/* Copyright (c) 2015-2023 Brent Farris.                                     */
/*                                                                           */
/* May all those that this source may reach be blessed by the LORD and find  */
/* peace and joy in life.                                                    */
/* Everyone who drinks of this water will be thirsty again; but whoever      */
/* drinks of the water that I will give him shall never thirst; John 4:13-14 */
/*                                                                           */
/* Permission is hereby granted, free of charge, to any person obtaining a   */
/* copy of this software and associated documentation files (the "Software"),*/
/* to deal in the Software without restriction, including without limitation */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,  */
/* and/or sell copies of the Software, and to permit persons to whom the     */
/* Software is furnished to do so, subject to the following conditions:      */
/*                                                                           */
/* The above copyright, blessing, biblical verse, notice and                 */
/* this permission notice shall be included in all copies or                 */
/* substantial portions of the Software.                                     */
/*                                                                           */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS   */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.    */
/* IN NO EVENT SHALL THE /* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY   */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE     */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                             */
/*****************************************************************************/

package project

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	projectTemplateFolder = "project_template"
	sourceFolder          = "src/source"
)

func CreateNew(path string) error {
	if stat, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return os.ErrExist
	}
	if err := unzipTemplate(path); err != nil {
		return err
	}
	if err := createCache(path); err != nil {
		return err
	}
	return createSource(path)
}

func unzipTemplate(into string) error {
	r, err := zip.OpenReader("project_template.zip")
	if err != nil {
		return err
	}
	defer r.Close()
	for _, file := range r.File {
		sf, err := file.Open()
		if err != nil {
			return err
		}
		defer sf.Close()
		toPath := filepath.Join(into, file.Name)
		os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
		df, err := os.Create(toPath)
		if err != nil {
			return err
		}
		defer df.Close()
		if _, err = io.Copy(df, sf); err != nil {
			return err
		}
	}
	return nil
}

func createCache(projectFolder string) error {
	if _, err := os.Stat(filepath.Join(projectFolder, "/.cache")); err != nil {
		return os.Mkdir(filepath.Join(projectFolder, "/.cache"), 0755)
	}
	return nil
}

func createSource(projectFolder string) error {
	sourceDir := filepath.Join(projectFolder, sourceFolder)
	err := os.Mkdir(sourceDir, 0755)
	if err != nil {
		return err
	}
	mainFile := filepath.Join(sourceDir, "/source.go")
	_, err = os.Stat(mainFile)
	if err == nil {
		return errors.New("source file already exists and should not")
	}
	f, err := os.Create(mainFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(`package source

import "kaiju/engine"

func Main(host *engine.Host) {
	
}
`)
	return err
}
