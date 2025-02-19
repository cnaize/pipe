package modify

type ModifyFn func(data map[string]any) error

var (
	NopModifier ModifyFn = func(data map[string]any) error { return nil }
)
