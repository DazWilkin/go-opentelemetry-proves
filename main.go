package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"

	cloudmonitoring "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
)

func main() {
	projectID := os.Getenv("PROJECT")
	if projectID == "" {
		log.Fatal("Environment variable `PROJECT` is required and was unset")
	}

	log.Printf("[main] ProjectID: %s", projectID)
	opts := []cloudmonitoring.Option{
		cloudmonitoring.WithProjectID(projectID),
	}
	popts := []push.Option{
		// push.WithResource(
		// 	resource.New(
		// 		label.String("env", "dev"),
		// 		label.String("author", "google"),
		// 		label.String("application", "example"),
		// 	),
		// ),
	}
	pusher, err := cloudmonitoring.InstallNewPipeline(opts, popts...)

	// opts := []stdout.Option{
	// 	stdout.WithPrettyPrint(),
	// }
	// pusher, err := stdout.InstallNewPipeline(opts, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer pusher.Stop()

	// meter := global.Meter("cloudmonitoring/example")
	meter := pusher.Provider().Meter("cloudmonitoring/example")

	counter := metric.Must(meter).NewInt64Counter("counter.foo")

	lemonsKey := label.Key("lemons")
	labels := []label.KeyValue{
		lemonsKey.String("test"),
		label.Key("key").String("value"),
		label.String("dog", "Freddie"),
	}

	ctx := context.Background()

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	timer := time.NewTicker(10 * time.Second)
	for range timer.C {
		v := r.Int63n(100)
		log.Printf("[main] v=%02d", v)
		counter.Add(ctx, v, labels...)
	}
}
