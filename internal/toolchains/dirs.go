package toolchains

import (
	"os"
	"path"
)

type dirs struct {
	ToolchainsDir string
	ExtractDir    string
}

func getDirs(create bool) (*dirs, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := path.Join(homeDir, ".local/share/goenv")

	dirs := &dirs{
		ToolchainsDir: path.Join(dataDir, "toolchains"),
		ExtractDir:    path.Join(dataDir, "extract"),
	}

	if create {

		if err := os.MkdirAll(dirs.ToolchainsDir, 0777); err != nil {
			return nil, err
		}

		if err := os.MkdirAll(dirs.ExtractDir, 0777); err != nil {
			return nil, err
		}
	} else {

		if _, err := os.Lstat(dirs.ToolchainsDir); err != nil {
			return nil, err
		}

		if _, err := os.Lstat(dirs.ExtractDir); err != nil {
			return nil, err
		}
	}

	return dirs, nil
}
