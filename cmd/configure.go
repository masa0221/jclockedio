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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configure called")
		configInit()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configInit() {
	maskType := newDataMaskType()
	config.Jobcan.Email = readInput("Jobcan E-mail", config.Jobcan.Email, maskType.Partial)
	config.Jobcan.Password = readInput("Jobcan Password", config.Jobcan.Password, maskType.Password)

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
