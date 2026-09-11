"use client";

import React from "react";
import { useGameStore } from "@/store/gameStore";

export function InvestigationOverlay() {
  const recentPulses = useGameStore((state) => state.recentPulses);
  const latestPulse = recentPulses.at(-1);

  return (
    <>
      <div className="pointer-events-none absolute left-6 top-4 z-10 flex items-center gap-3 rounded-lg border border-white/10 px-3.5 py-2 text-[11px] text-gray-400 glass-panel">
        <span>Kéo chuột để di chuyển camera</span>
        <span>•</span>
        <span>Lăn chuột để phóng to/thu nhỏ</span>
        <span>•</span>
        <span className="text-truth-glow">Bấm vào người hoặc công trình để tương tác</span>
      </div>

      {latestPulse && (
        <div className="absolute right-6 top-4 z-10 flex items-center gap-3 rounded-xl border border-amber-500/30 px-4 py-2.5 shadow-xl animate-fade-in glass-panel-elevated">
          <div className="h-2.5 w-2.5 animate-ping rounded-full bg-amber-400" />
          <div className="text-xs">
            <span className="font-semibold text-amber-300">{latestPulse.speakerName}</span>
            <span className="text-gray-400"> vừa kể cho </span>
            <span className="font-semibold text-gray-200">{latestPulse.listenerName}</span>
          </div>
        </div>
      )}
    </>
  );
}
