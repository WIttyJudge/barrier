package action

import (
	"barrier/internal/action/exit"
	"barrier/internal/hostsfile"

	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func (a *Action) Restore(ctx *cli.Context) error {
	log.Info().Msg("restoring hosts file from backup..")

	hosts, err := hostsfile.New()
	if err != nil {
		return exit.Error(exit.Hostsfile, err, "failed to process hosts file")
	}

	if err := hosts.Restore(); err != nil {
		return exit.Error(exit.Hostsfile, err, "failed to restore hosts file")
	}

	log.Info().Msg("restoring operation is successful")

	return nil
}
