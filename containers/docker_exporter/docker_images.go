package main

import (
	"encoding/json"
)

// APIImages represent an image returned in the ListImages call.
type APIImages struct {
	ID          string            `json:"Id" yaml:"Id" toml:"Id"`
	RepoTags    []string          `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty" toml:"RepoTags,omitempty"`
	Created     int64             `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
	Size        int64             `json:"Size,omitempty" yaml:"Size,omitempty" toml:"Size,omitempty"`
	VirtualSize int64             `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty" toml:"VirtualSize,omitempty"`
	ParentID    string            `json:"ParentId,omitempty" yaml:"ParentId,omitempty" toml:"ParentId,omitempty"`
	RepoDigests []string          `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty" toml:"RepoDigests,omitempty"`
	Labels      map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`
}

func (c *Docker) ListImages() ([]APIImages, error) {
	var images []APIImages
	resp, err := c.client.Get("http://unix/v1.24/images/json?all=0")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return nil, err
	}
	return images, nil
	//io.Copy(os.Stdout, response.Body)
}
