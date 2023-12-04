package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Module struct remains the same as before

type Docker struct {
	cli *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{cli: cli}, nil
}

func (d *Docker) PullImage(module Module) error {
	ctx := context.Background()
	_, err := d.cli.ImagePull(ctx, module.ContainerImage, types.ImagePullOptions{})
	return err
}

func (d *Docker) RunContainer(module Module) error {
	ctx := context.Background()
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: module.ContainerImage,
		Cmd:   strings.Split(module.Command, " "),
		Env:   convertMapToSlice(module.EnvVariables),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: module.RestartPolicy,
		},
	}, nil, nil, module.ContainerName)
	if err != nil {
		return err
	}

	return d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}

func (d *Docker) ListContainers() ([]types.Container, error) {
	ctx := context.Background()
	return d.cli.ContainerList(ctx, types.ContainerListOptions{})
}

func convertMapToSlice(m map[string]string) []string {
	var env []string
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

// StopContainer stops a container given its name.
func (d *Docker) StopContainer(containerName string) error {
	ctx := context.Background()
	timeout := 10
	return d.cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer removes a container given its name.
func (d *Docker) RemoveContainer(containerName string) error {
	ctx := context.Background()
	return d.cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{})
}
