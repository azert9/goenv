package toolchains

import (
	"errors"
	"os"
	"path"
)

func getRef(dirs *dirs, ref string) (*string, error) {

	link, err := os.Readlink(path.Join(dirs.ToolchainsDir, ref))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &link, nil
}

func updateRefs(dirs *dirs) error {

	// TODO: this will not behave properly when having toolchain versions that cannot be parsed.

	versions, err := List()
	if err != nil {
		return err
	}

	// finding the last version that is not "latest"
	var latestVersion *string
	for i := len(versions) - 1; i >= 0; i-- {
		if versions[i] != "latest" {
			latestVersion = &versions[i]
			break
		}
	}

	latestSymlinkPath := path.Join(dirs.ToolchainsDir, "latest")

	curLatestVersion, err := getRef(dirs, "latest")
	if err != nil {
		return err
	}

	if (latestVersion == nil && curLatestVersion == nil) || (latestVersion != nil && curLatestVersion != nil && *latestVersion == *curLatestVersion) {

		// nothing to do

	} else {

		if err := os.Remove(latestSymlinkPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if latestVersion != nil {
			if err := os.Symlink(*latestVersion, latestSymlinkPath); err != nil {
				return err
			}
		}
	}

	return nil
}
