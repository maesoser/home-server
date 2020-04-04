package tcping

import (
	"context"
	"fmt"
	shaker "github.com/tevino/tcp-shaker"
	"log"
	"math"
	"net"
	"time"
)

// NewTCPinger returns a new Pinger struct pointer
func NewTCPinger(addr string, port string) (*Pinger, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	return &Pinger{
		ipaddr:      ipaddr,
		addr:        addr,
		Interval:    time.Second,
		Timeout:     time.Second * 3,
		Count:       -1,
		Port:        port,
		PacketsSent: 0,
		PacketsRecv: 0,
	}, nil
}

// Pinger represents ICMP packet sender/receiver
type Pinger struct {
	Interval    time.Duration
	Timeout     time.Duration
	Port        string
	Count       int
	PacketsSent int
	PacketsRecv int
	rtts        []time.Duration
	Source      string
	ipaddr      *net.IPAddr
	addr        string
}

// Statistics represent the stats of a currently running or finished
// pinger operation.
type Statistics struct {
	PacketsRecv int
	PacketsSent int
	PacketLoss  float64
	IPAddr      *net.IPAddr
	Addr        string
	Rtts        []time.Duration
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	StdDevRtt   time.Duration
}

// SetIPAddr sets the ip address of the target host.
func (p *Pinger) SetIPAddr(ipaddr *net.IPAddr) {
	p.ipaddr = ipaddr
	p.addr = ipaddr.String()
}

// IPAddr returns the ip address of the target host.
func (p *Pinger) IPAddr() *net.IPAddr {
	return p.ipaddr
}

// SetAddr resolves and sets the ip address of the target host, addr can be a
// DNS name like "www.google.com" or IP like "127.0.0.1".
func (p *Pinger) SetAddr(addr string) error {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return err
	}
	p.SetIPAddr(ipaddr)
	p.addr = addr
	return nil
}

// Addr returns the string ip address of the target host.
func (p *Pinger) Addr() string {
	return p.addr
}

func (p *Pinger) Run() {
	c := shaker.NewChecker()

	ctx, stopChecker := context.WithCancel(context.Background())
	defer stopChecker()
	go func() {
		if err := c.CheckingLoop(ctx); err != nil {
			fmt.Println("checking loop stopped due to fatal error: ", err)
		}
	}()

	<-c.WaitReady()
	for i := 0; i < p.Count; i++ {
		p.PacketsSent++
		start := time.Now()
		err := c.CheckAddr(p.addr+":"+p.Port, p.Timeout)
		elapsed := time.Since(start)
		p.rtts = append(p.rtts, elapsed)
		if err != nil {
			log.Println(p.addr+":"+p.Port, err)
		} else {
			p.PacketsRecv++
		}
		time.Sleep(p.Interval)
	}

}

// Statistics returns the statistics of the pinger. This can be run while the
// pinger is running or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func (p *Pinger) Statistics() *Statistics {
	loss := float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.rtts) > 0 {
		min = p.rtts[0]
		max = p.rtts[0]
	}
	for _, rtt := range p.rtts {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsSent: p.PacketsSent,
		PacketsRecv: p.PacketsRecv,
		PacketLoss:  loss,
		Rtts:        p.rtts,
		Addr:        p.addr,
		IPAddr:      p.ipaddr,
		MaxRtt:      max,
		MinRtt:      min,
	}
	if len(p.rtts) > 0 {
		s.AvgRtt = total / time.Duration(len(p.rtts))
		var sumsquares time.Duration
		for _, rtt := range p.rtts {
			sumsquares += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		s.StdDevRtt = time.Duration(math.Sqrt(
			float64(sumsquares / time.Duration(len(p.rtts)))))
	}
	return &s
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}
