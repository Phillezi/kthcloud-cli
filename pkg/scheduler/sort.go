package scheduler

import (
	"fmt"
)

type Dependent interface {
	GetDependencies() []string
}

func SortByDependencies[T Dependent](nodes map[string]T) ([]string, error) {
	visited := map[string]bool{}
	temp := map[string]bool{}
	var result []string

	var visit func(string) error
	visit = func(n string) error {
		if temp[n] {
			return fmt.Errorf("circular dependency detected at %s", n)
		}

		if !visited[n] {
			temp[n] = true

			for _, dep := range nodes[n].GetDependencies() {
				if _, ok := nodes[dep]; !ok {
					return fmt.Errorf("missing dependency %s for %s", dep, n)
				}
				if err := visit(dep); err != nil {
					return err
				}
			}

			visited[n] = true
			temp[n] = false
			result = append(result, n)
		}

		return nil
	}

	for k := range nodes {
		if !visited[k] {
			if err := visit(k); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func Batch[T Dependent](s Scheduler, dependents map[string]T, mapper func(key string, dependent T) Job) error {
	order, err := SortByDependencies(dependents)
	if err != nil {
		return err
	}

	for _, key := range order {
		if err := s.AddNamed(key, mapper(key, dependents[key])); err != nil {
			return err
		}
	}

	return nil
}
