package cmd

type user struct {
	ID string `json:"id"`
}

type timesheet struct {
	Notes     string `json:"notes"`
	Hours     string `json:"hours"`
	ProjectID string `json:"project_id"`
	TaskID    string `json:"task_id"`
	SpentAt   string `json"spend_at`
}
