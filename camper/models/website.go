package models

type Website struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Hash        string `json:"hash"`
	Time        int    `json:"time"`
	Cron_Job_Id string `json:"cron_job_id"`
}
