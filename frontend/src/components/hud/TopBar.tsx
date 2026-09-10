"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { sound } from "@/services/sound";
import {
  Clock,
  AlertTriangle,
  FileCheck2,
  Award,
  BookOpen,
  Send,
  Flag,
  Volume2,
  VolumeX,
} from "lucide-react";

export function TopBar() {
  const session = useGameStore((s) => s.session);
  const { openNotebookModal, openCorrectionModal, openSubmitTruthModal } =
    useUIStore();
  const [isMuted, setIsMuted] = useState(sound.getIsMuted());

  const handleToggleAudio = () => {
    sound.resume();
    const muted = sound.toggleMute();
    setIsMuted(muted);
    if (!muted) {
      sound.playClick();
    }
  };

  if (!session) return null;

  // Format Game Time starting from 17:58:00
  const totalSeconds = 17 * 3600 + 58 * 60 + session.game_second;
  const hours = Math.floor(totalSeconds / 3600) % 24;
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  const timeStr = `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
    2,
    "0"
  )}:${String(seconds).padStart(2, "0")}`;

  const falsePercent = Math.round(session.primary_false_narrative_ratio * 100);
  const evidencePercent = Math.round(session.evidence_strength * 100);
  const credibilityPercent = Math.round(session.player_credibility * 100);

  // Meter warning state
  let meterColor = "bg-blue-500";
  let meterTextColor = "text-blue-400";
  let meterPulse = "";
  if (falsePercent >= 70) {
    meterColor = "bg-red-600";
    meterTextColor = "text-red-400 font-bold";
    meterPulse = "animate-pulse glow-danger";
  } else if (falsePercent >= 60) {
    meterColor = "bg-orange-500";
    meterTextColor = "text-orange-400 font-semibold";
  } else if (falsePercent >= 50) {
    meterColor = "bg-yellow-500";
    meterTextColor = "text-yellow-400";
  }

  return (
    <header className="fixed top-0 left-0 right-0 h-16 glass-panel border-b border-white/10 z-30 px-6 flex items-center justify-between shadow-2xl">
      {/* Title & Time */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-truth animate-ping" />
          <span className="font-extrabold tracking-widest text-sm text-white/90">
            CITY OF LIES
          </span>
        </div>

        <div className="flex items-center gap-2 px-3 py-1 bg-black/40 rounded-md border border-white/10 font-mono text-sm text-gray-200">
          <Clock className="w-4 h-4 text-truth-glow" />
          <span>{timeStr}</span>
        </div>
      </div>

      {/* Center Vitals: False Belief / Evidence / Credibility */}
      <div className="flex items-center gap-8">
        {/* False Narrative Ratio */}
        <div className={`flex flex-col gap-1 w-48 ${meterPulse}`}>
          <div className="flex justify-between text-xs font-mono">
            <span className="flex items-center gap-1 text-gray-400">
              <AlertTriangle className="w-3.5 h-3.5 text-danger" /> Tin đồn sai lệch
            </span>
            <span className={meterTextColor}>{falsePercent}% / 75%</span>
          </div>
          <div className="w-full h-2 bg-black/60 rounded-full overflow-hidden border border-white/10">
            <div
              className={`h-full transition-all duration-500 rounded-full ${meterColor}`}
              style={{ width: `${Math.min(100, (falsePercent / 75) * 100)}%` }}
            />
          </div>
        </div>

        {/* Evidence Strength */}
        <div className="flex flex-col gap-1 w-44">
          <div className="flex justify-between text-xs font-mono">
            <span className="flex items-center gap-1 text-gray-400">
              <FileCheck2 className="w-3.5 h-3.5 text-truth-glow" /> Bằng chứng
            </span>
            <span className="text-truth-glow font-mono font-semibold">
              {evidencePercent}% / 70%
            </span>
          </div>
          <div className="w-full h-2 bg-black/60 rounded-full overflow-hidden border border-white/10">
            <div
              className="h-full bg-truth rounded-full transition-all duration-500"
              style={{ width: `${Math.min(100, evidencePercent)}%` }}
            />
          </div>
        </div>

        {/* Player Credibility */}
        <div className="flex flex-col gap-1 w-36">
          <div className="flex justify-between text-xs font-mono">
            <span className="flex items-center gap-1 text-gray-400">
              <Award className="w-3.5 h-3.5 text-credibility-glow" /> Uy tín
            </span>
            <span className="text-credibility-glow font-mono font-semibold">
              {credibilityPercent}%
            </span>
          </div>
          <div className="w-full h-2 bg-black/60 rounded-full overflow-hidden border border-white/10">
            <div
              className="h-full bg-credibility rounded-full transition-all duration-500"
              style={{ width: `${Math.min(100, credibilityPercent)}%` }}
            />
          </div>
        </div>
      </div>

      {/* Action Buttons & Audio Toggle */}
      <div className="flex items-center gap-3">
        <button
          onClick={() => {
            sound.playPaperRustle();
            openNotebookModal("evidence");
          }}
          className="flex items-center gap-2 px-3.5 py-1.5 bg-surface-elevated hover:bg-surface-border text-gray-200 text-xs font-medium rounded-lg border border-white/10 transition-colors shadow-sm"
        >
          <BookOpen className="w-4 h-4 text-truth-glow" />
          Sổ tay điều tra
        </button>

        <button
          onClick={() => {
            sound.playClick();
            openCorrectionModal();
          }}
          className="flex items-center gap-2 px-3.5 py-1.5 bg-amber-950/40 hover:bg-amber-900/60 text-amber-200 text-xs font-medium rounded-lg border border-amber-600/40 transition-colors shadow-sm"
        >
          <Send className="w-4 h-4 text-amber-400" />
          Đính chính tin đồn
        </button>

        <button
          onClick={() => {
            sound.playClick();
            openSubmitTruthModal();
          }}
          className="flex items-center gap-2 px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-semibold rounded-lg shadow-md transition-all glow-truth"
        >
          <Flag className="w-4 h-4" />
          Nộp sự thật
        </button>

        <div className="w-[1px] h-6 bg-white/10 mx-1" />

        {/* Audio Mute/Unmute */}
        <button
          onClick={handleToggleAudio}
          title={isMuted ? "Bật âm thanh (Sound FX & Ambient)" : "Tắt âm thanh"}
          className={`p-2 rounded-lg border transition-all ${
            isMuted
              ? "bg-red-950/40 border-red-500/30 text-red-400 hover:bg-red-900/50"
              : "bg-surface-elevated border-white/10 text-truth-glow hover:bg-surface-border"
          }`}
        >
          {isMuted ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4 animate-pulse" />}
        </button>
      </div>
    </header>
  );
}
