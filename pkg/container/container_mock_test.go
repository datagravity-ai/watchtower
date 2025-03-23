package container

import (
	types "github.com/docker/docker/api/types/image"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type MockContainerUpdate func(*dockerContainer.InspectResponse, *types.InspectResponse)

func MockContainer(updates ...MockContainerUpdate) *Container {
	containerInfo := dockerContainer.InspectResponse{
		ContainerJSONBase: &dockerContainer.ContainerJSONBase{
			ID:         "container_id",
			Image:      "image",
			Name:       "test-containrrr",
			HostConfig: &dockerContainer.HostConfig{},
		},
		Config: &dockerContainer.Config{
			Labels: map[string]string{},
		},
	}
	image := types.InspectResponse{
		ID:     "image_id",
		Config: &dockerContainer.Config{},
	}

	for _, update := range updates {
		update(&containerInfo, &image)
	}
	return NewContainer(&containerInfo, &image)
}

func WithPortBindings(portBindingSources ...string) MockContainerUpdate {
	return func(c *dockerContainer.InspectResponse, i *types.InspectResponse) {
		portBindings := nat.PortMap{}
		for _, pbs := range portBindingSources {
			portBindings[nat.Port(pbs)] = []nat.PortBinding{}
		}
		c.HostConfig.PortBindings = portBindings
	}
}

func WithImageName(name string) MockContainerUpdate {
	return func(c *dockerContainer.InspectResponse, i *types.InspectResponse) {
		c.Config.Image = name
		i.RepoTags = append(i.RepoTags, name)
	}
}

func WithLinks(links []string) MockContainerUpdate {
	return func(c *dockerContainer.InspectResponse, i *types.InspectResponse) {
		c.HostConfig.Links = links
	}
}

func WithLabels(labels map[string]string) MockContainerUpdate {
	return func(c *dockerContainer.InspectResponse, i *types.InspectResponse) {
		c.Config.Labels = labels
	}
}

func WithContainerState(state dockerContainer.State) MockContainerUpdate {
	return func(cnt *dockerContainer.InspectResponse, img *types.InspectResponse) {
		cnt.State = &state
	}
}

func WithHealthcheck(healthConfig dockerContainer.HealthConfig) MockContainerUpdate {
	return func(cnt *dockerContainer.InspectResponse, img *types.InspectResponse) {
		cnt.Config.Healthcheck = &healthConfig
	}
}

func WithImageHealthcheck(healthConfig dockerContainer.HealthConfig) MockContainerUpdate {
	return func(cnt *dockerContainer.InspectResponse, img *types.InspectResponse) {
		img.Config.Healthcheck = &healthConfig
	}
}
