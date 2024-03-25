package main

import (
	"io/fs"
	"os"
	"regexp"
	"strconv"
)

type Regexp struct {
	regexp.Regexp
}

func CompileRegex(s string) (*Regexp, error) {
	logger.Debugf("Compiling regex (%q)", s)

	regex, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}

	return &Regexp{*regex}, nil
}

func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := CompileRegex(string(b))
	if err != nil {
		return err
	}

	*r = *regex

	return nil
}

func (r *Regexp) MarshalText() ([]byte, error) {
	return []byte(r.Regexp.String()), nil
}

type FileMode struct {
	fs.FileMode
}

func (m *FileMode) UnmarshalText(b []byte) error {
	logger.Debugf("Parsing file mode (%q)", string(b))

	mode, err := strconv.ParseUint(string(b), 8, 32)
	m.FileMode = os.FileMode(mode)
	return err
}
