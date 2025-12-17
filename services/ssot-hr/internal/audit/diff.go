package audit

import "reflect"

// Diff returns a shallow diff between before/after for JSON-friendly maps.
// For nested maps, it marks the key as changed (doesn't compute deep diffs).
func Diff(before, after map[string]any) map[string]any {
	out := map[string]any{}
	keys := map[string]struct{}{}
	for k := range before { keys[k] = struct{}{} }
	for k := range after { keys[k] = struct{}{} }

	for k := range keys {
		bv, bok := before[k]
		av, aok := after[k]
		if !bok && aok {
			out[k] = map[string]any{"from": nil, "to": av}
			continue
		}
		if bok && !aok {
			out[k] = map[string]any{"from": bv, "to": nil}
			continue
		}
		if !reflect.DeepEqual(bv, av) {
			out[k] = map[string]any{"from": bv, "to": av}
		}
	}
	return out
}
