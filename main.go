package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/slogflags"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	duration           = flag.Duration("duration", 0, "Number of seconds between executions")
	removeUnusedImages = flag.Bool("allimages", false, "Remove all unused images, not just dangling ones")
	removeContainer    = flag.Bool("containers", true, "Remove containers that are not running")
)

func main() {
	envflag.Parse()
	slogflags.Logger(slogflags.WithSetDefault(true))
	slog.Info("Starting", "duration", duration.String(), "allimages", *removeUnusedImages, "containers", *removeContainer)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		slog.Error("unable to create docker client", "error", err)
		os.Exit(1)
	}
	if *duration < time.Minute {
		cleanAndOutput(cli)
		return
	}
	for {
		cleanAndOutput(cli)
		time.Sleep(*duration)
	}
}

func cleanAndOutput(cli *client.Client) {
	err := cleanAll(cli)
	if err != nil {
		slog.Error("unable to perform cleanup", "error", err)
		os.Exit(1)
	}
}

func cleanAll(cli *client.Client) error {
	if *removeContainer {
		err := cleanContainers(cli)
		if err != nil {
			return err
		}
	}
	err := cleanImages(cli)
	if err != nil {
		return err
	}
	return nil
}

func cleanContainers(cli *client.Client) error {
	pruneFilter := filters.NewArgs()
	containerReport, err := cli.ContainersPrune(
		context.Background(),
		pruneFilter,
	)
	if err != nil {
		return err
	}
	if len(containerReport.ContainersDeleted) > 0 {
		slog.Info("Removed containers", "space_reclaimed", humanReadableSize(containerReport.SpaceReclaimed))
	}
	return nil
}

func cleanImages(cli *client.Client) error {
	pruneFilter := filters.NewArgs()
	if *removeUnusedImages {
		pruneFilter.Add("dangling", "false")
	} else {
		pruneFilter.Add("dangling", "true")
	}
	imageReport, err := cli.ImagesPrune(
		context.Background(),
		pruneFilter,
	)
	if err != nil {
		return err
	}
	if len(imageReport.ImagesDeleted) > 0 {
		slog.Info("Removed images", "space_reclaimed", humanReadableSize(imageReport.SpaceReclaimed))
	}
	return nil
}

func humanReadableSize(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
