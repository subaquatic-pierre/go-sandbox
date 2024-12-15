package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/codeclysm/extract"
	util "github.com/subaquatic-pierre/gotainer/util"
)

func MountImageFs(imageName, containerDir string) {
	imagePath := fmt.Sprintf("./assets/%s/%s.tar.gz", imageName, imageName)
	// fmt.Printf("Extracting %s %s\n", source, dst)
	r, err := os.Open(imagePath)
	if err != nil {
		log.Printf("directory path does not exist: %s\n", err)
		return
	}
	defer r.Close()

	ctx := context.Background()
	archiveErr := extract.Archive(ctx, r, containerDir, nil)
	if archiveErr != nil {
		log.Println("unable to extract image", err)

	}
}

func CreateTempDir(path string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	prefix := nonAlphanumericRegex.ReplaceAllString(path, "_")
	dir, err := os.MkdirTemp("./tmp", prefix)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("created tmp dir for container: %s\n", dir)
	return dir
}

func DirExists(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist)
}

func RemoveTempDir(path string) {
	if !DirExists(path) {
		log.Printf("directory path does not exist: %s\n", path)
		return
	}

	err := os.Remove(path)
	if err != nil {
		log.Println("there was an error removing temp dir", err)
		return
	}

	log.Println("removed tmp directory with path: ", path)
}

func Cwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func ExecContainer(root, call string) {
	// get handle to old root
	oldRootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	// close handle after return
	defer oldRootHandle.Close()

	// create container command
	cmd := exec.Command(call)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// change root to new container root
	chdirErr := syscall.Chdir(root)
	if chdirErr != nil {
		log.Println("unable to change directory into new root", chdirErr)
		return
	}

	fmt.Println("container dir", root)
	cwd := Cwd()

	containerPath := filepath.Join(cwd, root)
	fmt.Println("containerPath", containerPath)
	chrootErr := syscall.Chroot(containerPath)
	if chrootErr != nil {
		log.Println("unable to change root", chrootErr)
		return
	}

	// run container command
	cmdErr := cmd.Run()
	if err != nil {
		log.Println("unable to run container command", cmdErr)
		return
	}

	// switch back to root and dir
	syscall.Fchdir(int(oldRootHandle.Fd()))
	syscall.Chroot(".")

}

// func MountImageFs(imageName, containerDir string) {
// 	cmdStr := fmt.Sprintf("docker export $(docker create %s) | tar -xC . --one-top-level=%s", imageName, containerDir)

// 	cmd := exec.Command(cmdStr)

// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println("unable to extract docker image to container directory", err)
// 	}
// }

func PullImage(imageName string) {
	imagePath := fmt.Sprintf("./assets/%s", imageName)
	if DirExists(imagePath) {
		log.Printf("image already exists for %s", imageName)
		return
	}

	cmd := exec.Command("./pull.sh")
	cmd.Env = append(cmd.Env, fmt.Sprintf("image=%s", imageName))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("unable to pull image", err)
		return
	}

	log.Printf("content stored in ./assets/%s", imageName)
}

func NewCgroup(containerID string) int {
	cgroups := "/sys/fs/cgroup/"
	groupSlice := fmt.Sprintf("%s.slice", containerID)
	groupPath := filepath.Join(cgroups, groupSlice)
	os.Mkdir(groupPath, 0755)

	// Write -1 to group procs, means parent groupd
	util.Must(os.WriteFile(filepath.Join(groupPath, "cgroup.procs"), []byte("-1"), 0700))

	// Set pid max
	util.Must(os.WriteFile(filepath.Join(groupPath, "pids.max"), []byte("20"), 0700))

	// Removes the new cgroup in place after the container exits
	util.Must(os.WriteFile(filepath.Join(groupPath, "notify_on_release"), []byte("1"), 0700))

	cgroupHandle, _ := os.Open(groupPath)
	return int(cgroupHandle.Fd())
}
