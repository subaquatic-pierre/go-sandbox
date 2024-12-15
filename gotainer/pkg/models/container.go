package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/codeclysm/extract"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/google/uuid"
	"github.com/subaquatic-pierre/gotainer/lib"
	"github.com/subaquatic-pierre/gotainer/util"
)

type Container struct {
	ID         string
	Image      string
	Config     *Config
	Cmd        string
	Root       string
	SystemRoot int
	Mounts     []*string
	CgroupName string
	Pid        uint64
}

func NewContainer(image string) (*Container, error) {
	config := Config{
		ImageDir:     filepath.Join(lib.Cwd(), "./images"),
		ContainerDir: filepath.Join(lib.Cwd(), "./containers"),
	}
	id := uuid.New()

	// ensure image exists to run container
	imagePath := filepath.Join(config.ImageDir, fmt.Sprintf("%s/%s.tar.gz", image, image))
	if !lib.DirExists(imagePath) {
		return nil, fmt.Errorf("unable to create container, image doesn't exist, %s", image)
	}

	// create empty dir on creation in which to mount fs
	err := os.Mkdir(filepath.Join(config.ContainerDir, id.String()), 0755)
	if err != nil {
		return nil, err
	}

	// get command from image dir if exists
	containerCmd := ""
	buf, err := os.ReadFile(filepath.Join(config.ImageDir, "/%s/%s-cmd", image, image))
	if err == nil {
		containerCmd = string(buf)
	}

	cgroupName := fmt.Sprintf("%s%s", image, strings.Split(id.String(), "-")[0])

	return &Container{
		ID:         id.String(),
		Image:      image,
		Config:     &config,
		Cmd:        containerCmd,
		Root:       filepath.Join(config.ContainerDir, id.String()),
		CgroupName: cgroupName,
		Pid:        uint64(os.Getpid()),
	}, nil
}

func (c *Container) Exec() {
	call := c.Cmd
	log.Println("call: ", call)

	// create container command
	cmd := exec.Command(call)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run container command
	err := cmd.Run()
	if err != nil {
		log.Println("unable to run container command", err)
		return
	}
}

func (c *Container) MountFs() {
	// get image tar for container
	imagePath := c.imagePath()
	r, err := os.Open(imagePath)
	if err != nil {
		log.Printf("directory path does not exist: %s\n", err)
		return
	}
	defer r.Close()

	ctx := context.Background()
	archiveErr := extract.Archive(ctx, r, c.Root, nil)
	if archiveErr != nil {
		log.Println("unable to extract image", err)
	}
}

func (c *Container) NewCgroup() {
	res := cgroup2.Resources{}
	cgroupName := fmt.Sprintf(c.CgroupName)
	m, err := cgroup2.NewSystemd("/", cgroupName, -1, &res)
	if err != nil {
		log.Println("error creating cgroup", err)
		return
	}

	m.AddProc(c.Pid)

	log.Printf("created cgroup: %s and added process: %d\n", c.CgroupName, c.Pid)

}

func (c *Container) Chroot() {
	// get handle to old root
	oldRootHandle, err := os.Open("/")
	if err != nil {
		panic(err)
	}
	// close handle after return
	defer oldRootHandle.Close()

	c.SystemRoot = int(oldRootHandle.Fd())

	// change root to new container root
	chdirErr := syscall.Chdir(c.Root)
	if chdirErr != nil {
		log.Println("unable to change directory into new root", chdirErr)
		return
	}
	syscall.Chroot(c.Root)

	log.Println("changed root to", c.Root)
}

// switch back to root and dir
func (c *Container) ResetRoot() {
	syscall.Fchdir(int(c.SystemRoot))
	syscall.Chroot(".")

	log.Println("reset root to .")
}

// mount "proc"
func (c *Container) MountVolumes(mounts []string) {
	for i := range mounts {
		mount := mounts[i]
		util.Must(syscall.Mount(mount, mount, mount, 0, ""))
	}

	log.Println("mounted volumes", mounts)
}

// unmount "proc"
func (c *Container) UnMountVolumes() {
	for vol := range c.Mounts {
		syscall.Unmount(*c.Mounts[vol], 0)
	}

	log.Println("unmounted volumes", c.Mounts)
}

func (c *Container) UnmountFs() {
	// Remove container dir
	err := os.RemoveAll(c.Root)
	if err != nil {
		log.Println(err)
	}

	log.Println("unmounted fs")
}

func (c *Container) RemoveCgroup() {
	// Remove Cgroup
	m, err := cgroup2.LoadSystemd("/", c.CgroupName)
	if err != nil {
		log.Println("error loading cgroup for deletion", err)
		return
	}
	err = m.DeleteSystemd()
	if err != nil {
		log.Println("error deleting cgroup", err)
		return
	}

	log.Println("deleted cgroup for container", c.CgroupName)

}

func (c *Container) Init() {
	c.MountFs()

	// create new cgroup
	c.NewCgroup()

	// change root
	c.Chroot()

	// mount volumes
	volumes := []string{"proc"}
	c.MountVolumes(volumes)
}

func (c *Container) Shutdown() {
	c.UnMountVolumes()

	c.ResetRoot()

	c.UnmountFs()

	log.Println("current directory", lib.Cwd())

	c.RemoveCgroup()
}

// func (c *Container) cgroupPath() string {
// 	cgroups := "/sys/fs/cgroup/"
// 	groupSlice := fmt.Sprintf("%s.slice", c.ID)
// 	groupPath := filepath.Join(cgroups, groupSlice)
// 	return groupPath
// }

func (c *Container) imagePath() string {
	return filepath.Join(c.Config.ImageDir, fmt.Sprintf("%s/%s.tar.gz", c.Image, c.Image))
}
