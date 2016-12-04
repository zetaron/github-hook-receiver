package cmd

// Copyright ©2016 Fabian Stegemann
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/adjust/rmq"
	"github.com/google/go-github/github"
	"github.com/rjz/githubhook"
)

var cfgFile string
var version = "1.0.0"

// RootCmd : This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "github-hook-receiver",
	Short: "Collects hooks from GitHub.",
	Long:  `Receives valid hooks from GitHub and published the payload to the queue to be processed by the queue-workers.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetDefault("redis.url", "redis:6379")
		viper.SetDefault("redis.database", 1)
		viper.SetDefault("host", ":80")

		log.Infof("Starting github-hook-receiver version: %s", version)

		if !viper.IsSet("github.secret") {
			log.Fatal("No GitHub secret defined.")
			return
		}

		connection := rmq.OpenConnection(
			"github-receive",
			"tcp",
			viper.GetString("redis.url"),
			viper.GetInt("redis.database"),
		)
		log.Info("Connected to redis")

		queue := connection.OpenQueue("deployment_events")
		log.Info("Opened queue")

		http.HandleFunc("/deployment", func(w http.ResponseWriter, r *http.Request) {
			hook, err := githubhook.Parse([]byte(viper.GetString("github.secret")), r)
			if err != nil {
				log.WithFields(log.Fields{
					"error":   err,
					"request": r,
				}).Error("Could not parse GitHub hook from request.")

				w.WriteHeader(http.StatusBadRequest)

				return
			}

			event := github.DeploymentEvent{}
			if err := json.Unmarshal(hook.Payload, &event); err != nil {
				log.WithFields(log.Fields{
					"error":   err,
					"payload": hook.Payload,
				}).Error("Could not parse DeploymentEvent from payload.")

				w.WriteHeader(http.StatusBadRequest)

				return
			}

			queue.PublishBytes(hook.Payload)
		})

		log.Fatal(http.ListenAndServe(viper.GetString("host"), nil))
	},
}

//Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings
	// Cobra supports Persistent Flags which if defined here will be global for your application

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/cutter/github-receive.yaml)")

	// Cobra also supports local flags which will only run when this action is called directly
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// Read in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// allow for nested environment variables
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigName("github-receive") // name of config file (without extension)
	viper.AddConfigPath("/etc/cutter")    // adding home directory as first search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
