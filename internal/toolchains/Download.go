package toolchains

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func httpDownload(url string, out io.Writer) error {

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		// TODO: warn on error
		response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status code: %s", response.Status)
	}

	if _, err := io.Copy(out, response.Body); err != nil {
		return err
	}

	return nil
}

func httpDownloadToPath(url string, destPath string) error {

	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		// TODO: warn on error
		file.Close()
	}()

	return httpDownload(url, file)
}

func extractArchive(archivePath string, destPath string) error {

	cmd := exec.Command("tar", "-xf", archivePath)
	cmd.Dir = destPath

	return cmd.Run()
}

func Download(version string) error {

	dirs, err := getDirs(true)
	if err != nil {
		return err
	}

	extractPath := path.Join(dirs.ExtractDir, version)

	destPath := path.Join(dirs.ToolchainsDir, version)

	if _, err = os.Lstat(destPath); err == nil {
		fmt.Printf("Already downloaded.\n")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	archiveExtension := ".tar.gz"

	url := fmt.Sprintf("https://go.dev/dl/%s.%s-%s%s", version, runtime.GOOS, runtime.GOARCH, archiveExtension)

	fmt.Printf("Downloading %s...\n", url)

	// TODO: remove the directory when done
	tmpdir, err := os.MkdirTemp("", "goenv")
	if err != nil {
		return err
	}
	defer func() {
		// TODO: warn on error
		os.RemoveAll(tmpdir)
	}()

	archivePath := path.Join(tmpdir, fmt.Sprintf("archive%s", archiveExtension))

	if err := httpDownloadToPath(url, archivePath); err != nil {
		return err
	}

	fmt.Printf("Extracting...\n")

	// TODO: handle existing extractPath (failed or concurrent download)

	if err := os.MkdirAll(extractPath, 0777); err != nil {
		return err
	}
	defer func() {
		// TODO: warn on error
		os.RemoveAll(extractPath)
	}()

	if err := extractArchive(archivePath, extractPath); err != nil {
		return err
	}

	if err := os.Rename(path.Join(extractPath, "go"), destPath); err != nil {
		return err
	}

	if err := updateRefs(dirs); err != nil {
		return err
	}

	fmt.Printf("Done!\n")

	return nil
}
