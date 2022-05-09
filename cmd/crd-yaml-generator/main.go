package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	// nolint:typecheck
	"github.com/nandajavarma/crd-yaml-generator/pkg/generator"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var logLevel string
var cfgFile string
var crdDefFile string

var rootCmd = &cobra.Command{ // nolint:gochecknoglobals
	PersistentPreRunE: configLogger,
	Use:               "cacophony-dht",
	Short:             "Let's make some noise",
	SilenceUsage:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read the yaml file
		// marshal it to CRD
		yamlfile, err := ioutil.ReadFile(crdDefFile)
		if err != nil {
			return err
		}

		extv1.AddToScheme(scheme.Scheme)
		crdData := extv1.CustomResourceDefinition{}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		stream, _ := ioutil.ReadFile(crdDefFile)
		obk, _, _ := decode(stream, nil, nil)
		fmt.Print(obk)

		err = yaml.Unmarshal(yamlfile, &crdData)
		if err != nil {
			return err
		}

		result, err := generator.GenerateYaml(crdData)
		if err != nil {
			return err
		}

		return ioutil.WriteFile("test.yaml", result, 0644)
	},
}

func configLogger(cmd *cobra.Command, args []string) error {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithField("log-level", logLevel).Fatal("incorrect log level")

		return fmt.Errorf("incorrect log level")
	}

	log.SetLevel(lvl)
	log.WithField("log-level", logLevel).Debug("log level configured")

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))

			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", "CDHT", envVarSuffix))
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)

			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
		}
	})
}

func init() {
	v := readConfigFile()

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l",
		"info", "set log level verbosity (options: debug, info, error, warning)")

	rootCmd.PersistentFlags().StringVarP(&crdDefFile, "crd-def-yaml", "c", "", `CRD definition file to generate YAML from`)
	rootCmd.MarkFlagRequired("crd-def-yaml") // we require the port to bind the service to

	bindFlags(rootCmd, v)
}

func readConfigFile() *viper.Viper {
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name (without extension).
		v.AddConfigPath(home)
		v.SetConfigType("yaml")
		v.SetConfigName("crd-yaml-generator")
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Info(err)
		}
	}

	v.SetEnvPrefix("CDHT")

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	return v
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("error in the cli. Exiting")
		os.Exit(1)
	}
}
