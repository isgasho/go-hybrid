package runtime

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/XSAM/go-hybrid/cmdutil"
	"github.com/XSAM/go-hybrid/environment"
	"github.com/XSAM/go-hybrid/errorw"
	"github.com/XSAM/go-hybrid/log"
	"github.com/XSAM/go-hybrid/log/zapfield"
	"github.com/XSAM/go-hybrid/metadata"
)

var flag Flag

func Start() {
	cmd := rootCmd()
	cmd.AddCommand(cmdutil.VersionCmd())
	cmd.AddCommand(flagCmd())
	cmd.AddCommand(logCmd())

	cmd.Execute()
}

func rootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:  metadata.AppName(),
		Long: "Example for go-hybrid.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flag.Environment.DevelopmentMode {
				environment.DevelopmentMode()
			} else {
				environment.ProductionMode()
			}

			if flag.Environment.JSONLogStyle {
				environment.JSONLogStyle()
			} else {
				environment.TextLogStyle()
			}
		},
	}

	// Default value
	flag = Flag{
		Number: 42,
	}
	err := cmdutil.ResolveFlagVariable(&cmd, &flag)
	if err != nil {
		log.BgLogger().Fatal("resolve flag", zapfield.Error(err))
	}

	return &cmd
}

func flagCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "flag",
		Short: "Print flag parse result",
		Run: func(cmd *cobra.Command, args []string) {
			log.BgLogger().Info("parse result", zap.Any("flag", flag))
		},
	}
	return &cmd
}

func logCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "log",
		Short: "Print logging sample",
		Run: func(cmd *cobra.Command, args []string) {
			// Background logger
			log.BgLogger().Info("-- background logger --")
			log.BgLogger().Debug("debug")
			log.BgLogger().Info("info")
			log.BgLogger().Warn("warn")
			log.BgLogger().Error("error")

			// Context logger
			log.BgLogger().Info("-- context logger --")
			ctx := log.WithKeyValue(context.Background(), "hello", "world")
			zapFields := []zap.Field{zap.Duration("duration", 42*time.Second), zap.Time("date", time.Now())}
			log.Logger(ctx).Debug("debug", zapFields...)
			log.Logger(ctx).Info("info", zapFields...)
			log.Logger(ctx).Warn("warn", zapFields...)
			log.Logger(ctx).Error("error", zapFields...)

			// Print errorw
			log.BgLogger().Info("-- print errorw --")
			err := errorw.NewMessage("error cause").
				WithField("key", "value").
				WithWrap("wrap1").
				WithWrap("wrap2")
			log.BgLogger().Debug("debug", zapfield.Error(err))
			log.BgLogger().Info("info", zapfield.Error(err))
			log.BgLogger().Warn("warn", zapfield.Error(err))
			log.BgLogger().Error("error", zapfield.Error(err))
		},
	}
	return &cmd
}
