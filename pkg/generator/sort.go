package generator

// NamedFields implements the sort.Sort methodset, allowing you to call
// sort.Sort(NamedFields) on a slice of []Field
type NamedFields []Field

func (x NamedFields) Len() int           { return len(x) }
func (x NamedFields) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x NamedFields) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
