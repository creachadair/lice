# lice -- a tool for generating license files

This is a straightforward command-line tool that knows how to generate files
containing license text, and how to inject license comments into source files.

Usage examples:

    # Write a LICENSE.txt containing a BSD 3-clause license.
    lice -write LICENSE.txt -L bsd3c

    # Add comment headers to Go files for the GPLv3.
    lice -edit -L gplv3 *.go

Run `lice -help` for online help.
