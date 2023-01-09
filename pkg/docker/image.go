/*
credit: https://github.com/containrrr/watchtower
*/
package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Image(ctx context.Context, imageName string) (*types.ImageInspect, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	imageInfo, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return nil, err
	}

	return &imageInfo, nil
}

func PullImage(ctx context.Context, imageName string, opts types.ImagePullOptions, handleOut func(io.ReadCloser) error) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, imageName, opts)
	if err != nil {
		return err
	}
	defer out.Close()

	if handleOut != nil {
		if err := handleOut(out); err != nil {
			return err
		}
	} else {
		if _, err := io.ReadAll(out); err != nil {
			return err
		}
	}

	return nil
}

func PullNewImage(ctx context.Context, imageName string, handleOut func(io.ReadCloser) error) error {
	if strings.HasPrefix(imageName, "sha256:") {
		return fmt.Errorf("container uses a pinned image, and cannot be updated")
	}

	opts, err := GetPullOptions(imageName)
	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	imageInfo, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return err
	}

	if match, err := CompareDigest(imageName, imageInfo.RepoDigests, opts.RegistryAuth); err != nil {
		// do nothing
	} else if match {
		return nil
	}

	return PullImage(ctx, imageName, opts, handleOut)
}

func HasNewImage(ctx context.Context, imageName string, currentImageID string) (bool, string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, currentImageID, err
	}
	defer cli.Close()

	newImageInfo, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return false, currentImageID, err
	}

	newImageID := newImageInfo.ID
	if newImageID == currentImageID {
		return false, currentImageID, nil
	}

	return true, newImageID, nil
}
