package jobs

import (
	"sync"
	"time"
)

type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type JobInfo struct {
	ID        string
	Type      string
	Status    JobStatus
	Progress  int
	StartedAt time.Time
	EndedAt   *time.Time
	Error     *string
}

type JobManager struct {
	mu   sync.Mutex
	jobs map[string]*JobInfo
}

func NewJobManager() *JobManager {
	return &JobManager{jobs: make(map[string]*JobInfo)}
}

func (m *JobManager) CreateJob(id, jt string) *JobInfo {
	m.mu.Lock()
	defer m.mu.Unlock()
	ji := &JobInfo{ID: id, Type: jt, Status: JobPending, StartedAt: time.Now().UTC(), Progress: 0}
	m.jobs[id] = ji
	return ji
}

func (m *JobManager) Update(id string, status JobStatus, progress int, errStr *string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if j, ok := m.jobs[id]; ok {
		j.Status = status
		j.Progress = progress
		if status == JobCompleted || status == JobFailed {
			now := time.Now().UTC()
			j.EndedAt = &now
			j.Error = errStr
		}
	}
}

func (m *JobManager) Get(id string) (*JobInfo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[id]
	return j, ok
}
