package db

import (
	"sync"
	"time"
)

type Agent struct {
	ID        string    `json:"id"`
	Hostname  string    `json:"hostname"`
	OS        string    `json:"os"`
	IP        string    `json:"ip"`
	LastSeen  time.Time `json:"last_seen"`
	Status    string    `json:"status"`
}

type Command struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"` // pending, completed
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

type Store struct {
	Agents   map[string]*Agent
	Commands map[string][]*Command // Map AgentID to list of commands
	mu       sync.RWMutex
}

var GlobalStore *Store

func Init() {
	GlobalStore = &Store{
		Agents:   make(map[string]*Agent),
		Commands: make(map[string][]*Command),
	}
}

func (s *Store) AddAgent(a *Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Agents[a.ID] = a
}

func (s *Store) GetAgent(id string) (*Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.Agents[id]
	return a, ok
}

func (s *Store) UpdateLastSeen(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if agent, ok := s.Agents[id]; ok {
		agent.LastSeen = time.Now()
		agent.Status = "active"
	}
}

func (s *Store) AddCommand(cmd *Command) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Commands[cmd.AgentID] = append(s.Commands[cmd.AgentID], cmd)
}

func (s *Store) GetPendingCommands(agentID string) []*Command {
	s.mu.Lock()
	defer s.mu.Unlock() // Lock full time to modify status
	
	var pending []*Command
	if cmds, ok := s.Commands[agentID]; ok {
		for _, cmd := range cmds {
			if cmd.Status == "pending" {
				pending = append(pending, cmd)
				cmd.Status = "sent"
			}
		}
	}
	return pending
}

func (s *Store) UpdateCommandResult(cmdID string, result string, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Inefficient search for now, better structure needed for prod
	for _, cmds := range s.Commands {
		for _, cmd := range cmds {
			if cmd.ID == cmdID {
				cmd.Result = result
				if errorMsg != "" {
					cmd.Result += "\nError: " + errorMsg
				}
				cmd.Status = "completed"
				return
			}
		}
	}
}

func (s *Store) GetAllAgents() []*Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	agents := make([]*Agent, 0, len(s.Agents))
	for _, a := range s.Agents {
		agents = append(agents, a)
	}
	return agents
}
