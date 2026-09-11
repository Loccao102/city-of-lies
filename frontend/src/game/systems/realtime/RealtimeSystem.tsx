"use client";

import { useEffect } from "react";
import { sound } from "@/services/sound";
import { useGameStore } from "@/store/gameStore";
import { GameWebSocket, type WSEventMessage } from "./GameWebSocket";

function handleMessage(sessionId: string, message: WSEventMessage) {
  const store = useGameStore.getState();
  const { type, payload, game_second } = message;

  if (type === "belief_stats.updated") {
    sound.updateTension(payload.primary_false_narrative_ratio);
    store.patchSession({
      game_second: game_second ?? store.session?.game_second ?? 0,
      primary_false_narrative_ratio: payload.primary_false_narrative_ratio,
      evidence_strength: payload.evidence_strength,
      player_credibility: payload.player_credibility,
    });
  } else if (type === "rumor.visualized") {
    sound.playRumorAlert();
    store.addPulse({
      speakerAgentId: payload.speaker_agent_id,
      listenerAgentId: payload.listener_agent_id,
      speakerName: payload.speaker_name,
      listenerName: payload.listener_name,
      locationId: payload.location_id,
      category: payload.category,
    });
  } else if (type === "agent.moved") {
    const agentId = payload.agent_id ?? payload.actor_agent_id;
    const locationId = payload.to_location_id ?? payload.location_id;
    if (agentId && locationId) store.moveAgent(agentId, locationId);
  } else if (type === "public_event.occurred") {
    const title = String(payload.title ?? "").toLowerCase();
    if (title.includes("cứu hỏa") || title.includes("cứu thương") || title.includes("siren") || title.includes("fire")) {
      sound.playSiren();
    }
  } else if (type === "session.status_changed") {
    store.patchSession({ status: payload.status });
  } else if (type === "game.won" || type === "game.lost") {
    if (type === "game.won") sound.playVictory();
    else sound.playDefeat();
    store.patchSession({ status: type === "game.won" ? "won" : "lost_false_belief" });
  }

  store.addEvent({
    id: `${message.sequence}-${Date.now()}`,
    session_id: sessionId,
    sequence: message.sequence,
    event_type: type,
    game_second,
    payload,
    occurred_at: new Date().toISOString(),
  });
}

/** Owns WebSocket lifecycle. Zustand stores state only; UI never opens sockets. */
export function RealtimeSystem() {
  const sessionId = useGameStore((state) => state.session?.id ?? null);

  useEffect(() => {
    if (!sessionId) return;

    const socket = new GameWebSocket(sessionId);
    const unsubscribe = socket.subscribe((message) => handleMessage(sessionId, message));
    socket.connect();

    return () => {
      unsubscribe();
      socket.disconnect();
    };
  }, [sessionId]);

  return null;
}
