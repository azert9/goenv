package env

import "path"

func (env *Env) GetGopath() (string, error) {

	p := path.Join(env.Path, "gopath")

	return p, nil
}
