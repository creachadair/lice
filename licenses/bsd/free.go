// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

package bsd

import "github.com/creachadair/lice/licenses"

func init() {
	licenses.Register(licenses.License{
		Name:    "FreeBSD software license",
		Slug:    "freebsd",
		URL:     "https://www.freebsd.org/copyright/freebsd-license.html",
		Text:    freetext,
		PerFile: licenses.PerFileNotice,
	})
}

const freetext = `
Copyright {{date "2006"}}, {{.Author}}
All Rights Reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
` + disclaimer + `{{if .Project}}
The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the {{.Project}} Project.
{{end}}`
