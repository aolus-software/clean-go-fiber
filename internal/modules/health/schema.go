package health

type HomeResponse struct {
	AppName string `json:"appName"`
	AppEnv  string `json:"appEnv"`
	Date    string `json:"date"`
}

type HealthServices struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

type HealthResponse struct {
	Status    string         `json:"status"`
	Timestamp string         `json:"timestamp"`
	Services  HealthServices `json:"services"`
}

type LiveResponse struct {
	Alive bool `json:"alive"`
}
