// Copyright (C) 2018 Michael J. Fromberger. All Rights Reserved.

// Package bsd describes BSD software licenses.
package bsd

import "github.com/creachadair/lice/licenses"

func init() {
	licenses.Register(licenses.License{
		Name:    "Modified BSD license (3-clause)",
		Slug:    "bsd3c",
		URL:     "https://directory.fsf.org/wiki/License:BSD-3-Clause",
		Text:    bsd3text,
		PerFile: licenses.PerFileNotice,
	})
}

const bsd3text = `
BSD 3-Clause License

Copyright (C) {{date "2006"}}, {{.Author}}
All Rights Reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    (1) Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

    (2) Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.

    (3) The name of the author may not be used to endorse or promote products
    derived from this software without specific prior written permission.
` + disclaimer

const disclaimer = `
THIS SOFTWARE IS PROVIDED BY THE AUTHOR "AS IS" AND ANY EXPRESS OR IMPLIED
WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO
EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT
OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY
OF SUCH DAMAGE.
`
