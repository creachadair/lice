// Copyright (C) 2018, Michael J. Fromberger
// All Rights Reserved.

// Program lice creates license files and injects license text into source.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"text/tabwriter"
	"time"

	"bitbucket.org/creachadair/goflags/enumflag"
	"bitbucket.org/creachadair/goflags/timeflag"
	"bitbucket.org/creachadair/lice/licenses"

	_ "bitbucket.org/creachadair/lice/licenses/apache"
	_ "bitbucket.org/creachadair/lice/licenses/bsd"
	_ "bitbucket.org/creachadair/lice/licenses/cc"
	_ "bitbucket.org/creachadair/lice/licenses/gpl"
	_ "bitbucket.org/creachadair/lice/licenses/mit"
)

var (
	indentStyle = enumflag.New("guess", "hash", "none", "slash", "star", "sstar", "xml")
	dateNow     = &timeflag.Value{Layout: "2006-01-02", Time: time.Now()}
	writeFile   = flag.String("write", "", "Write a license file at this path")
	slug        = flag.String("L", "", "License to use (use -list for a list)")
	doForce     = flag.Bool("f", false, "Force overwrite of existing files")
	doEdit      = flag.Bool("edit", false, "Edit license text into non-flag argument files")
	doList      = flag.Bool("list", false, "List available licenses")
	viewLicense = flag.String("view", "", "View license text")

	userName string

	indent = map[string]licenses.Indenting{
		"hash":  licenses.IPrefix("# "),                    // like bash, Python, Perl
		"slash": licenses.IPrefix("// "),                   // like C++, Go, Java
		"star":  licenses.IComment("/*", "   ", " */"),     // like C
		"sstar": licenses.IComment("/*", " * ", " */"),     // like C
		"xml":   licenses.IComment("<!--", "   ", "  -->"), // like HTML, XML
	}
)

func init() {
	flag.Var(indentStyle, "i", indentStyle.Help("Indentation style"))
	flag.Var(dateNow, "date", dateNow.Help("Copyright date for attribution"))

	u, err := user.Current()
	if err != nil {
		log.Panicf("Unable to determine current user: %v", err)
	}
	flag.StringVar(&userName, "author", u.Name, "Copyright author for attribution")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage: %[1]s [-list | -view <license>]
       %[1]s -L <license> -write <file>
       %[1]s -L <license> -edit <file1> <file2> ...

Generate license text for source code. With -list, the available license types
are listed. With -write, the tool writes the text of a license to the specified
file, substituting in the -author and -date information as necessary.

If -edit is set, any additional files named on the command line are edited in
place to insert a comment containing a per-file license annotation, if the
selected license type has one.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// If a list is requested, do that and exit early.
	if *doList {
		if *doEdit || *viewLicense != "" || *writeFile != "" {
			log.Fatal("You may not combine -write, -edit, or -view with -list")
		}
		fmt.Println("Available licenses:")
		tw := tabwriter.NewWriter(os.Stdout, 8, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
		licenses.List(func(lic licenses.License) {
			fmt.Fprint(tw, lic.Slug, "\t", lic.Name, "\t", lic.URL, "\n")
		})
		tw.Flush()
		return
	} else if *viewLicense != "" {
		*slug = *viewLicense
	} else if *slug == "" && *viewLicense == "" {
		log.Fatal("You must specify a license to use with -L")
	}

	lic := licenses.Lookup(*slug)
	if lic == nil {
		log.Fatalf("Unknown license type %q (use -list for a list)", *slug)
	}

	cfg := &licenses.Config{
		Author: userName,
		Time:   dateNow.Time,
	}

	// View a license.
	if *viewLicense != "" {
		if err := lic.WriteText(os.Stdout, cfg); err != nil {
			log.Fatalf("Viewing license: %v", err)
		}
	}

	// Write a license to a file.
	if *writeFile != "" {
		oflag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		if !*doForce {
			oflag |= os.O_EXCL
		}
		f, err := os.OpenFile(*writeFile, oflag, 0644)
		if err != nil {
			log.Fatalf("Opening license file: %v", err)
		}
		err = lic.WriteText(f, cfg)
		cerr := f.Close()
		if err != nil {
			log.Fatalf("Writing license file: %v", err)
		} else if cerr != nil {
			log.Fatalf("Closing license file: %v", cerr)
		}
		fmt.Fprintf(os.Stderr, "Wrote %s to %s\n", lic.Name, *writeFile)
	}

	// Edit license tags into other files, if available.
	if !*doEdit || flag.NArg() == 0 || lic.PerFile == "" {
		return
	}
	hasErr := false
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Opening file: %v [skipped]", err)
			hasErr = true
			continue
		}
		func() {
			defer f.Close()
			if err := lic.EditFile(f, cfg, chooseIndent(path)); err != nil {
				log.Printf("Editing file: %v", err)
				hasErr = true
			} else {
				fmt.Fprintf(os.Stderr, "Added %s to %s\n", lic.Name, path)
			}
		}()
	}

	if hasErr {
		os.Exit(1)
	}
}

// chooseIndent picks a suitable indenting rule for a file. If an indenting
// rule was specified by the user, use that; otherwise if the user asked us to
// guess, do so based on its file extension. If no indenting rule can be
// inferred, fall back to undecorated text.
func chooseIndent(path string) licenses.Indenting {
	in, ok := indent[indentStyle.Key()]
	if ok {
		return in
	} else if indentStyle.Key() != "guess" {
		return nil
	}
	switch filepath.Ext(path) {
	case "", ".sh", ".py", ".pl", ".rb": // N.B. includes no extension
		return indent["hash"]
	case ".cc", ".cpp", ".go", ".java", ".js":
		return indent["slash"]
	case ".c", ".h":
		return indent["star"]
	case ".htm", ".html", ".xhtml":
		return indent["xml"]
	default:
		return nil
	}
}
