package toolchains

import (
	"os"
	"path"
)

func GetPath(version string) (string, error) {

	dirs, err := getDirs(false)
	if err != nil {
		return "", err
	}

	toolchainPath := path.Join(dirs.ToolchainsDir, version)

	if _, err := os.Lstat(toolchainPath); err != nil {
		// TODO: if the toolchain is not downloaded, return a more specific error
		return "", err
	}

	return toolchainPath, nil
}
