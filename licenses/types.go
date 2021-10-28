// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

// Package licenses defines a base type and support functions for describing
// software licenses.
package licenses

// For a list of licenses and comments about them, see:
// https://www.gnu.org/licenses/license-list.en.html

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// PerFileNotice is a generic per-file license statement that can be added to
// any license that does not have more specific language to recommend.
const PerFileNotice = `
Copyright (C) {{date "2006"}} {{.Author}}. All Rights Reserved.
`

// A License describes a software license.
//
// A package that implements a license should call license.Register during init
// with a value of this type, suitably populated with values corresponding to
// the details of that license.
type License struct {
	// A human-readable name of the license.
	// For example: "Apache License, Version 2.0".
	Name string

	// A slug used to identify the license. This value must be unique across all
	// registered licenses. Good slug values should be short, ideally one word,
	// with no spaces.
	Slug string

	// A URL to a description of the license (optional).
	URL string

	// The text of the license (template, required).
	Text string

	// Additional license text that must be inserted into each file covered by
	// the license (template, optional).
	PerFile string
}

// Config carries parameters to be expanded by text templates for a license.
type Config struct {
	// The name of the author, to whom copyright is attributed.
	Author string

	// The name of the project to which the license is attached, if different
	// from the author. Example: "FreeBSD".
	Project string

	// The current time. The template can render this field using the "time" and
	// "date" functions provided in the function map.
	Time time.Time
}

// newTemplate parses a text template initialized with the helpers provided by
// c, and returns a function that will execute the template into an io.Writer
// using c as its context.
func (c Config) newTemplate(text string) (func(io.Writer) error, error) {
	t, err := template.New("text").Funcs(template.FuncMap{
		"date": c.Time.Format,
		"time": c.Time.Format,
	}).Parse(text)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %v", err)
	}
	return func(w io.Writer) error {
		return t.Execute(w, c)
	}, nil
}

func cleanup(text string) *block {
	return newBlock(text).trimSpace().untabify(0).leftJust()
}

// WriteText renders the main license text to w.
func (lic *License) WriteText(w io.Writer, c *Config) error {
	if lic == nil {
		return errors.New("no license found")
	}
	clean := cleanup(lic.Text).append("") // ensure file ends with a newline
	write, err := c.newTemplate(clean.String())
	if err != nil {
		return err
	}
	return write(w)
}

// EditFile edits the per file license text into f. If the license has no
// per-file text, this does nothing without error. The indent controls how the
// text is indented or commented; if indent == nil it is inserted verbatim.
func (lic *License) EditFile(f *os.File, c *Config, indent Indenting) error {
	if lic == nil || lic.PerFile == "" {
		return nil
	}

	// Find where the file is located so we can create a tempfile in the same
	// directory.
	abs, err := filepath.Abs(f.Name())
	if err != nil {
		return err
	}

	// Seek to the beginning of the old file, so we can copy it fully.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Generate the per-file license text at the head of the file.  Ensure there
	// is a blank separating the license text from anything else below it.
	clean := indent.fix(cleanup(lic.PerFile)).append("\n")
	write, err := c.newTemplate(clean.String())
	if err != nil {
		return err
	}

	// Create a tempfile to receive the annotated file.
	tmp, err := os.CreateTemp(filepath.Dir(abs), filepath.Base(abs)+"~*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// Write the annotation to tmp, then copy the original file after it.  Sync
	// to ensure the write is committed, then close and replace the original.
	err = write(tmp)
	if err == nil {
		_, err = io.Copy(tmp, f)
		if err == nil {
			err = tmp.Sync()
		}
	}
	cerr := tmp.Close()
	if err != nil {
		return err
	} else if cerr != nil {
		return cerr
	}
	return os.Rename(tmp.Name(), f.Name())
}
