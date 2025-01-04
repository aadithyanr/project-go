package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

type GitHubResponse struct {
	DefaultBranch string `json:"default_branch"`
}

func VerifyBranchName(username string, reponame string, branchName string) bool {
	response, err := http.Get("https://api.github.com/repos/" + username + "/" + reponame + "/branches/" + branchName)
	if err != nil {
		check(err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

func GetMainBranchName(username string, reponame string) (string, error) {
	response, err := http.Get("https://api.github.com/repos/" + username + "/" + reponame)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)
	var obj GitHubResponse
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}
	return obj.DefaultBranch, nil
}

func InitRepo(path string) error {
	gitBin, _ := exec.LookPath("git")

	cmd := &exec.Cmd{
		Path:   gitBin,
		Args:   []string{gitBin, "init", path},
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}
	err := cmd.Run()
	return err
}

func CheckIfGitInstalled() bool {
	gitBin, _ := exec.LookPath("git")
	return gitBin != ""
}

func main() {
	var dstPath string
	var branchName string
	isGitInstalled := CheckIfGitInstalled()

	if !isGitInstalled {
		fmt.Println("Git is not installed. Please install git and try again.")
		return
	}

	branchNamePtr := flag.String("branch", "master", "Branch name to clone")
	outPtr := flag.String("out", "", "Destination path")
	shouldNotInitPtr := flag.Bool("no-init", false, "Don't initialize the repository")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Please provide a repository URL")
		return
	}

	repoUrl := args[0]
	splitGitUrl := strings.Split(repoUrl, "/")
	if len(splitGitUrl) < 5 {
		fmt.Println("Invalid repository URL format")
		return
	}
	username := splitGitUrl[3]
	reponame := splitGitUrl[4]

	if *outPtr == "" {
		mainBranchName, err := GetMainBranchName(username, reponame)
		check(err)
		branchName = mainBranchName
	} else {
		existsBranch := VerifyBranchName(username, reponame, *branchNamePtr)
		if !existsBranch {
			fmt.Println("Branch name does not exist. Please provide a valid branch name.")
			os.Exit(1)
		}
		branchName = *branchNamePtr
	}

	tempDir := ".loading-temp"
	err := os.Mkdir(tempDir, 0755)
	check(err)
	defer os.RemoveAll(tempDir)

	fileUrl := "https://github.com/" + username + "/" + reponame + "/archive/refs/heads/" + branchName + ".zip"
	zipFilePath := filepath.Join(tempDir, "repo.zip")
	DownloadFile(zipFilePath, fileUrl)

	unzipPath := filepath.Join(tempDir, "unzipped")
	Unzip(zipFilePath, unzipPath)

	files, err := ioutil.ReadDir(unzipPath)
	check(err)
	if len(files) == 0 {
		fmt.Println("No files found in the repository archive")
		return
	}

	repoDir := files[0]
	repoDirPath := filepath.Join(unzipPath, repoDir.Name())
	err = os.Rename(repoDirPath, dstPath)
	check(err)

	if isGitInstalled && !*shouldNotInitPtr {
		InitRepo(dstPath)
	}
}
