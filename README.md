# projectgo

**projectgo** is a command-line tool written in Go for downloading, unzipping, and initializing GitHub repositories. its main purpose is to simplify the process of cloning GitHub repositories, verifying branches, and setting up a local git repository. the project is still a work in progress and is currently incomplete. i'm actively working on fixing issues and adding more features.

## features

- download a zip archive of a specific branch from a github repository.
- extract the zip archive to a specified directory.
- verify if a specific branch exists in the repository.
- initialize the downloaded repository as a local git repository.

## usage

### prerequisites

1. ensure [go](https://golang.org/dl/) is installed on your system.
2. ensure [git](https://git-scm.com/downloads) is installed and added to your system's path.
3. ensure you have an active internet connection.

### running the tool

1. build the project:
   ```bash
   go build -o projectgo main.go
   ```

2. execute the tool:
   ```bash
   ./projectgo [flags] <repository-url>
   ```

### flags

- `-branch`: specify the branch name to download (default: `master`).
- `-out`: specify the destination folder for the downloaded repository.
- `-no-init`: skip git initialization after download.

### examples

1. download the default branch:
   ```bash
   ./projectgo https://github.com/aadithyanr/yomama-jokes-api
   ```

2. specify a branch to download:
   ```bash
   ./projectgo -branch=main https://github.com/aadithyanr/yomama-jokes-api
   ```

3. set a custom output directory:
   ```bash
   ./projectgo -out=my-repo https://github.com/aadithyanr/yomama-jokes-api
   ```

4. skip git initialization:
   ```bash
   ./projectgo -no-init https://github.com/aadithyanr/yomama-jokes-api
   ```

## current status

this project is **incomplete**, and several features and functionalities are not yet fully implemented or tested. i'm still actively working on fixing bugs and improving the tool. some known issues include:

- error handling is limited in some areas.
- branch verification may fail due to missing api checks.
- not all repository url formats are supported.

## contribution

feel free to fork this project and contribute by submitting issues or pull requests. any help is appreciated!

## license

this project is released under the [mit license](https://opensource.org/licenses/mit).
