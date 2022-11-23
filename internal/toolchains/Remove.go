package toolchains

import (
	"errors"
	"os"
)

func Remove(version string) error {

	dirs, err := getDirs(false)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return err
		}
	}

	p, err := GetPath(version)
	if err != nil {
		// TODO: return nil if the toolchain was not downloaded?
		return err
	}

	if err := os.RemoveAll(p); err != nil {
		return err
	}

	if err := updateRefs(dirs); err != nil {
		return err
	}

	return nil
}
