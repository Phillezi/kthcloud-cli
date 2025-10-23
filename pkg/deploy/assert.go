package deploy

import "fmt"

func HandleAndAssert[T any](resp any, stage ...string) (T, error) {
	var zero T

	resourceName := fmt.Sprintf("%T", resp)
	actualStage := "notprovided"
	if len(stage) > 0 && stage[0] != "" {
		actualStage = stage[0]
	}

	obj, err := HandleAPIResponse(resourceName, actualStage, resp)
	if err != nil {
		return zero, err
	}

	typedObj, ok := obj.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected type from %s %s: got %T, expected %T",
			resourceName, actualStage, obj, zero)
	}

	return typedObj, nil
}
