package toolchains

import (
	"os"
)

func Remove(version string) error {

	p, err := GetPath(version)
	if err != nil {
		// TODO: return nil if the toolchain was not downloaded?
		return err
	}

	if err := os.RemoveAll(p); err != nil {
		return err
	}

	return nil
}
