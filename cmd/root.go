package cmd

import (
	"fmt"
	"time"

	"github.com/rs/xlog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cobra-example/pkg/config"
	"cobra-example/pkg/server"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cobra-example",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			logLevelStr := viper.GetString("LOG_LEVEL")
			logLevel, err := xlog.LevelFromString(logLevelStr)
			if err != nil {
				return fmt.Errorf("Invalid log level: '%s'. Error: %v\n", logLevelStr, err)
			}

			cfg := &config.Config{
				LogLevel:    logLevel,
				HttpPort:    viper.GetUint("PORT"),
				HttpTimeout: viper.New().GetDuration("TIMEOUT"),
				MongoURI:    viper.GetString("DB_URI"),
			}

			s := server.New(cfg)

			return s.Run()
		},
	}

	rootCmd.AddCommand(newVersion())

	viper.AutomaticEnv()
	// viper.SetEnvPrefix("APP")

	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set the logging level ('debug', 'info', 'warn', 'error', 'fatal') ENV: LOG_LEVEL")
	viper.BindPFlag("LOG_LEVEL", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().String("db-uri", "mongodb://localhost:27017", "The URI to connect to DB ENV: DB_URI")
	viper.BindPFlag("DB_URI", rootCmd.PersistentFlags().Lookup("db-uri"))

	rootCmd.PersistentFlags().UintP("port", "p", 8080, "HTTP server port ENV: PORT")
	viper.BindPFlag("PORT", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.PersistentFlags().DurationP("http-timeout", "t", 5*time.Second, "Maximum time to allow one HTTP request ENV: TIMEOUT")
	viper.BindPFlag("TIMEOUT", rootCmd.PersistentFlags().Lookup("port"))

	return rootCmd
}
