# 3d-segmentation-visualizer
visualize the outputs from the vision service's GetObjectPointClouds

<img width="400" alt="Screenshot 2023-09-21 at 1 27 34 PM" src="https://github.com/viam-labs/3d-segmentation-visualizer/assets/8298653/62399349-4f11-4c71-b4bc-b80b23d709bb">
<img width="400" alt="Screenshot 2023-09-21 at 1 28 46 PM" src="https://github.com/viam-labs/3d-segmentation-visualizer/assets/8298653/00da23e9-460e-4b79-bef3-6dc12ae8bb35">

## API

This is a viam camera component that fulfills the [Camera API](https://pkg.go.dev/go.viam.com/rdk@v0.9.1/components/camera).

## Config

There are two necessary attributes needed to use the visualizer:

- `vision_service_name`: the name of the vision service that implements GetObjectPointClouds
- `camera_name`: the name of the source camera that the GetObjectPointClouds will used on

```
    {
      "namespace": "rdk",
      "attributes": {
        "vision_service_name": "obstacle_service",
        "camera_name": "pointcloud_cam"
      },
      "depends_on": [],
      "name": "viz",
      "model": "viam-labs:camera:3d-segmentation-visualizer",
      "type": "camera"
    },
```
