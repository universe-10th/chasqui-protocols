package protocols

import "errors"

var ErrCircularDependencies = errors.New("a circular dependency was detected among protocols")

func traverseDependency(dependency Protocol, traversed Protocols, collected map[Protocol]int) error {
	if _, ok := traversed[dependency]; ok {
		return ErrCircularDependencies
	} else if _, ok := collected[dependency]; ok {
		return nil
	}

	traversed[dependency] = true
	defer delete(traversed, dependency)

	for dependency := range dependency.Dependencies() {
		if err := traverseDependency(dependency, traversed, collected); err != nil {
			return err
		}
	}
	collected[dependency] = len(collected)
	return nil
}

func flatten(dependencies []Protocol) ([]Protocol, error) {
	traversed := make(Protocols)
	collected := make(map[Protocol]int)

	for _, dependency := range dependencies {
		if err := traverseDependency(dependency, traversed, collected); err != nil {
			return nil, err
		}
	}

	count := len(collected)
	flat := make([]Protocol, count)
	for protocol, index := range collected {
		flat[index] = protocol
	}
	return flat, nil
}
