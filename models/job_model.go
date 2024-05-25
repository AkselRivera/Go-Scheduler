package models

type Job struct {
	ID            string      `json:"_id"`
	Command       string      `json:"command"`
	CreationDate  string      `json:"creation_date"`
	CustomerID    string      `json:"customer_id"`
	ExecutionDate string      `json:"execution_date"`
	HostId        string      `json:"hostId"`
	Interpreter   string      `json:"interpreter"`
	LastAttemptAt interface{} `json:"last_attempt_at"`
	Platform      string      `json:"platform"`
	RelatedJobId  string      `json:"related_job_id"`
	ScenarioId    string      `json:"scenario_id"`
	Status        string      `json:"status"`
	TaskType      string      `json:"task_type"`
	Timeout       float64     `json:"timeout"`
}

type JobDetails struct {
	Job           Job    `json:"job"`
	CustomerToken string `json:"customer_token"`
	SoardId       string `json:"soard_id"`
	Freeze_time   int    `json:"freeze_time,omitempty"`
}

type JobResults struct {
	ID            string                 `json:"_id"`
	Batuta        map[string]interface{} `json:"batuta"`
	Report        string                 `json:"report_job_id,omitempty"`
	LastAttemptAt string                 `json:"last_attempted_at"`
}
