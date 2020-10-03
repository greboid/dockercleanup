package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/kouhin/envflag"
)

var (
	duration           = flag.Duration("duration", 0, "Number of seconds between executions")
	removeUnusedImages = flag.Bool("allimages", false, "Remove all unused images, not just dangling ones")
	removeContainer    = flag.Bool("containers", true, "Remove containers that are not running")
)

func main() {
	err := envflag.Parse()
	if err != nil {
		log.Fatalf("unable to parse config location: %s", err.Error())
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("unable to create docker client: %s", err.Error())
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
		log.Fatalf("unable to clean perform cleanup: %s", err.Error())
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
		log.Printf("Removed %s",
			humanReadableSize(containerReport.SpaceReclaimed),
		)
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
		log.Printf("Removed %s",
			humanReadableSize(imageReport.SpaceReclaimed),
		)
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
