package model

// Various AWS services return different representations for tags, so this struct standardizes them from this point downstream
type Tag struct {
	Key   string
	Value string
}

func (t Tag) Render() []string {
	return []string{
		t.Key,
		t.Value,
	}
}

type Tags []Tag

func (t Tags) Headers() []string {
	return []string{
		"KEY",
		"VALUE",
	}
}

func (t Tags) Get(key string) (string, bool) {
	for _, v := range t {
		if v.Key == key {
			return v.Value, true
		}
	}
	return "", false
}
