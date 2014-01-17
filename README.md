#Crane

Crane is a management tool for Docker containers written in Go. It's main goal is to help orchestrate multiple Docker containers. It is a command line utility that allows to automate execution of multiple commands(both Docker's commands and container's commands) across a cluster of containers. 

All of it can be packed into a single executable file that can be easily distributed!


#Why Crane?

 Are you tired of setting up a cluster of containers every time you start your machine?
 
 Would you like to be able to recreate your containers on a different machine with a single command?
 
 Would you like to be able to run a set of commands across many containers with a single command?

If you answered "yes" to at least one of above questions then Crane is something you might want to give a spin!


#Features

</br>
- Run a number of docker commands (pull,commit,run and more) on a cluster of containers with one simple command! You can run a number of unix commands across multiple containers with a single command!
- One executable is all you need!
- Build/Destroy multi-container docker environments using single command.
- Easily launch and manage multiple copies of the same environment.
- Declarative TOML format file to specify a cluster of container configurations. You can share this file with others and recreate your environment with a single command!
- Easily launch multiple instances of the same container.
- Share data between the host machine and containers running in the environment
- Create daemonized containers that can be easily accessed and modified afterwards.
- Ability to freeze containers into immutable images and saving your container's changes step by step (similar to Dockerfile commits after each command).
- SSH to daemonized containers.
- Many more...


#Example Usage

Crane was created to improve collaboration between developers. With a single configuration file (let's face it: it does not get simpler than TOML format) and a single executable you can export and recreate your environment on another machines with a single command!

You can define a cluster of containers and run multiple commands inside of them (possibly different per each container) using (again) a single command!


# Requirements

Crane manages Docker containers so you have to install Docker first. Detailed Docker installation instructions can be found [here](http://docs.docker.io/en/latest/installation).

You need [Go](http://golang.org/) set up as well.

# How to use Crane

Say NO! to complicated and lengthy installation and configuration processes!

All Crane source files and dependencies are compiled into a single executable file! This way you don't need to worry about installing and configuring any dependencies yourself.

In order to build the executable from source you need to run:
```
go build crane.go commands.go version.go
```
This will generate an executable called "crane". In order to be able to use it you need to make it executable:

```
chmod a+x crane
```

Now you are set to go! Just check everything is ok:
```
$./crane version
Crane v0.0.1.dev
```

You can also use ```go install```  and add your $GOPATH to the $PATH so you can run crane anywhere in the filesystem without the "./" prefix.

You can configure Docker containers using a single file called Cranefile.toml. Crane can generate a simple Cranefile.toml for you if you feel lazy. 

In case you missed it I say it again: you need only a SINGLE executable file and that's it! Simplesss...

#Example usage

Example Cranefile.toml

```
[containers]
[containers.firstContainer]
IMAGE = "orobix/sshfs_startup_key2"
DOCKERFILE = "."
GRAPHICAL = true
DAEMONIZED = true
CWD = "/home/foo" #Leave empty if not needed
DNS = ""#172.25.0.10" #Leave empty if not needed
PASSWORD = "orobix2013"#Leave empty if not needed
USERNAME = "root"
PORTS = [[49153, 22],[49653, 80]]

MOUNTPOINTS = [["/home/piotr/node-simple", "/mnt/node-simple","rw"],["/home/piotr/colors","/mnt/colors","ro"]]

COMMANDS = [["first","echo firstContainerfirstScript"],["second","echo firstContainerSecondScript"]]

[containers.secondContainer]
IMAGE = "orobix/sshfs_startup_key2"
DOCKERFILE = "."
GRAPHICAL = true
DAEMONIZED = false
CWD = "" #"/home/foo" #Leave empty if not needed
DNS = "" #172.25.0.10" #Leave empty if not needed
PASSWORD = "orobix2013"#Leave empty if not needed
USERNAME = "root"
PORTS = [[49154, 22],[49654, 80]]

MOUNTPOINTS = [["/home/piotr/node-simple", "/mnt/node-simple","rw"],["/home/piotr/colors","/mnt/colors","ro"]]
COMMANDS = [["first","echo secondContainerfirstScript"],["second","echo secondContainerSecondScript"]]



```
Example dialog involving daemonized containers:

```
$sudo docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
$crane version
Crane v0.0.1.dev
$crane start -a
2014/01/16 12:39:19 start.go:89: ▶ N 0x1a  Successfully started container "firstContainer"...
2014/01/16 12:39:20 start.go:89: ▶ N 0x30  Successfully started container "secondContainer"...
$sudo docker ps
CONTAINER ID        IMAGE                              COMMAND                CREATED             STATUS              PORTS                                          NAMES
540dbe8b1c72        orobix/sshfs_startup_key2:latest   /bin/bash -c /usr/sb   5 seconds ago       Up 5 seconds        0.0.0.0:49154->22/tcp, 0.0.0.0:49654->80/tcp   clever_tesla        
190b33bd4137        orobix/sshfs_startup_key2:latest   /bin/bash -c /usr/sb   7 seconds ago       Up 6 seconds        0.0.0.0:49153->22/tcp, 0.0.0.0:49653->80/tcp   boring_shockley     
$crane enter firstContainer
root@190b33bd4137:~# ls
ls
testfile
root@190b33bd4137:~# touch Newfile
touch Newfile
root@190b33bd4137:~# ls
ls
Newfile  testfile
root@190b33bd4137:~# exit
exit
exit
$crane destroy
2014/01/16 12:40:07 destroy.go:113: ▶ N 0xb  Kill command output:
190b33bd4137908a5f1bede19c8e45f2a827b826ca101db76a45c4015756869f
540dbe8b1c72ac41f7b438b68f0fb2bf60bc9a50b359a45b0f5d4ce6ce9ef8f7
2014/01/16 12:40:08 destroy.go:126: ▶ N 0xd  Remove command output:
190b33bd4137908a5f1bede19c8e45f2a827b826ca101db76a45c4015756869f
540dbe8b1c72ac41f7b438b68f0fb2bf60bc9a50b359a45b0f5d4ce6ce9ef8f7
$sudo docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

Running multiple commands across multiple containers:

```
$crane runall -c="first;second"
2014/01/16 12:59:29 utils.go:13: ▶ N 0x1a  
#####COMMAND OUTPUT######
firstContainerfirstScript
firstContainerSecondScript

#####END OF COMMAND OUTPUT#####

2014/01/16 12:59:30 utils.go:13: ▶ N 0x35  
#####COMMAND OUTPUT######
secondContainerfirstScript
secondContainerSecondScript

#####END OF COMMAND OUTPUT#####

$crane runall -o="echo works"
2014/01/16 13:00:07 utils.go:13: ▶ N 0x1a  
#####COMMAND OUTPUT######
works

#####END OF COMMAND OUTPUT#####

2014/01/16 13:00:08 utils.go:13: ▶ N 0x35  
#####COMMAND OUTPUT######
works

#####END OF COMMAND OUTPUT#####


```
Commiting containers and removing images
```
$crane start firstContainer
2014/01/16 13:04:01 start.go:89: ▶ N 0x1a  Successfully started container "firstContainer"...

$crane enter firstContainer
root@ef040db83aec:~# touch haha
touch haha
root@ef040db83aec:~# ls
ls
haha  testfile
root@ef040db83aec:~# exit
exit
exit

$crane freeze firstContainer::updatedDockerImage
2014/01/16 13:05:54 utils.go:13: ▶ N 0x9  
#####COMMAND OUTPUT######
b1de5e229f9d4e69ef93e1c462aed3b2a5667db9a4d7c60d228afc02916a9a47

#####END OF COMMAND OUTPUT#####

$sudo docker images| grep "updatedDockerImage"
updatedDockerImage                       latest              b1de5e229f9d        25 seconds ago      1.797 GB

$crane rmi updatedDockerImage
2014/01/16 13:07:10 utils.go:13: ▶ N 0x8  
#####COMMAND OUTPUT######
Untagged: b1de5e229f9d4e69ef93e1c462aed3b2a5667db9a4d7c60d228afc02916a9a47
Deleted: b1de5e229f9d4e69ef93e1c462aed3b2a5667db9a4d7c60d228afc02916a9a47

#####END OF COMMAND OUTPUT#####

$sudo docker images| grep "updatedDockerImage"
```
Updating image of container and pulling multiple images at once:

```
$crane run -u firstContainer:#"echo first; echo second"
2014/01/16 13:29:13 utils.go:13: ▶ N 0x1c  
#####COMMAND OUTPUT######
first
second

#####END OF COMMAND OUTPUT#####

2014/01/16 13:29:13 freeze.go:86: ▶ N 0x2d  Successfully froze container "firstContainer" into image "alpha" with id "9dec4701263f76206323625ef05e3bf19082ae975981d85ddbc3714301ed6033"

$crane build -a
2014/01/16 13:30:01 utils.go:13: ▶ N 0xe  
#####COMMAND OUTPUT######
Successfully build following images:
[alpha orobix/sshfs_startup_key2]
#####END OF COMMAND OUTPUT#####
```


# Configuration

##Cranefile

Crane can be configured using Cranefile.toml. Cranefile must be located in the same directory as the 'crane' executable bianry file.

Example Cranefile.toml:

    [containers]
    [containers.firstContainer]
    IMAGE = "orobix/sshfs_startup_key2"
    DOCKERFILE = "."
    GRAPHICAL = true
    DAEMONIZED = true
    CWD = "/home/foo"
    DNS = "172.25.0.10"
    USERNAME = "root"
    PASSWORD = "orobix2013"
    PORTS = [[49153, 22],[49653, 80]]
    MOUNTPOINTS = [["/home/piotr/node-simple", "/mnt/node-simple","rw"],["/home/piotr/colors","/mnt/colors","ro"]]
    COMMANDS = [["first","echo firstContainerfirstScript"],["second","echo firstContainerSecondScript"]]

Flags explained:

IMAGE (string) The name of the docker image. The docker image takes precedence before the Dockerfile. This means that if a provided image already exists in the system it will be used by the crane without building a new image using Dockerfile. If the image does not exist crane will search for it in the docker public repository.If it will not be able to find it there it will build a new image using Dockerfile provided with the image name provided.To sum up the order in which crane decides which image to use is as follows:

1. Use image from host system.
2. Search for image in the docker public repository.If it is present in the public repository it will be automatically pulled during run.
3. Build new image using the Dockerfile.

DOCKERFILE (string) Path to the directory where Dockerfile is stored. If the Dockerfile is located in the same directory as executable binary use ".".

GRAPHICAL (boolean) To Be Added soon.

DAEMONIZED (boolean) Determines if a container is daemonized or not.

CWD(string) Working directory inside the container. Leave empty("") is not needed.

DNS(string) Address of the DNS. Leave empty("") if not needed.

USERNAME(string) Username used when inside a container(make sure that the used image has this user set up). Default username is "root"

PASSWORD(string) Some images require password to login. Leave empty("") is not needed. 

PORTS(array of digit arrays) A list of port redirection.A port redirect is specified as [PUBLIC,PRIVATE], where TCP port PUBLIC will be redirected to TCP port PRIVATE.You can create multiple port redirections.The public port can be omitted, in which case a random public port will be allocated. Please remember that deamonized containers require port 22 to be exposed in order to interact with the host.

MOUNTPOINTS(array of string arrays) A list of mountpoints. Every element consists of 3 arguments: a first argument is the host mountpoint's absolute path, second argument is the remote mountpoint's absolute path(it will be created if not existing) and the last argument specifies if the host's filesysytem: should be mounted as a read-only ("ro") of read-write ("rw").

COMMANDS(array of string arrays) A list of commands to be executed inside the container. Every element consists of 2 elements: command identifier and command itself.

##State file
The state file (.crane) is used by Crane to keep track of all containers. It MUST NOT be modified by the user (unless you know what you are doing).

The structure of the file is pretty straightforward:

    [statecontainers]
    [statecontainers.firstContainer]
    ID = "12233445"
    IP = "172.234.1.1"
    [statecontainers.secondContainer]
    ID = "123456df"
    IP = "not_daemonized_has_no_ip"

Flags explained:

ID - holds a container ID.

IP - holds a container IP. Please not that non-daemonized and stopped containers won't have an IP address. In state file it will be reflected with the value "not_daemonized_has_no_ip". 

## Crane Commands

###Build
Builds a docker image using Dockerfiles specified in the Cranefile.

    crane build [options] <containerName1> <containerName2>
    
Builds docker images using Dockerfiles specified in Cranefile for <containerName1> and <containerName2>

*Note*<br />
Docker is a smart beast and it will detect if you are trying to build an image using the same Dockerfile multiple times with different image names. In such case Docker will use already existing image to create a new image with different name but in fact it is the same image all way long with multiple names. If interested you can find out more about docker images/layers magic [here](http://docs.docker.io/en/latest/terms/image/).  

*Another Note*<br />
IF you are testing Crane locally AND would like to build docker images using Dockerfiles AND are stuck in a network that does not let you ping even Google DNS 8.8.8.8 in order to build docker images using Dockerfile you have to restart docker daemon with appropriate DNS server:<br />

In terminal separable to the terminal when you run Crane run docker daemon in with the dns option: 

    sudo docker -d -dns x.x.x.x
    
You may have to delete the previous docker process first by using:

    sudo rm /var/run/docker.pid


Options:

-a (--all)   Builds docker images using Dockerfiles from all containers specified in the Cranefile.



###Create
Generate an example Cranefile.toml

    crane create

Please note that if Cranefile already exists it will be overwritten. For more details about the Cranefile format and use see section  " Configuration-Cranefile" above.

###Debug Mode
You can use crane in a debug mode which provides much more information about what is happenning behind the scenes and can help you diagnose problems.

    crane -d
    
Using "-d" option will run crane in the debug mode. This option can be used in conjunction will all commands presented below.


###Destroy

Destroy uses "docker kill" and "docker rm" commands to destroy containers. 

Firstly crane kills the running/stopped container(s) using "docker kill" command. As a result the container is no longer running but can be started again. Then crane uses "docker rm" command to completely remove the container from the system. 

Please note that you can destroy only containers that were created by the crane itself (not manually by docker) and present in the state file(in order to perform an action on a container you must have container ID).

    crane destroy
    
Destroys all containers defined in the Cranefile.

    crane destroy <Container1> <Container2>
    
Destroys containers specified by the user.

###Enter
Presents the user with the interactive command line prompt inside a chosen container (you can enter only one container at a time).

    crane enter <containerName>

In case of daemonized containers it is necessary to "start" them first before trying to enter them.


###Freeze

    crane freeze [options] <containerName1> <containerName2>
Transforms (using docker lingo "commits") containers into immutable images using names defined in Cranefile (overwrites existing images).

    crane freeze [options] <containerName1>::<imageName1> <containerName2>::<imageName2>
Transforms (using docker lingo "commits") containers into immutable images with chosen names.

Both options can be used within a single crane command:

    crane freeze [options] <containername1> <containerName2>:<imageName2>

Options:

-a (--all) : Freeze all containers defined in the Cranefile into immutable images. Useful when you want to save all your work in one go.

Please note that this command will use existing image names as per Cranefile(overwrite existing images). If you want to freeze a container with a specific name that is different from the original image please use the "freeze" command.


###Pull
Pulls images from the docker public repository.

    crane pull [options] <image1> <image2>
    
Pulls chosen images from the public repository.
    
Options:

-a(--all) : Pulls all images defined in the Cranefile from the docker public repository.

###Rmi
Removes docker images specified by the user.

    crane rmi [options] <containerName1> <containerName2>
    
Removes docker images specified in the Cranefile for <containerName1> and <containerName2>.

Options:

-a (--all) Remove all images specified in the Cranefile for all containers.

###Run

Run command is defined in Cranefile or typed from the command line inside specified containers.

In case of daemonized containers you need to start them first in order to run commands inside of them.

    crane run [options] <containerName1>:<commandId1>,<commandId2> <containerName2>:<commandId1>
    
Run commands defined in the Cranefile for each container.

    crane run [options] <containerName1>:#"<command1>;<command2>" ... etc

Run commands entered by the user in the command line (no need to add them to the Cranefile).
   
Please note different delimiters in both cases.

The user can use the mix both methods for different containers using a single command:

    crane run [options] <containerName1>:<command1>,<command2> <containerName2>:#"<command1>;<command2>"

Options:

-s="" (--save="") : In order to make the changes made to the container permanent(e.g. you modified a file inside a container that you want to stay this way next time you use that container) you need to transform the container into an immutable image. You can either create a completely new image by providing the name of the new container through using this option or you can simply overwrite the existing image by adding your changes to it e.g.:

    crane run -s="goose,dog" gosling:firstCommand,secondCommand puppy:anotherCommand
    
If you are going to create a completely new image that you are going to reuse with the same container name please remember to change the Cranefile appropriately to reflect those changes(chosen container Image variable).

-u (--update) : It's like an -s option which always overwrites existing images. When present containers are transformed to images which will replace existing images used to create those containers.
    
###Runall
    
This command is used to automate multiple run commands.
    
    
    crane runall

    
Runs all commands in all containers defined in the Cranefile.
    
    crane runall -c="<CranefileCommand1>;<CranefileCommand2>"
    
Runs specified Cranefile commands in all containers.
    

    crane runall -o="<ownCommand>;<ownCommand>"
    
Runs own (not defined in the Cranefile but typed on the commandline) commands in all containers.
    
    
    crane runall -l="<containerName1>;<containerName2>"
    
Runs all Cranefile's commands in specified containers.
    
Please note that All arguments are parsed sequentially so e.g. 
    
    crane runall -l="firstContainer;secondContainer"
    
Will execute all commands in the firstContainer and once finished it will execute all commands in the second container.
    
The options **cannot** be used simultaneously within a single "runall" command (but tou can call runall multiple times if you need to use multiple options).

###Start
        
    crane start [options] <containerName1> <containerName2>
    
Starts all daemonized containers defined in the Cranefile.

Options:

-a(--all) : Starts all daemonized containers defined in the Cranefile.

-f(--force) : Crane assumes that the image already exists in the host system and does not attempt to download it from the docker public repository (useful for offline mode).


###Version

    crane version
    crane -v
    crane --version
Shows the current version of Crane.

#Status 
</br>
Early development, use at your own risk!

#Licence
</br>
MIT license http://www.opensource.org/licenses/mit-license.php/

Copyright (C) 2014 Piotr Chudzik 

Orobix srl http://www.orobix.com

This project was created as a part of the [REVAMMAD](http://revammad.blogs.lincoln.ac.uk/) project funded by the Marie Curie ITN. 
