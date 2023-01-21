package twshim

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
)

const baseReleaseURL = "https://api.github.com/repos/tailwindlabs/tailwindcss/releases" // without trailing slash

// Command returns an exec.Cmd ready for execution.
// If the binary of the given release does not exist in downloadRoot, it is downloaded first.
func Command(downloadRoot, releaseTag, assetName string, arg ...string) (*exec.Cmd, error) {
	dir := path.Join(downloadRoot, releaseTag)
	bin := path.Join(dir, assetName)
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
		if err := DownloadCLI(releaseTag, assetName, bin); err != nil {
			return nil, fmt.Errorf("downloading tailwing %s: %w", releaseTag, err)
		}
	} else if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd, nil
}

// RuntimeAssetName derives the name of the GitHub asset that holds the tailwindcss standalone CLI binary for
// the operating system and architecture this func is called on.
func RuntimeAssetName() (string, error) {
	dist := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	// Add more distributions as needed.
	switch dist {
	case "linux/amd64":
		return "tailwindcss-linux-x64", nil
	case "windows/amd64":
		return "tailwindcss-windows-x64", nil
	default:
		return "", fmt.Errorf("distribution not supported: %s", dist)
	}
}

// DownloadCLI downloads a named asset from a release by tag to the desired destination and makes it executable.
func DownloadCLI(tag, asset, dest string) error {
	rURL := baseReleaseURL + "/tags/" + tag
	r, err := fetchRelease(rURL)
	if err != nil {
		return err
	}

	a, err := r.findAsset(asset)
	if err != nil {
		return err
	}
	aURL := baseReleaseURL + "/assets/" + strconv.Itoa(a.ID)
	if err := downloadFile(aURL, dest); err != nil {
		return err
	}

	// Make it executable by user and group.
	err = os.Chmod(dest, 0770)
	if err != nil {
		return err
	}

	return nil
}

func fetchRelease(url string) (githubRelease, error) {
	resp, err := http.Get(url)
	if err != nil {
		return githubRelease{}, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("bad status: %s", resp.Status)
	}
	var r githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return githubRelease{}, err
	}
	return r, nil
}

type githubRelease struct {
	Name   string        `json:"name"`
	Assets []githubAsset `json:"assets"`
	Tag    string        `json:"tag_name"`
}

func (r *githubRelease) findAsset(name string) (githubAsset, error) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, nil
		}
	}
	return githubAsset{}, fmt.Errorf("no such asset: %s", name)
}

type githubAsset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func downloadFile(url string, dest string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	// Setting Accept header to application/octet-stream instructs GitHub's API to send the asset's binary.
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
