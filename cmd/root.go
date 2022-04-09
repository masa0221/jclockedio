/*
Copyright © 2022 Masashi Tsuru

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var config Config

type JobcanConfig struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
}

type ChatworkConfig struct {
	Send     bool   `toml:"send"`
	ApiToken string `toml:"apitoken"`
	RoomId   string `toml:"roomid"`
}

type OutputConfig struct {
	Format string `toml:"format"`
}

type Config struct {
	Jobcan   JobcanConfig   `toml:"Jobcan"`
	Chatwork ChatworkConfig `toml:"Chatwork"`
	Output   OutputConfig   `toml:"Output"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jclockedio",
	Short: "Jobcan clocked in/out tool",
	Long: `This is a things for clocked in/out to Jobcan, then it results are send to Chatwork.
"jclockedio" means Jobcan clocked in/out.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jclockedio)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode. Output details of process")
}

func initConfig() {
	viper.SetConfigType("toml")

	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jclockedio" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jclockedio")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Config file is not found.\nIt be execute configure process.")
		os.Exit(0)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Can't unmarshal config:", err)
		os.Exit(1)
	}
}
