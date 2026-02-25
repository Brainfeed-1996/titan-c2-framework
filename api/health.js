export default function handler(req, res) { res.status(200).json({ status: 'ok', service: 'titan-c2-framework', timestamp: new Date().toISOString() }); }
