package warnings

import "errors"

type Warning error

func New(text string) Warning {
	return errors.New(text)
}

func Is(warning, target Warning) bool {
	return errors.Is(warning, target)
}

func Join(warnings ...Warning) error {
	errs := make([]error, len(warnings))
	for i, w := range warnings {
		errs[i] = w
	}
	return errors.Join(errs...)
}
