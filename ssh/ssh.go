package ssh

import (
	"code.google.com/p/go.crypto/ssh"
	log "github.com/SnowRipple/crane/logger"
	"io"
	"os"
)

var logger = log.GetLogger()

type clientPassword string

func (p clientPassword) Password(user string) (string, error) {
	return string(p), nil
}

type TerminalModes map[uint8]uint32

const (
	VINTR         = 1
	VQUIT         = 2
	VERASE        = 3
	VKILL         = 4
	VEOF          = 5
	VEOL          = 6
	VEOL2         = 7
	VSTART        = 8
	VSTOP         = 9
	VSUSP         = 10
	VDSUSP        = 11
	VREPRINT      = 12
	VWERASE       = 13
	VLNEXT        = 14
	VFLUSH        = 15
	VSWTCH        = 16
	VSTATUS       = 17
	VDISCARD      = 18
	IGNPAR        = 30
	PARMRK        = 31
	INPCK         = 32
	ISTRIP        = 33
	INLCR         = 34
	IGNCR         = 35
	ICRNL         = 36
	IUCLC         = 37
	IXON          = 38
	IXANY         = 39
	IXOFF         = 40
	IMAXBEL       = 41
	ISIG          = 50
	ICANON        = 51
	XCASE         = 52
	ECHO          = 53
	ECHOE         = 54
	ECHOK         = 55
	ECHONL        = 56
	NOFLSH        = 57
	TOSTOP        = 58
	IEXTEN        = 59
	ECHOCTL       = 60
	ECHOKE        = 61
	PENDIN        = 62
	OPOST         = 70
	OLCUC         = 71
	ONLCR         = 72
	OCRNL         = 73
	ONOCR         = 74
	ONLRET        = 75
	CS7           = 90
	CS8           = 91
	PARENB        = 92
	PARODD        = 93
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

const SSH_PORT = ":22"

func SshConnect(sshAddress, username, passwordString, sshCommand string) {

	logger.Debug("Trying to set up ssh connection with SSHAddress:" + sshAddress + ",Username:" + username + ",password:" + passwordString + ",SSH command:" + sshCommand + ".")

	// An SSH client is represented with a slete. Currently only
	// the "password" authentication method is supported.
	// To authenticate with the remote server you must pass at least one
	// implementation of ClientAuth via the Auth field in ClientConfig.

	password := clientPassword(passwordString)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.ClientAuth{
			// ClientAuthPassword wraps a ClientPassword implementation
			// in a type that implements ClientAuth.
			ssh.ClientAuthPassword(password),
		},
	}
	client, err := ssh.Dial("tcp", sshAddress+SSH_PORT, config)
	if err != nil {
		logger.Fatal("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	defer client.Close()
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		logger.Fatal("unable to create session: %s", err)
	}
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ECHO:          0,     // disable echoing
		TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		logger.Fatal("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		logger.Fatal("Unable to setup stdin for session: %v\n", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		logger.Fatal("Unable to setup stdout for session: %v\n", err)
	}

	stderr, err := session.StdoutPipe()
	if err != nil {
		logger.Fatal("Unable to setup stdout for session: %v\n", err)
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(stdin, os.Stdin)
	go io.Copy(os.Stderr, stderr)

	logger.Debug("The following command that will be executed during this SSH session:" + sshCommand)
	if err := session.Run(sshCommand); err != nil {
		logger.Fatal("SSH Error:Failed to run: " + err.Error())
	}
}
