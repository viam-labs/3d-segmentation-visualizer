package segmentationvisualizer

import (
	"context"
	"fmt"
	"image/color"

	"github.com/golang/geo/r3"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/edaniels/golog"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
)

// ModelName is the name of the model
const ModelName = "3d-segmentation-visualizer"

var (
	// Model is the full resource name of the model
	Model            = resource.NewModel("viam-labs", "camera", ModelName)
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(
		camera.API,
		Model,
		resource.Registration[camera.Camera, *Config]{Constructor: newVisualizer},
	)
}

// Config specifies which camera and which service should be used to do the overlay
type Config struct {
	CameraName  string `json:"camera_name"`
	ServiceName string `json:"vision_service_name"`
}

// Validate will ensure that both the underlying camera and service are present
func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.CameraName == "" {
		return nil, fmt.Errorf(`expected "camera_name" attribute for %s %q`, ModelName, path)
	}
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf(`expected "vision_service_name" attribute for %s %q`, ModelName, path)
	}

	return []string{cfg.CameraName, cfg.ServiceName}, nil
}

type visualizer struct {
	resource.Named
	camera.VideoSource
	cameraName string
	service    vision.Service
	logger     golog.Logger
}

func newVisualizer(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger golog.Logger,
) (camera.Camera, error) {
	v := &visualizer{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	if err := v.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return camera.FromVideoSource(conf.ResourceName(), v), nil
}

func (v *visualizer) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	v.service = nil
	v.cameraName = ""
	cfg, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return errors.Errorf("Could not assert proper config for %s", ModelName)
	}
	// get the source camera
	v.cameraName = cfg.CameraName
	cam, err := camera.FromDependencies(deps, cfg.CameraName)
	if err != nil {
		return errors.Wrapf(err, "unable to get camera %v for %s", cfg.CameraName, ModelName)
	}
	vs, ok := cam.(camera.VideoSource)
	if !ok {
		return errors.Wrapf(err, "camera %v is not a video source for %s", cfg.CameraName, ModelName)
	}
	v.VideoSource = vs
	// get the source service
	v.service, err = vision.FromDependencies(deps, cfg.ServiceName)
	if err != nil {
		return errors.Wrapf(err, "unable to get vision service %v for %s", cfg.ServiceName, ModelName)
	}
	return nil
}

// NextPointCloud function calls a segmenter service on the underlying camera and returns a pointcloud.
func (v *visualizer) NextPointCloud(ctx context.Context) (pointcloud.PointCloud, error) {
	// get the service
	clouds, err := v.service.GetObjectPointClouds(ctx, v.cameraName, map[string]interface{}{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not get point clouds from the vision service")
	}
	if clouds == nil {
		return pointcloud.New(), nil
	}

	// merge pointclouds with a color overlay
	merged := pointcloud.New()
	palette := colorful.FastWarmPalette(len(clouds))
	for i, cluster := range clouds {
		col, ok := color.NRGBAModel.Convert(palette[i]).(color.NRGBA)
		if !ok {
			return nil, errors.Errorf("could not convert color %T to NRGBA", palette[i])
		}
		colR, colG, colB := int(col.R), int(col.G), int(col.B)
		cluster.Iterate(0, 0, func(v r3.Vector, d pointcloud.Data) bool {
			r, g, b := d.RGB255()
			mixR := uint8((int(r) + colR) / 2)
			mixG := uint8((int(g) + colG) / 2)
			mixB := uint8((int(b) + colB) / 2)
			err = merged.Set(v, pointcloud.NewColoredData(color.NRGBA{R: mixR, G: mixG, B: mixB}))
			return err == nil
		})
		if err != nil {
			return nil, err
		}
	}
	return merged, nil
}

// DoCommand simply echos whatever was sent.
func (v *visualizer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return cmd, nil
}

// Close closes the underlying stream.
func (v *visualizer) Close(ctx context.Context) error {
	err1 := v.service.Close(ctx)
	err2 := v.VideoSource.Close(ctx)
	return multierr.Combine(err1, err2)
}
