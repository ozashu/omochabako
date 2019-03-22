package main

import (
        "fmt"
        "os"
        "os/exec"
        "syscall"

        "github.com/docker/docker/pkg/reexec"
)

func init() {
        reexec.Register("omInitialisation", omInitialisation)
        if reexec.Init() {
                os.Exit(0)
        }
}

func omInitialisation() {
        fmt.Printf("\n>> namespace setup code goes here <<\n\n")
        omRun()
}

func omRun() {
        cmd := exec.Command("/bin/sh")

        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.Env = []string{"PS1=[omochabako]%"}

        if err := cmd.Run(); err != nil {
                fmt.Printf("Error running the /bin/sh command - %s\n", err)
                os.Exit(1)
        }
}

func main() {
        cmd := reexec.Command("omInitialisation")

        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.Env = []string{"PS1=[omochabako]%"}

        cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWNS |
                        syscall.CLONE_NEWUTS |
                        syscall.CLONE_NEWPID |
                        syscall.CLONE_NEWNET |
                        syscall.CLONE_NEWUSER,
                UidMappings: []syscall.SysProcIDMap{
                        {
                                ContainerID: 0,
                                HostID:      os.Getuid(),
                                Size:        1,
                        },
                },
                GidMappings: []syscall.SysProcIDMap{
                        {
                                ContainerID: 0,
                                HostID:      os.Getgid(),
                                Size:        1,
                        },
                },
        }
        if err := cmd.Run(); err != nil {
                fmt.Printf("Error running the reexec.Command - %s\n", err)
                os.Exit(1)
        }
}
