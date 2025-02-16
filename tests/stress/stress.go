package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	rpsFlag := flag.Int("rps", 100, "Requests per second")
	durationFlag := flag.Duration("duration", 1*time.Second, "Test duration (e.g. 5s, 1m)")
	flag.Parse()

	if *rpsFlag <= 0 {
		log.Fatal("rps must be a positive integer")
	}
	if *durationFlag <= 0 {
		log.Fatal("duration must be a positive value")
	}

	targetURL := "http://localhost:8080/api/auth"
	authPayload := AuthRequest{
		Username: "test_user",
		Password: "test_password",
	}

	payloadBytes, err := json.Marshal(authPayload)
	if err != nil {
		log.Fatalf("Error creating JSON body: %v", err)
	}

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    targetURL,
		Body:   payloadBytes,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	})

	rate := vegeta.Rate{Freq: *rpsFlag, Per: time.Second}
	duration := *durationFlag

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "Auth Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	extractMetrics(&metrics)
}

func extractMetrics(metrics *vegeta.Metrics) {
	rps := metrics.Rate
	fmt.Printf("RPS (Requests Per Second): %.2f\n", float64(rps))

	p50 := metrics.Latencies.P50
	p95 := metrics.Latencies.P95
	p99 := metrics.Latencies.P99
	fmt.Printf("P50 Latency: %v\n", p50)
	fmt.Printf("P95 Latency: %v\n", p95)
	fmt.Printf("P99 Latency: %v\n", p99)

	successRate := metrics.Success * 100
	fmt.Printf("Success Rate: %.2f%%\n", successRate)
}
