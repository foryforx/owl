package api

// Health describes the health of a service
type HealthResponse struct {
	Healths []Health `json:"healths"`
}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
	Details any    `json:"details,omitempty"`
}
