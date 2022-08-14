/*
Copyright © 2022 Mohamed Hammad Youssef mmhy2003@hotmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mmhy2003/underscoreai/config"
	"github.com/spf13/cobra"
)

const HFAPIENDPOINT = "https://api-inference.huggingface.co/models/bigscience/bloom"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "_",
	Short: "A command-line tool to help shell users describe their command and get it via AI",
	Long:  `A command-line tool to help shell users describe their command and get it via AI`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			prompt := strings.Join(args, " ")
			result := GetResult(prompt)
			fmt.Println(ProcessResult(result))
			// fmt.Println(result.GeneratedText)
		}
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.underscoreai.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type HFResult struct {
	GeneratedText string `json:"generated_text,omitempty"`
}

type OptionsStruct struct {
	UseCache     bool `json:"use_cache,omitempty"`
	WaitForModel bool `json:"wait_for_model,omitempty"`
}

type ParametersStruct struct {
	DoSample         bool    `json:"do_sample,omitempty"`
	ReturnFullText   bool    `json:"return_full_text,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	MaxLength        int     `json:"max_length,omitempty"`
	MaxNewTokens     int     `json:"max_new_tokens,omitempty"`
	StoppingCriteria string  `json:"stopping_criteria,omitempty"`
}

type RequestStruct struct {
	Inputs     string           `json:"inputs,omitempty"`
	Parameters ParametersStruct `json:"parameters,omitempty"`
	Options    OptionsStruct    `json:"options,omitempty"`
}

func LoadPromptContext() string {
	// read pre-set prompts from file
	promptContext, err := ioutil.ReadFile(config.Config.PromptContextPath)
	if err != nil {
		log.Fatal(err)
	}

	promptContextStr := string(promptContext)

	return promptContextStr
}

func GetResult(prompt string) HFResult {
	// load prompt context
	var promptContext = LoadPromptContext()
	newPrompt := "...\nP: " + prompt + "\nA:"
	finalPrompt := promptContext + newPrompt

	data := RequestStruct{
		Inputs: finalPrompt,
		Parameters: ParametersStruct{
			Temperature:      0.3,
			MaxLength:        100,
			MaxNewTokens:     64,
			StoppingCriteria: "...",
		},
		Options: OptionsStruct{
			UseCache:     true,
			WaitForModel: true,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		HFAPIENDPOINT,
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+config.Config.HFAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results []HFResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	result := results[0]

	return result
}

func ProcessResult(result HFResult) string {
	// split result.GeneratedText at ...
	var resultStr = strings.Split(result.GeneratedText, "...")

	// get last element of resultStr
	var lastElement = resultStr[len(resultStr)-2]

	// remove new-line character from beginning and end of lastElement
	lastElement = strings.Trim(lastElement, "\n")

	// split lastElement at \n
	var lastElementPortion = strings.Split(lastElement, "\n")[1]

	// remove prompt from beginning of result
	lastElementPortion = strings.TrimPrefix(lastElementPortion, "A: ")

	return lastElementPortion
}
