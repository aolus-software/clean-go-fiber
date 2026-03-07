package health

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type HomeResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
