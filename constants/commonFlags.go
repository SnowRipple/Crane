package constants

type CommonFlags struct {
	DebugMode bool `short:"d" long:"debug" description:"When crane is used in the debug mode a lot of extra information is provided during program execution." `

	Version    bool `short:"v" long:"version" description:"Shows the information about the crane version you are using."`
	ForceImage bool `short:"f" long:"force" description:"If chosen, crane will assume that the chosen image already exists in the host system(useful for offline mode)" `

	All bool `short:"a" long:"all" description:"Performs operation for all available arguments"`

	Save string `short:"s" long:"save" description:"Transforms(commits) containers into immutable images that can be reused in the future.For multiple containers provide comma separated names."`

	Update bool `short:"u" long:"update" description:"transforms(commits) containers into immutable images that willreplace existing images."`

	RunAllCommand string `short:"c" long:"commands" description:"To be used alongside runall.Run specified commands from the Cranefile across all containers"`

	RunAllContainer string `short:"l" long: "containers" description:"To be used alongside runall command.Run all commands in specified containers"`

	RunAllOwnCommand string `short:"o" long:"owncommands" description:"To be used alongside runall.Run specified own(not Cranefile) commands across all containers"`
}
