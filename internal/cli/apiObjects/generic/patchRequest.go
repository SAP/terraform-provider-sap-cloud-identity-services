package generic

type PatchRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

type PatchRequestBody struct {
	Operations []PatchRequest `json:"operations"`
}
