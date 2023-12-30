package helpers

import (
	"fmt"
	"net"
	"time"

	logger "github.com/dannonb/go-network-monitor/logger"

	probing "github.com/prometheus-community/pro-bing"
)

const (
	numRetries = 3
	timeout    = time.Second * 2
	Interval   = time.Second * 10
)

type PingResponse struct {
	AvgTime time.Duration
	MinTime time.Duration
	MaxTime time.Duration
	Jitter  int64

	PacketsSent     int
	PacketsReceived int

	ResponseTimes []time.Duration
	Host          string
}

// func Ping(hosts []string, packets int) []PingResponse {
// 	fmt.Println("STARTING PING")
// 	defer func(start time.Time) {
// 		fmt.Printf("FROM PING: time=%v", time.Since(start))
// 	}(time.Now())

// 	var responseList []PingResponse

// 	for _, h := range hosts {
// 		responseList = append(responseList, SinglePing(h, packets))
// 	}

// 	return responseList
// }

// func SinglePing(host string, packets int) PingResponse {
// 	fmt.Println("SINGLE PING: ", host)
// 	pinger, err := ping.NewPinger(host)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	pinger.Count = packets
// 	pinger.SetPrivileged(true)
// 	pinger.Run()

// 	return parseToPingResponse(pinger.Statistics(), host)
// }

func parseToPingResponse(stats *probing.Statistics, host string) PingResponse {
	return PingResponse{
		AvgTime:         stats.AvgRtt,
		MinTime:         stats.MinRtt,
		MaxTime:         stats.MaxRtt,
		Jitter:          stats.MaxRtt.Milliseconds() - stats.MinRtt.Milliseconds(),
		PacketsSent:     stats.PacketsSent,
		PacketsReceived: stats.PacketsRecv,
		ResponseTimes:   stats.Rtts,
		Host:            host,
	}
}

func Monitor(host string) {
	fmt.Printf("Checking %s\n", host)
	for i := 0; i < numRetries; i++ {
		_, err := ping(host)
		if err != nil {
			logger.Printf("Ping failed for %s: %s\n", host, err)
			fmt.Printf("Ping failed for %s: %s\n", host, err)
			time.Sleep(timeout)
			continue
		}
		fmt.Printf("Ping successful for %s\n", host)
		logger.Printf("Ping successful for %s\n", host)
		return
	}
	logger.Printf("Unable to reach %s after %d retries\n", host, numRetries)
	fmt.Printf("Unable to reach %s after %d retries\n", host, numRetries)
}

func ping(host string) (PingResponse, error) {
	ipAddr := net.ParseIP(host)
	if ipAddr == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Println("failed to resolve domain name")
			return PingResponse{}, err
		}
		ipAddr = ips[0]
	}

	pinger, err := probing.NewPinger(ipAddr.String())
	if err != nil {
		fmt.Println("failed to initialize ping")
		return PingResponse{}, err
	}

	pinger.Count = 3
	pinger.Timeout = 3

	pinger.Run()

	stats := parseToPingResponse(pinger.Statistics(), host)

	return stats, nil
}
