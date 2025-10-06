package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/heimdalr/dag"
)

// Directed acyclic graph scheduler

type DAGS struct {
	ctx context.Context
	dag *dag.DAG
}

type Visitor struct {
	ctx context.Context

	revertMu sync.RWMutex
	revert   string
}

func (v *Visitor) Revert() string {
	v.revertMu.RLock()
	defer v.revertMu.RUnlock()

	return v.revert
}

func (v *Visitor) Visit(vertex dag.Vertexer) {
	v.revertMu.RLock()
	if v.revert != "" {
		v.revertMu.RUnlock()
		return
	}
	v.revertMu.RUnlock()
	id, value := vertex.Vertex()
	if vertexValue, ok := value.(Job); ok && vertexValue != nil {
		log.Default().Println("Start: ", id)
		if err := vertexValue.Run(v.ctx); err != nil {
			log.Default().Println("FAIL: ", id, " Err: ", err.Error())
			v.revertMu.RLock()
			if v.revert != "" {
				v.revertMu.RUnlock()
				return
			}
			v.revertMu.RUnlock()
			v.revertMu.Lock()
			defer v.revertMu.Unlock()
			v.revert = id
			return
		}
		log.Default().Println("End: ", id)
	} else {
		panic("invalid type in DAG!!!")
	}
}

func New(ctx context.Context) *DAGS {
	dags := DAGS{
		dag: dag.NewDAG(),
		ctx: ctx,
	}

	return &dags
}

func (d DAGS) Add(j Job, deps ...string) (id string, err error) {
	id, err = d.dag.AddVertex(j)
	if err != nil {
		return id, err
	}

	for _, dep := range deps {
		if erro := d.dag.AddEdge(dep, id); err != nil {
			err = errors.Join(err, erro)
		}
	}

	return
}

func (d DAGS) Start() error {
	ctx, cancel := context.WithCancel(d.ctx)
	defer cancel()
	var visitor dag.Visitor = &Visitor{ctx: ctx}

	d.dag.ReduceTransitively()

	d.dag.OrderedWalk(visitor)

	if vi, ok := visitor.(*Visitor); ok {
		id := vi.Revert()
		if id != "" {
			fmt.Println("FAILURE HAS OCC, WILL REVERT, THE FIRST FAILURE WAS IN", id)
			// revert from the id
		}
	} else {
		fmt.Println("Visitor doesnt impl revert")
	}

	return ctx.Err()
}
