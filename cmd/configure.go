/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Generate and regenerate configure file",
	Long: `Generate and regenerate configure file

Output format:
Output format is use for adit's stdout and Chatwork message.

You can use variables below.
----
{{ .clock }}           -- clock time
{{ .beforeStatus }}    -- before clocked in/out status
{{ .afterStatus }}     -- after clocked in/out status
----
For example:
[{{ .clock }}] {{ .beforeStatus }} => {{ .afterStatus }}
`,
	Run: func(cmd *cobra.Command, args []string) {
		configInit()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func configInit() {
	maskType := newDataMaskType()
	config.Jobcan.Email = readInput("Jobcan E-mail", config.Jobcan.Email, maskType.Partial)
	config.Jobcan.Password = readInput("Jobcan Password", config.Jobcan.Password, maskType.Password)
	config.Output.Format = readInput("Output format", config.Output.Format, maskType.None)

	if readInputYN("Do you send to Chatwork?") {
		config.Chatwork.Send = true
		config.Chatwork.ApiToken = readInput("Chatwork API Token", config.Chatwork.ApiToken, maskType.Partial)
		config.Chatwork.RoomId = readInput("Chatwork room_id", config.Chatwork.RoomId, maskType.None)
	} else {
		config.Chatwork.Send = false
	}
	filepath := os.ExpandEnv("$HOME") + "/.jclockedio"
	saveConfig(filepath, config)
	fmt.Printf("\nCreated!(%v)\nEnjoy your work🌸\n", filepath)
}

type DataMaskType struct {
	Password string
	Partial  string
	None     string
}

func newDataMaskType() DataMaskType {
	return DataMaskType{"password", "partial", "none"}
}

func readInputYN(label string) bool {
	fmt.Printf("%v (y)es/(n)o: ", label)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	answer := false
	if input != "" && (input[0:1] == "Y" || input[0:1] == "y") {
		answer = true
	}

	return answer
}

func readInput(label string, defaultValue string, maskType string) string {
	msg := ""
	if defaultValue == "" {
		msg = fmt.Sprintf("%v: ", label)
	} else {
		outputValue := ""
		dataMaskType := newDataMaskType()
		if maskType == dataMaskType.Password {
			outputValue = strings.Repeat("*", len(defaultValue))
		} else if maskType == dataMaskType.Partial {
			strLength := len(defaultValue) - 2
			if strLength <= 0 {
				outputValue = strings.Repeat("*", len(defaultValue))
			} else {
				outputValue = defaultValue[0:1] + strings.Repeat("*", strLength) + defaultValue[strLength+1:]
			}
		} else {
			outputValue = defaultValue
		}
		msg = fmt.Sprintf("%v [%v]: ", label, outputValue)
	}
	fmt.Printf(msg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()
	if input == "" {
		input = defaultValue
	}

	return input
}

func saveConfig(filepath string, config Config) {
	bs, err := toml.Marshal(config)
	if err != nil {
		fmt.Printf("Unable to marshal config: %v", err)
	}
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Can't create config file:", err)
		os.Exit(1)
	}
	fmt.Fprintln(file, string(bs))
}
