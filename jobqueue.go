package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type TrackInfo struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Link   string `json:"link"`
}

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID         string      `json:"job_id"`
	Query      string      `json:"query"`
	Status     JobStatus   `json:"status"`
	Result     []TrackInfo `json:"result,omitempty"`
	Error      string      `json:"error,omitempty"`
	SubmittedAt time.Time   `json:"-"` // For potential future use (e.g. TTL)
	OutputFilePath string   `json:"-"` // Internal: path to the CLI's output file
}

var (
	jobQueue      chan *Job
	jobsMap       map[string]*Job
	jobsMapMutex  = &sync.Mutex{}
	appConfig     *AppConfig
)

func InitializeJobQueue(bufferSize int, cfg *AppConfig) {
	jobQueue = make(chan *Job, bufferSize)
	jobsMap = make(map[string]*Job)
	appConfig = cfg // Store config for worker
	go worker()
}

func SubmitJob(query string) *Job {
	jobsMapMutex.Lock()
	defer jobsMapMutex.Unlock()

	jobID := uuid.New().String()
	outputFilePath := appConfig.TempOutputDir + "/soulvibe_out_" + jobID + ".json"

	job := &Job{
		ID:         jobID,
		Query:      query,
		Status:     StatusQueued,
		SubmittedAt: time.Now(),
		OutputFilePath: outputFilePath,
	}
	jobsMap[jobID] = job
	jobQueue <- job // Send to worker queue
	return job
}

func GetJobStatus(jobID string) (*Job, bool) {
	jobsMapMutex.Lock()
	defer jobsMapMutex.Unlock()

	job, exists := jobsMap[jobID]
	return job, exists
}

func worker() {
	for job := range jobQueue {
		jobsMapMutex.Lock()
		job.Status = StatusProcessing
		jobsMap[job.ID] = job // Update status in map
		jobsMapMutex.Unlock()

		// Execute the CLI command
		result, err := RunCliCommand(appConfig, job.Query, job.OutputFilePath, job.ID)

		jobsMapMutex.Lock()
		if err != nil {
			job.Status = StatusFailed
			job.Error = err.Error()
		} else {
			job.Status = StatusCompleted
			job.Result = result
		}
		jobsMap[job.ID] = job // Update with final status and result/error
		jobsMapMutex.Unlock()
	}
}