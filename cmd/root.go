/*
Copyright © 2024 takak
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"regexp/syntax"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "randrex",
	Short: "Usage: randrex [OPTION]... Regexp",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runRandrex,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.randrex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("pattern", "p", "", "")
	rootCmd.Flags().IntP("number", "n", 1, "Number of print")
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

		// Search config in home directory with name ".randrex" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".randrex")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func runRandrex(cmd *cobra.Command, args []string) {
	pattern, _ := cmd.Flags().GetString("pattern")

	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return
	}

	generatedStr := genRandomStringFromRegex(re)
	fmt.Println(generatedStr)
}

func genRandomStringFromRegex(re *syntax.Regexp) string {
	var result []rune
	generateFromNode(re, &result)
	return string(result)
}

func generateFromNode(node *syntax.Regexp, result *[]rune) {
	switch node.Op {
	case syntax.OpCharClass:
		*result = append(*result, rune(node.Rune[rand.Intn(len(node.Rune))]))
	case syntax.OpRepeat:
		for i := 0; i < node.Min; i++ {
			generateFromNode(node.Sub[0], result)
		}
	case syntax.OpConcat:
		for _, sub := range node.Sub {
			generateFromNode(sub, result)
		}
	case syntax.OpAlternate:
		generateFromNode(node.Sub[rand.Intn(len(node.Sub))], result)
	case syntax.OpCapture:
		for _, sub := range node.Sub {
			generateFromNode(sub, result)
		}
	default:
		// 他のノード型は無視する
	}
}
