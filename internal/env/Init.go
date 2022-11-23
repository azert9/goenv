package env

import (
	_ "embed"
	"github.com/azert9/goenv/internal/toolchains"
	"github.com/azert9/goenv/internal/utils"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

//go:embed activate.in.sh
var activateScriptTemplate string

// TODO: make private if possible
type ShellString string

func (str ShellString) Escape() ShellString {
	return ShellString(utils.ShellEscape(string(str)))
}

// TODO: make private if possible
type TemplateParams struct {
	Dir    ShellString
	GoRoot ShellString
	GoPath ShellString
}

func RenderTemplate(tmpl *template.Template, targetPath string, params *TemplateParams) error {

	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			// TODO: just issue a warning
			panic(err)
		}
	}(file)

	if err := tmpl.Execute(file, params); err != nil {
		return err
	}

	return nil
}

func (env *Env) Init() error {

	// TODO: remove the whole directory on error

	envAbsPath, err := filepath.Abs(env.Path)
	if err != nil {
		return err
	}

	gopathAbsPath := path.Join(envAbsPath, "gopath")
	gorootAbsPath := path.Join(envAbsPath, "goroot")

	if err := os.Mkdir(envAbsPath, 0777); err != nil {
		return err
	}

	if err := os.Mkdir(gopathAbsPath, 0777); err != nil {
		return err
	}

	// TODO: handle situations where no toolchain is available

	toolchainPath, err := toolchains.GetPath("latest")
	if err != nil {
		return err
	}

	if err := os.Symlink(toolchainPath, gorootAbsPath); err != nil {
		return err
	}

	tmpl := template.Must(template.New("").Parse(activateScriptTemplate))

	params := TemplateParams{
		Dir:    ShellString(envAbsPath),
		GoRoot: ShellString(gorootAbsPath),
		GoPath: ShellString(gopathAbsPath),
	}

	if err := RenderTemplate(tmpl, path.Join(envAbsPath, "activate.sh"), &params); err != nil {
		return err
	}

	return nil
}
