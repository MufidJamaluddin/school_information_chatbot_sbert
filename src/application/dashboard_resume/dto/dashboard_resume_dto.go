package dto

type DashboardResumeDTO struct {
	TotalUser     uint64 `json:"totalUser"`
	TotalQuestion uint64 `json:"totalQuestion"`
	TotalGreeting uint64 `json:"totalGreeting"`
	TotalAdmin    uint64 `json:"totalAdmin"`
}
