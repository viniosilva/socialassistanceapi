package service

type HealthStatus string

const (
	HealthStatusUp   HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

type HealthResponse struct {
	Status HealthStatus `json:"status"`
}
