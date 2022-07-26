package kiku

// Organization is the struct has ID and name.
type Organization struct {
	ID   int    `json:"organizationId"` // 組織ID
	Name string `json:"name"`           // 組織名
}
