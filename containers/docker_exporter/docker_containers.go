package main

import (
	"context"
	"encoding/json"
)

type ListContainersOptions struct {
	All     bool
	Size    bool
	Limit   int
	Since   string
	Before  string
	Filters map[string][]string
	Context context.Context
}

// APIPort is a type that represents a port mapping returned by the Docker API
type APIPort struct {
	PrivatePort int64  `json:"PrivatePort,omitempty" yaml:"PrivatePort,omitempty" toml:"PrivatePort,omitempty"`
	PublicPort  int64  `json:"PublicPort,omitempty" yaml:"PublicPort,omitempty" toml:"PublicPort,omitempty"`
	Type        string `json:"Type,omitempty" yaml:"Type,omitempty" toml:"Type,omitempty"`
	IP          string `json:"IP,omitempty" yaml:"IP,omitempty" toml:"IP,omitempty"`
}

// APIMount represents a mount point for a container.
type APIMount struct {
	Name        string `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Source      string `json:"Source,omitempty" yaml:"Source,omitempty" toml:"Source,omitempty"`
	Destination string `json:"Destination,omitempty" yaml:"Destination,omitempty" toml:"Destination,omitempty"`
	Driver      string `json:"Driver,omitempty" yaml:"Driver,omitempty" toml:"Driver,omitempty"`
	Mode        string `json:"Mode,omitempty" yaml:"Mode,omitempty" toml:"Mode,omitempty"`
	RW          bool   `json:"RW,omitempty" yaml:"RW,omitempty" toml:"RW,omitempty"`
	Propagation string `json:"Propagation,omitempty" yaml:"Propagation,omitempty" toml:"Propagation,omitempty"`
	Type        string `json:"Type,omitempty" yaml:"Type,omitempty" toml:"Type,omitempty"`
}

// APIContainers represents each container in the list returned by
// ListContainers.
type APIContainers struct {
	ID         string            `json:"Id" yaml:"Id" toml:"Id"`
	Image      string            `json:"Image,omitempty" yaml:"Image,omitempty" toml:"Image,omitempty"`
	Command    string            `json:"Command,omitempty" yaml:"Command,omitempty" toml:"Command,omitempty"`
	Created    int64             `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
	State      string            `json:"State,omitempty" yaml:"State,omitempty" toml:"State,omitempty"`
	Status     string            `json:"Status,omitempty" yaml:"Status,omitempty" toml:"Status,omitempty"`
	Ports      []APIPort         `json:"Ports,omitempty" yaml:"Ports,omitempty" toml:"Ports,omitempty"`
	SizeRw     int64             `json:"SizeRw,omitempty" yaml:"SizeRw,omitempty" toml:"SizeRw,omitempty"`
	SizeRootFs int64             `json:"SizeRootFs,omitempty" yaml:"SizeRootFs,omitempty" toml:"SizeRootFs,omitempty"`
	Names      []string          `json:"Names,omitempty" yaml:"Names,omitempty" toml:"Names,omitempty"`
	Labels     map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`
	Networks   NetworkList       `json:"NetworkSettings,omitempty" yaml:"NetworkSettings,omitempty" toml:"NetworkSettings,omitempty"`
	Mounts     []APIMount        `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`
}

// PortBinding represents the host/container port mapping as returned in the
// `docker inspect` json
type PortBinding struct {
	HostIP   string `json:"HostIp,omitempty" yaml:"HostIp,omitempty" toml:"HostIp,omitempty"`
	HostPort string `json:"HostPort,omitempty" yaml:"HostPort,omitempty" toml:"HostPort,omitempty"`
}

// PortMapping represents a deprecated field in the `docker inspect` output,
// and its value as found in NetworkSettings should always be nil
type PortMapping map[string]string

// ContainerNetwork represents the networking settings of a container per network.
type ContainerNetwork struct {
	Aliases             []string `json:"Aliases,omitempty" yaml:"Aliases,omitempty" toml:"Aliases,omitempty"`
	MacAddress          string   `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
	GlobalIPv6PrefixLen int      `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
	GlobalIPv6Address   string   `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
	IPv6Gateway         string   `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
	IPPrefixLen         int      `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
	IPAddress           string   `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
	Gateway             string   `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
	EndpointID          string   `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
	NetworkID           string   `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
}

// NetworkSettings contains network-related information about a container
type NetworkSettings struct {
	Networks               map[string]ContainerNetwork `json:"Networks,omitempty" yaml:"Networks,omitempty" toml:"Networks,omitempty"`
	IPAddress              string                      `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
	IPPrefixLen            int                         `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
	MacAddress             string                      `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
	Gateway                string                      `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
	Bridge                 string                      `json:"Bridge,omitempty" yaml:"Bridge,omitempty" toml:"Bridge,omitempty"`
	PortMapping            map[string]PortMapping      `json:"PortMapping,omitempty" yaml:"PortMapping,omitempty" toml:"PortMapping,omitempty"`
	Ports                  map[APIPort][]PortBinding   `json:"Ports,omitempty" yaml:"Ports,omitempty" toml:"Ports,omitempty"`
	NetworkID              string                      `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
	EndpointID             string                      `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
	SandboxKey             string                      `json:"SandboxKey,omitempty" yaml:"SandboxKey,omitempty" toml:"SandboxKey,omitempty"`
	GlobalIPv6Address      string                      `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
	GlobalIPv6PrefixLen    int                         `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
	IPv6Gateway            string                      `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
	LinkLocalIPv6Address   string                      `json:"LinkLocalIPv6Address,omitempty" yaml:"LinkLocalIPv6Address,omitempty" toml:"LinkLocalIPv6Address,omitempty"`
	LinkLocalIPv6PrefixLen int                         `json:"LinkLocalIPv6PrefixLen,omitempty" yaml:"LinkLocalIPv6PrefixLen,omitempty" toml:"LinkLocalIPv6PrefixLen,omitempty"`
	SecondaryIPAddresses   []string                    `json:"SecondaryIPAddresses,omitempty" yaml:"SecondaryIPAddresses,omitempty" toml:"SecondaryIPAddresses,omitempty"`
	SecondaryIPv6Addresses []string                    `json:"SecondaryIPv6Addresses,omitempty" yaml:"SecondaryIPv6Addresses,omitempty" toml:"SecondaryIPv6Addresses,omitempty"`
}

// NetworkList encapsulates a map of networks, as returned by the Docker API in
// ListContainers.
type NetworkList struct {
	Networks map[string]ContainerNetwork `json:"Networks" yaml:"Networks,omitempty" toml:"Networks,omitempty"`
}

func (c *Docker) ListContainers() ([]APIContainers, error) {
	var containers []APIContainers
	resp, err := c.client.Get("http://unix/v1.24/containers/json?all=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, err
	}
	return containers, nil
	//io.Copy(os.Stdout, response.Body)
}
