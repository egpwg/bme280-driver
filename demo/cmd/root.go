package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "BME280 driver testing tool",
	Long:  `BOSCH BME280 driver testing tool`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		exit := false
		var err error

		for !exit {
			exit, err = Command(func() string {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("bme280 > ")
				in, err := reader.ReadString('\n')
				if err != nil {
					panic(err)
				}
				return strings.TrimSpace(in)
			}())
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		// err = model.CloseBus()
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Print temperature & humidity & pressure values",
	Long:  `Print all sensor values: temperature & humidity & pressure`,
	Run: func(cmd *cobra.Command, args []string) {
		err := All()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var temperatureCmd = &cobra.Command{
	Use:   "t",
	Short: "Print temperature value",
	Long:  `Print sensor value: temperature`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Temperature()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var humidityCmd = &cobra.Command{
	Use:   "h",
	Short: "Print humidity value",
	Long:  `Print sensor value: humidity`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Humidity()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var pressureCmd = &cobra.Command{
	Use:   "p",
	Short: "Print pressure value",
	Long:  `Print sensor value: pressure`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Pressure()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(allCmd, temperatureCmd, humidityCmd, pressureCmd)

	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.demo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".demo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".demo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
