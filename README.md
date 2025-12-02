# `3d-segmentation-visualizer` modular component

This module implements the [Viam Camera API](https://docs.viam.com/dev/reference/apis/components/camera/) in a `viam-labs:camera:3d-segmentation-visualizer` model.

With this model, you can visualize the outputs from a vision service's `GetObjectPointClouds` method by overlaying colored segmentation data onto point clouds from a camera.

<img width="400" alt="Screenshot 2023-09-21 at 1 27 34 PM" src="https://github.com/viam-labs/3d-segmentation-visualizer/assets/8298653/62399349-4f11-4c71-b4bc-b80b23d709bb">
<img width="400" alt="Screenshot 2023-09-21 at 1 28 46 PM" src="https://github.com/viam-labs/3d-segmentation-visualizer/assets/8298653/00da23e9-460e-4b79-bef3-6dc12ae8bb35">

Navigate to the **CONFIGURE** tab of your machine's page.

Click the **+** button, select **Component or service**, then select the `camera / 3d-segmentation-visualizer` model provided by the [`3d-segmentation-visualizer` module](https://app.viam.com/module/viam-labs/3d-segmentation-visualizer).

Click **Add module**, enter a name for your camera, and click **Create**.

## Configure your `3d-segmentation-visualizer` camera

On the new component panel, copy and paste the following attribute template into your camera's **Attributes** box:

```json
{
  "camera_name": "<string>",
  "vision_service_name": "<string>"
}
```

### Attributes

The following attributes are available for `viam-labs:camera:3d-segmentation-visualizer` cameras:

| Name | Type   | Inclusion | Description |
| ---- | ------ | --------- | ----------- |
| `camera_name`         | string | Required  | The name of the source camera that `GetObjectPointClouds` uses as input |
| `vision_service_name` | string | Required  | The name of the vision service that implements `GetObjectPointClouds` |

### Example Configuration

```json
{
  "vision_service_name": "obstacle_service",
  "camera_name": "pointcloud_cam"
}
```