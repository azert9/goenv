package main

import (
	_ "embed"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/azert9/goenv/internal/env"
	"github.com/azert9/goenv/internal/toolchains"
	"github.com/azert9/goenv/internal/utils"
	"os"
	"path"
)

// TODO: use a config file
const defaultGoenvDir = ".goenv"

type Options struct {
	EnvDir           string
	ToolchainVersion string
	Subcommand       func(options *Options) error
}

func SubcommandInit(options *Options) error {

	// TODO: allow choosing the go version, etc...

	goenv := env.Env{
		Path: options.EnvDir,
	}

	if err := goenv.Init(); err != nil {
		return err
	}

	fmt.Printf("Activate with: . %s\n", utils.ShellEscape(path.Join(options.EnvDir, "activate.sh")))

	return nil
}

func SubcommandToolchainList(options *Options) error {

	results, err := toolchains.List()
	if err != nil {
		return err
	}

	if len(results) > 0 {

		fmt.Printf("Available toolchains:\n")

		for _, result := range results {
			fmt.Printf("* %s\n", result)
		}
	} else {

		fmt.Printf("No toolchain available.\n")
	}

	return nil
}

func SubcommandToolchainDownload(options *Options) error {
	return toolchains.Download(options.ToolchainVersion)
}

func SubcommandToolchainUse(options *Options) error {

	goenv := env.Env{
		Path: options.EnvDir,
	}

	goroot, err := toolchains.GetPath(options.ToolchainVersion)
	if err != nil {
		return err
	}

	if err := goenv.SetGoroot(goroot); err != nil {
		return err
	}

	return nil
}

func SubcommandToolchainRemove(options *Options) error {

	if err := toolchains.Remove(options.ToolchainVersion); err != nil {
		return err
	}

	fmt.Printf("Removed toolchain %s. Be careful to update any goenv that could use it.")

	return nil
}

func SubcommandShow(options *Options) error {

	goenv := env.Env{
		Path: options.EnvDir,
	}

	fmt.Printf("Location: %s\n", goenv.Path)

	// TODO: show gopath?

	goroot, err := goenv.GetGoroot()
	if err != nil {
		return err
	}
	fmt.Printf("GOROOT: %s\n", goroot)

	gopath, err := goenv.GetGopath()
	if err != nil {
		return err
	}
	fmt.Printf("GOPATH: %s\n", gopath)

	return nil
}

func parseProgramArgs(args []string, out *Options) error {

	parser := argparse.NewParser(args[0], "") // TODO

	dirOpt := parser.String("d", "dir", &argparse.Options{
		Required: false,
		Help:     "Path to the environment directory.",
	})

	initSubcommand := parser.NewCommand("init", "Initialize a new environment.")
	showSubcommand := parser.NewCommand("show", "Show information about the environment.")
	toolchainSubcommand := parser.NewCommand("toolchain", "Manage Go toolchains.")
	toolchainListSubcommand := toolchainSubcommand.NewCommand("list", "List downloaded toolchains.")
	toolchainUseSubcommand := toolchainSubcommand.NewCommand("use", "Configure the environment to use the specified toolchain.")
	toolchainDownloadSubcommand := toolchainSubcommand.NewCommand("download", "Download a Go toolchain.")
	toolchainRemoveSubcommand := toolchainSubcommand.NewCommand("remove", "Remove a downloaded Go toolchain.")

	toolchainVersionOpt := toolchainSubcommand.StringPositional(&argparse.Options{
		Help: "Toolchain version to use.",
	})

	if err := parser.Parse(args); err != nil {
		return err
	}

	out.EnvDir = *dirOpt
	out.ToolchainVersion = *toolchainVersionOpt

	toolchainVersionRequired := false

	switch {
	case initSubcommand.Happened():
		out.Subcommand = SubcommandInit
	case showSubcommand.Happened():
		out.Subcommand = SubcommandShow
	case toolchainListSubcommand.Happened():
		out.Subcommand = SubcommandToolchainList
	case toolchainUseSubcommand.Happened():
		out.Subcommand = SubcommandToolchainUse
		toolchainVersionRequired = true
	case toolchainDownloadSubcommand.Happened():
		out.Subcommand = SubcommandToolchainDownload
		toolchainVersionRequired = true
	case toolchainRemoveSubcommand.Happened():
		out.Subcommand = SubcommandToolchainRemove
		toolchainVersionRequired = true
	default:
		panic("missing case for subcommand")
	}

	if toolchainVersionRequired && *toolchainVersionOpt == "" {
		return fmt.Errorf("missing positional argument")
	}

	return nil
}

func main() {

	// TODO: always check that toolchain versions do not contain slashes

	var options Options

	if err := parseProgramArgs(os.Args, &options); err != nil {
		// TODO
		panic(err)
	}

	if options.EnvDir == "" {
		value, found := os.LookupEnv("GOENV_PATH")
		if found {
			options.EnvDir = value
		}
	}

	if options.EnvDir == "" {
		options.EnvDir = defaultGoenvDir
	}

	if err := options.Subcommand(&options); err != nil {
		// TODO
		panic(err)
	}
}
