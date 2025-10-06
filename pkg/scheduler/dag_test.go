package scheduler_test

import (
	"context"
	"log"
	"os"
	"os/signal"
	"testing"

	"github.com/kthcloud/cli/pkg/scheduler"
)

func TestDag(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	dag := scheduler.New(ctx)

	helloID, err := dag.Add(&scheduler.JobImpl{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}

	worldID, err := dag.Add(&scheduler.JobImpl{Value: "world"}, helloID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = dag.Add(&scheduler.JobImpl{Value: "Im independant"})
	if err != nil {
		log.Fatal(err)
	}

	independant2, err := dag.Add(&scheduler.JobImpl{Value: "Im also independant"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = dag.Add(&scheduler.JobImpl{Value: "Im dependant on \"Im also independant\""}, independant2)
	if err != nil {
		log.Fatal(err)
	}

	fooID, err := dag.Add(&scheduler.JobImpl{Value: "foo"}, worldID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = dag.Add(&scheduler.JobImpl{Value: "Bar", Fail: true}, fooID)
	if err != nil {
		log.Fatal(err)
	}

	dag.Start()
}
