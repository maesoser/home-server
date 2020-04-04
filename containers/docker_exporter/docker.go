package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type ContainerStats struct {
	info  APIContainers
	stats Stats
}

type Docker struct {
	client *http.Client
}

func (c *Docker) Start(endpoint string) {
	c.client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", endpoint)
			},
		},
	}
	log.Println("Connected to docker")
}

func (c *Docker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	States := map[string]int{
		"created":    1,
		"restarting": 2,
		"running":    3,
		"paused":     4,
		"exited":     5,
		"dead":       6,
	}

	imgs, err := c.ListImages()
	if err != nil {
		log.Printf("Error obtaining image list: %s\n", err)
	}
	fmt.Fprintf(w, "docker_img_count %d\n", len(imgs))
	for _, img := range imgs {
		fmt.Fprintf(w, "docker_img_size_bytes{id=\"%s\", repotag=\"%s\"} %d\n", img.ID, img.RepoTags[0], img.Size)
		fmt.Fprintf(w, "docker_vimg_size_bytes{id=\"%s\", repotag=\"%s\"} %d\n", img.ID, img.RepoTags[0], img.VirtualSize)
	}

	containers, err := c.ListContainers()
	if err != nil {
		log.Printf("Error obtaining container list: %s\n", err)
	}
	fmt.Fprintf(w, "docker_container_count %d\n", len(containers))
	start := time.Now()
	for _, container := range containers {
		stats, err := c.GetStats(container.ID)
		if err != nil {
			log.Printf("Error obtaining container stats: %v\n", err)
		}
		ioread, iowrite := stats.ReadIO()
		nettx, netrx := stats.ReadNet()
		usage, limit, percent := stats.ReadMem()
		total, system, ncpus, pids := stats.ReadCPU()
		fmt.Fprintf(w, "docker_container_status{id=\"%s\", name=\"%s\"} %d\n", container.Names[0], container.ID, States[container.State])
		fmt.Fprintf(w, "docker_container_created{id=\"%s\", name=\"%s\"} %d\n", container.Names[0], container.ID, container.Created)
		fmt.Fprintf(w, "docker_container_net_tx_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], nettx)
		fmt.Fprintf(w, "docker_container_net_rx_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], netrx)
		fmt.Fprintf(w, "docker_container_mem_limit_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], limit)
		fmt.Fprintf(w, "docker_container_mem_usage_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], usage)
		fmt.Fprintf(w, "docker_container_mem_usage_pcnt{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], percent)
		fmt.Fprintf(w, "docker_container_io_read_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], ioread)
		fmt.Fprintf(w, "docker_container_io_write_bytes{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], iowrite)
		fmt.Fprintf(w, "docker_container_pids_count{id=\"%s\", name=\"%s\"} %d\n", container.ID, container.Names[0], pids)
		fmt.Fprintf(w, "docker_container_cpu_total{id=\"%s\", name=\"%s\"} %f\n", container.ID, container.Names[0], total)
		fmt.Fprintf(w, "docker_container_cpu_system{id=\"%s\", name=\"%s\"} %f\n", container.ID, container.Names[0], system)
		fmt.Fprintf(w, "docker_container_cpu_count{id=\"%s\", name=\"%s\"} %f\n", container.ID, container.Names[0], ncpus)
	}
	log.Printf("Collected info for %d containers (%.2f secs)\n", len(containers), time.Since(start).Seconds())
}

func (stats *Stats) ReadCPU() (float64, float64, float64, int) {
	ncpus := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)
	Pids := int(stats.PidsStats.Current)
	return total, system, ncpus, Pids
}

func (stats *Stats) ReadMem() (int64, int64, int) {
	usage := int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
	limit := int64(stats.MemoryStats.Limit)
	percent := percent(float64(usage), float64(limit))
	return usage, limit, percent
}

func (stats *Stats) ReadNet() (int64, int64) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	return rx, tx
}

func (stats *Stats) ReadIO() (int64, int64) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IOServiceBytesRecursive {
		if blk.Op == "Read" {
			read = int64(blk.Value)
		}
		if blk.Op == "Write" {
			write = int64(blk.Value)
		}
	}
	return read, write
}
