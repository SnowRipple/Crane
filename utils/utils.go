package utils

import (
	"fmt"
	log "github.com/SnowRipple/crane/logger"
)

var logger = log.GetLogger()

//Prints output of the command in a form that is easy to read
func PrintCommandOutput(output []byte) {

	logger.Notice("\n#####COMMAND OUTPUT######\n%v\n#####END OF COMMAND OUTPUT#####\n", string(output))
}

//Extracts error message from the container output if present.
//If not it uses host's error message only
func ExtractContainerMessage(message []byte, err error) string {

	var cumulativeMessage string
	//Container's error message
	if message != nil && len(message) > 0 {
		cumulativeMessage = fmt.Sprintf("\nContainer's message:\n%s", string(message))
	}
	//Host's error message
	if err != nil {
		cumulativeMessage = cumulativeMessage + fmt.Sprintf("\nHost Error:\n%s", err.Error())
	}

	return cumulativeMessage
}
