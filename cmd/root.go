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

var PRINT_MAX int = 10

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "randrex",
	Short: "Usage: randrex [generate|parse] [OPTION]... Regexp",
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate strings matching a regular expression",
	Run:   runGenerateCommand,
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse regular expression",
	Run:   runParseCommand,
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
	generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().IntP("number", "n", 1, "Number of print")
	rootCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(parseCmd)
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

func runGenerateCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No arguments provided.")
		return
	}
	pattern := args[0]

	number, err := cmd.Flags().GetInt("number")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get number:", number)
		return
	}

	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse pattern:", pattern)
		return
	}

	for i := 0; i < number; i++ {
		generatedStr := genRandomStringFromRegex(re)
		fmt.Println(generatedStr)
	}
}

func genRandomStringFromRegex(re *syntax.Regexp) string {
	var result []rune
	generateFromNode(re, &result)
	return string(result)
}

func generateFromNode(node *syntax.Regexp, result *[]rune) {
	switch node.Op {
	case syntax.OpLiteral:
		*result = append(*result, node.Rune...)
	case syntax.OpCharClass, syntax.OpAnyCharNotNL:
		var randRunes []rune
		for i := 0; i < len(node.Rune); i += 2 {
			if node.Rune[i+1] != node.Rune[i] {
				randRunes = append(randRunes, rune(rand.Intn(int(node.Rune[i+1]-node.Rune[i]))+int(node.Rune[i])))
			} else {
				randRunes = append(randRunes, node.Rune[i])
			}
		}
		if len(randRunes) != 0 {
			*result = append(*result, rune(randRunes[rand.Intn(len(randRunes))]))
		} else {
			*result = append(*result, rune(rand.Intn(96)+31)) // from " " to "~" in ASCII
		}
	case syntax.OpStar, syntax.OpPlus, syntax.OpRepeat:
		min := 0
		max := PRINT_MAX
		switch node.Op {
		case syntax.OpRepeat:
			min = node.Min
			if node.Max >= 0 {
				max = node.Max
			}
		case syntax.OpPlus:
			min = 1
		}
		j := min
		if max != min {
			j = rand.Intn(max-min) + min
		}
		for i := 0; i < j; i++ {
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
		// ignore other case
	}
}

func runParseCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No arguments provided.")
		return
	}
	pattern := args[0]

	re, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse pattern:", pattern)
		return
	}

	printRegexNode(re, 0)
}

func printRegexNode(node *syntax.Regexp, indent int) {
	fmt.Printf("%*sop: %s\n", indent, "", node.Op)
	if node.Op == syntax.OpRepeat {
		fmt.Printf("%*smin,max: %d,%d\n", indent, "", node.Min, node.Max)
	}
	fmt.Printf("%*sflags: %v\n", indent, "", node.Flags)
	fmt.Printf("%*srune: %v\n", indent, "", string(node.Rune))
	fmt.Printf("%*ssub: %d\n", indent, "", len(node.Sub))

	for _, sub := range node.Sub {
		printRegexNode(sub, indent+2)
	}
}
