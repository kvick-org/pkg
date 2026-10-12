package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

var pseudo = regexp.MustCompile(`^v[0-9]+\.(0\.0-|\d+\.\d+-([^+]*\.)?0\.)\d{14}-[A-Za-z0-9]+(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

const (
	Devel = "devel"
)

type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

type Build struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
}

type Runtime struct {
	Distro string `json:"distro"`
	OS     string `json:"os"`
	Arch   string `json:"arch"`
}

type Info struct {
	Build   Build   `json:"build" text:"Build"`
	Runtime Runtime `json:"runtime" text:"Runtime"`
}

func (i Info) Output(w io.Writer, format OutputFormat) error {
	switch format {
	case OutputFormatText:
		fmt.Fprintf(w, "%s version %s %s\n", filepath.Base(os.Args[0]), i.Build.Version, i.Build.Commit)
		return nil
	case OutputFormatJSON:
		b, err := json.Marshal(&i)
		if err != nil {
			return err
		}
		fmt.Fprint(w, string(b))
		return nil
	default:
		return fmt.Errorf("unknown output format %s", format)
	}
}

func Load() (Info, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return Info{}, errors.New("could not read build info")
	}
	info := Info{
		Build: Build{
			GoVersion: bi.GoVersion,
			Version:   Devel,
		},
		Runtime: Runtime{
			Distro: getDistro(),
		},
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "GOARCH":
			info.Runtime.Arch = s.Value
		case "GOOS":
			info.Runtime.OS = s.Value
		case "vcs.revision":
			info.Build.Commit = s.Value
		case "vcs.modified":
			modified, err := strconv.ParseBool(s.Value)
			if err != nil {
				return Info{}, err
			}
			if bi.Main.Version == "" {
				continue
			}
			if modified {
				continue
			}
			if pseudo.MatchString(bi.Main.Version) {
				continue
			}
			info.Build.Version = bi.Main.Version
		}
	}
	return info, nil
}

func getDistro() string {
	unknownDistro := "unknown"
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return unknownDistro
	}
	for line := range strings.SplitSeq(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], `"`)
		}
	}
	return unknownDistro
}
