package main

import (
        "log"
        "os"
        "os/exec"
        "syscall"
)

func Run() {
        cmd := exec.Command("/proc/self/exe", "init")

        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.Env = []string{"PS1=[omochabako]%"}

        cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWNS |
                        syscall.CLONE_NEWUTS |
                        syscall.CLONE_NEWIPC |
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
                log.Fatalf("Error running the cmd.Run - %s\n", err)
        }
        os.Exit(0)
}

func Initialisation() error {
        if err := syscall.Sethostname([]byte("omochabako")); err != nil {
                log.Fatalf("Error setting hostname - %s\n", err)
        }
        if err := syscall.Exec("/bin/sh", []string{"/bin/sh"}, os.Environ()); err != nil {
                log.Fatalf("Error exec: $s\n", err)
        }
        return nil
}

func Usage() {
        log.Fatalf("Usage: %s run\n", os.Args[0])
}

func main() {
        switch os.Args[1] {
        case "run":
                Run()
        case "init":
                if err := Initialisation(); err != nil {
                        log.Fatalf("Error: %s init\n", err)
                }
                os.Exit(0)
        default:
                Usage()
        }
}
