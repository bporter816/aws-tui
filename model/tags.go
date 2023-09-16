package model

// Various AWS services return different representations for tags, so this struct standardizes them from this point downstream
type Tag struct {
	Key   string
	Value string
}
