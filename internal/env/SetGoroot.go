package env

import (
	"os"
	"path"
)

func (env *Env) SetGoroot(target string) error {

	p := path.Join(env.Path, "goroot")

	if err := os.Remove(p); err != nil {
		return err
	}

	if err := os.Symlink(target, p); err != nil {
		return err
	}

	return nil
}
