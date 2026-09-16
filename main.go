// Package main is a module which serves the 3d segmentation visualizer
package main

import (
	"context"

	"github.com/viam-labs/3d-segmentation-visualizer/segmentationvisualizer"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("3d-segmentation-visualizer"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	vizMod, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	err = vizMod.AddModelFromRegistry(ctx, camera.API, segmentationvisualizer.Model)
	if err != nil {
		return err
	}

	err = vizMod.Start(ctx)
	defer vizMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
