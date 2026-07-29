package dep

import (
	"fmt"
	"slices"
)

type DependencySorter interface {
	Name() string
	Dependencies() []string
}

// Sort returns the list of sorted items based on
// depdendencies. If there is an error (such as a cycle)
// then the original unsorted list is returned.
func Sort[T DependencySorter](all []T) ([]T, error) {
	s := newDepSorter(all)
	// TODO: This whole thng could be more efficient
	// but not my focus right now.
	inDegree := s.calculateInDegrees()
	queue := s.initializeQueue(inDegree)
	sorted, err := s.process(inDegree, queue)
	if err != nil {
		return all, err
	}
	all = all[:0]
	for _, name := range sorted {
		all = append(all, s.entryValue(name))
	}
	return all, nil
}

type depEntry[T DependencySorter] struct {
	name     string
	deps     []string
	inDegree int
	t        T
}

func newDepSorter[T DependencySorter](all []T) *depSorter[T] {
	s := &depSorter[T]{}
	s.entries = make([]depEntry[T], 0, len(all))
	for _, a := range all {
		e := depEntry[T]{name: a.Name(), deps: a.Dependencies(), t: a}
		s.entries = append(s.entries, e)
	}
	return s
}

type depSorter[T DependencySorter] struct {
	entries []depEntry[T]
}

func (s *depSorter[T]) entryValue(name string) T {
	for _, e := range s.entries {
		if e.name == name {
			return e.t
		}
	}
	// Total error. Should never be possible.
	var t T
	return t
}

func (s *depSorter[T]) calculateInDegrees() map[string]int {
	inDegree := make(map[string]int)
	for _, e := range s.entries {
		inDegree[e.name] = 0
	}
	for _, e := range s.entries {
		for range e.deps {
			inDegree[e.name]++
		}
	}
	return inDegree
}

func (s *depSorter[T]) initializeQueue(inDegree map[string]int) []string {
	queue := []string{}
	for itemName, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, itemName)
		}
	}
	return queue
}

func (s *depSorter[T]) process(inDegree map[string]int, queue []string) ([]string, error) {
	sortedList := []string{}
	itemsProcessed := 0
	for len(queue) > 0 {
		// Dequeue an item
		itemName := queue[0]
		queue = queue[1:]

		// Add to sorted list
		sortedList = append(sortedList, itemName)
		itemsProcessed++

		// Update dependencies
		for _, e := range s.entries {
			dependentName := e.name
			if slices.Contains(e.deps, itemName) {
				inDegree[dependentName]--
				if inDegree[dependentName] == 0 {
					queue = append(queue, dependentName)
				}
			}
		}
	}

	// check for cycles
	if itemsProcessed != len(s.entries) {
		return nil, fmt.Errorf("cycle detected; cannot sort items")
	}

	return sortedList, nil
}
