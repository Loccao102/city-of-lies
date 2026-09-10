import { create } from "zustand";
import {
  SessionResponse,
  AgentSummary,
  NotebookData,
  WorldEvent,
} from "@/types/game";
import { GameWebSocket } from "@/services/websocket";
import { getNotebook, getAgents } from "@/services/api";
import { sound } from "@/services/sound";

export interface RumorPulse {
  id: string;
  speakerAgentId: string;
  listenerAgentId: string;
  speakerName: string;
  listenerName: string;
  locationId: string;
  category: string;
  timestamp: number;
}

interface GameStore {
  session: SessionResponse | null;
  agents: AgentSummary[];
  notebook: NotebookData | null;
  events: WorldEvent[];
  recentPulses: RumorPulse[];
  ws: GameWebSocket | null;

  setSession: (session: SessionResponse) => void;
  setAgents: (agents: AgentSummary[]) => void;
  setNotebook: (notebook: NotebookData) => void;
  addEvent: (evt: WorldEvent) => void;
  addPulse: (pulse: Omit<RumorPulse, "id" | "timestamp">) => void;
  refreshState: () => Promise<void>;
  connectWS: (sessionId: string) => void;
  disconnectWS: () => void;
}

export const useGameStore = create<GameStore>((set, get) => ({
  session: null,
  agents: [],
  notebook: null,
  events: [],
  recentPulses: [],
  ws: null,

  setSession: (session) => set({ session }),
  setAgents: (agents) => set({ agents }),
  setNotebook: (notebook) => set({ notebook }),

  addEvent: (evt) => {
    set((state) => ({
      events: [evt, ...state.events.slice(0, 49)], // Keep last 50 events
    }));
  },

  addPulse: (pulseData) => {
    const pulse: RumorPulse = {
      ...pulseData,
      id: Math.random().toString(36).substring(7),
      timestamp: Date.now(),
    };
    set((state) => ({
      recentPulses: [...state.recentPulses.slice(-15), pulse],
    }));
  },

  refreshState: async () => {
    const session = get().session;
    if (!session) return;
    try {
      const [nb, ags] = await Promise.all([
        getNotebook(session.id),
        getAgents(session.id),
      ]);
      set({ notebook: nb, agents: ags });
    } catch (e) {
      console.warn("Error refreshing state:", e);
    }
  },

  connectWS: (sessionId: string) => {
    get().disconnectWS();
    const ws = new GameWebSocket(sessionId);

    ws.subscribe((msg) => {
      const { type, payload, game_second } = msg;

      if (type === "belief_stats.updated") {
        sound.updateTension(payload.primary_false_narrative_ratio);
        set((state) => {
          if (!state.session) return state;
          return {
            session: {
              ...state.session,
              game_second: game_second ?? state.session.game_second,
              primary_false_narrative_ratio: payload.primary_false_narrative_ratio,
              evidence_strength: payload.evidence_strength,
              player_credibility: payload.player_credibility,
            },
          };
        });
      } else if (type === "rumor.visualized") {
        sound.playRumorAlert();
        get().addPulse({
          speakerAgentId: payload.speaker_agent_id,
          listenerAgentId: payload.listener_agent_id,
          speakerName: payload.speaker_name,
          listenerName: payload.listener_name,
          locationId: payload.location_id,
          category: payload.category,
        });
      } else if (type === "public_event.occurred") {
        const title = (payload.title ?? "").toLowerCase();
        if (title.includes("cứu hỏa") || title.includes("cứu thương") || title.includes("siren") || title.includes("fire")) {
          sound.playSiren();
        }
      } else if (type === "session.status_changed") {
        set((state) => {
          if (!state.session) return state;
          return {
            session: {
              ...state.session,
              status: payload.status,
            },
          };
        });
      } else if (type === "game.won" || type === "game.lost") {
        if (type === "game.won") {
          sound.playVictory();
        } else {
          sound.playDefeat();
        }
        set((state) => {
          if (!state.session) return state;
          return {
            session: {
              ...state.session,
              status: type === "game.won" ? "won" : "lost_false_belief",
            },
          };
        });
      }

      // Record in event log
      get().addEvent({
        id: Math.random().toString(),
        session_id: sessionId,
        sequence: msg.sequence,
        event_type: type,
        game_second,
        payload,
        occurred_at: new Date().toISOString(),
      });
    });

    ws.connect();
    set({ ws });
  },

  disconnectWS: () => {
    const ws = get().ws;
    if (ws) {
      ws.disconnect();
      set({ ws: null });
    }
  },
}));
