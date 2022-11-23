package env

import (
	"os"
	"path"
)

func (env *Env) GetGoroot() (string, error) {

	p := path.Join(env.Path, "goroot")

	target, err := os.Readlink(p)
	if err != nil {
		return "", err
	}

	return target, nil
}
