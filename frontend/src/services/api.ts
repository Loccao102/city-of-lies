import {
  SessionResponse,
  ScenarioSummary,
  AgentSummary,
  AgentDetail,
  NotebookData,
  DialogueResponse,
  InspectResponse,
  Publication,
  SubmitTruthRequest,
  TruthSubmissionResult,
} from "@/types/game";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

export async function listScenarios(): Promise<ScenarioSummary[]> {
  const res = await fetch(`${API_BASE}/scenarios`);
  if (!res.ok) {
    throw new Error("Failed to fetch scenarios list");
  }
  return res.json();
}

export async function createSession(scenarioId = "random", seed?: number): Promise<SessionResponse> {
  const res = await fetch(`${API_BASE}/sessions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ scenario_id: scenarioId, seed }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Failed to create session: ${text}`);
  }
  return res.json();
}


export async function getSession(sessionId: string): Promise<SessionResponse> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}`);
  if (!res.ok) throw new Error("Failed to fetch session");
  return res.json();
}

export async function getAgents(sessionId: string): Promise<AgentSummary[]> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/agents`);
  if (!res.ok) throw new Error("Failed to fetch agents");
  return res.json();
}

export async function getAgentDetail(sessionId: string, agentId: string): Promise<AgentDetail> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/agents/${agentId}`);
  if (!res.ok) throw new Error("Failed to fetch agent details");
  return res.json();
}

export async function interviewAgent(
  sessionId: string,
  agentId: string,
  message: string
): Promise<DialogueResponse> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/agents/${agentId}/interviews`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Interview failed: ${text}`);
  }
  return res.json();
}

export async function inspectLocation(
  sessionId: string,
  locationId: string
): Promise<InspectResponse> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/locations/${locationId}/inspect`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Inspection failed: ${text}`);
  }
  return res.json();
}

export async function getNotebook(sessionId: string): Promise<NotebookData> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/notebook`);
  if (!res.ok) throw new Error("Failed to fetch notebook");
  return res.json();
}

export async function publishCorrection(
  sessionId: string,
  data: { challenged_claim_id: string; evidence_ids: string[]; message: string }
): Promise<Publication> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/corrections`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Publish correction failed: ${text}`);
  }
  return res.json();
}

export async function submitTruth(
  sessionId: string,
  data: SubmitTruthRequest
): Promise<TruthSubmissionResult> {
  const res = await fetch(`${API_BASE}/sessions/${sessionId}/truth-submissions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Submit truth failed: ${text}`);
  }
  return res.json();
}
