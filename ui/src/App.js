import React, { useState, useEffect } from 'react';
import { Container, Typography, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip, Box } from '@mui/material';

// Mock data generator since we don't have the full API connected in this snippet
const mockAgents = [
  { id: 'agent-a1b2', hostname: 'DESKTOP-FINANCE', platform: 'windows', last_seen: '2s ago', status: 'active' },
  { id: 'agent-c3d4', hostname: 'ubuntu-prod-01', platform: 'linux', last_seen: '45s ago', status: 'active' },
  { id: 'agent-e5f6', hostname: 'macbook-dev', platform: 'darwin', last_seen: '5m ago', status: 'dormant' },
];

function App() {
  const [agents, setAgents] = useState(mockAgents);

  useEffect(() => {
    // In a real app, this would poll /api/agents
    const interval = setInterval(() => {
      // Refresh logic
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 'bold', color: '#1976d2' }}>
          TITAN <span style={{ color: '#555' }}>FRAMEWORK</span>
        </Typography>
        <Typography variant="subtitle1" color="textSecondary">
          Advanced Command & Control Dashboard
        </Typography>
      </Box>

      <Paper elevation={3}>
        <Box sx={{ p: 2, borderBottom: '1px solid #eee' }}>
          <Typography variant="h6">Active Agents</Typography>
        </Box>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Agent ID</TableCell>
                <TableCell>Hostname</TableCell>
                <TableCell>Platform</TableCell>
                <TableCell>Last Seen</TableCell>
                <TableCell>Status</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {agents.map((agent) => (
                <TableRow key={agent.id} hover>
                  <TableCell sx={{ fontFamily: 'monospace' }}>{agent.id}</TableCell>
                  <TableCell>{agent.hostname}</TableCell>
                  <TableCell>
                    <Chip label={agent.platform} size="small" color={agent.platform === 'windows' ? 'primary' : 'default'} />
                  </TableCell>
                  <TableCell>{agent.last_seen}</TableCell>
                  <TableCell>
                    <Chip 
                      label={agent.status} 
                      color={agent.status === 'active' ? 'success' : 'warning'} 
                      size="small" 
                      variant="outlined"
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Container>
  );
}

export default App;
