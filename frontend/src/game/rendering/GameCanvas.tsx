"use client";

import React, { useEffect, useRef } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { GameScene } from "../core/GameScene";
import { InvestigationOverlay } from "../features/investigation/InvestigationOverlay";

export function GameCanvas() {
  const mountRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<GameScene | null>(null);
  const agents = useGameStore((state) => state.agents);
  const openInterviewModal = useUIStore((state) => state.openInterviewModal);
  const openInspectModal = useUIStore((state) => state.openInspectModal);

  useEffect(() => {
    if (!mountRef.current) return;

    const scene = new GameScene(mountRef.current, {
      onAgentSelected: openInterviewModal,
      onLocationSelected: openInspectModal,
    });
    sceneRef.current = scene;
    scene.setAgents(useGameStore.getState().agents);

    return () => {
      scene.dispose();
      sceneRef.current = null;
    };
  }, [openInterviewModal, openInspectModal]);

  useEffect(() => {
    sceneRef.current?.setAgents(agents);
  }, [agents]);

  return (
    <div className="relative h-full w-full">
      <div ref={mountRef} className="h-full w-full cursor-grab active:cursor-grabbing" />
      <InvestigationOverlay />
    </div>
  );
}
