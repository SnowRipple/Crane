package container

import "fmt"

//Model of a container defined in the Cranefile.
type Container struct {
	Image       string
	Dockerfile  string
	Graphical   bool
	Daemonized  bool
	Cwd         string
	Dns         string
	Password    string
	Username    string
	Ports       [][]int
	Mountpoints [][]string
	Commands    [][]string
}

func (container *Container) String() string {
	return fmt.Sprintf("Image name: %s\nDockerfile: %s\nGraphical?: %t\nDaemonized?: %t\nWorking Directory: %s\nDNS: %%s\nPassword: %s\nUsername: %s\n,Ports: \n,%s\nMountpoints: \n%v\n,Commands:\n%v\n", container.Image, container.Dockerfile, container.Graphical, container.Daemonized, container.Cwd, container.Dns, container.Password, container.Username, container.Ports, container.Mountpoints, container.Commands)
}

//Model of a container defined in the .crane file.
type StateContainer struct {
	ID string
	IP string
}

func (stateContainer *StateContainer) String() string {
	return fmt.Sprintf("StateContainer ID: %s\nStateContainer IP: %s\n", stateContainer.ID, stateContainer.IP)
}
