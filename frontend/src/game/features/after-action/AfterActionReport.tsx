"use client";

import React from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { createSession } from "@/services/api";
import {
  Trophy,
  Skull,
  RotateCcw,
  Sparkles,
  Award,
  AlertTriangle,
  FileCheck2,
  Clock,
} from "lucide-react";

export function AfterActionReport() {
  const session = useGameStore((s) => s.session);
  const setSession = useGameStore((s) => s.setSession);
  const connectWS = useGameStore((s) => s.connectWS);
  const { isAARModalOpen, aarResult, closeAARModal } = useUIStore();

  const isWon = session?.status === "won" || aarResult?.won;
  const isLost = session?.status === "lost_false_belief" || session?.status === "lost_timeout";

  if (!isAARModalOpen && !isWon && !isLost) return null;
  if (!session) return null;

  const handleRestart = async (sameSeed: boolean) => {
    try {
      const seed = sameSeed ? session.simulation_seed : Math.floor(Math.random() * 1000000);
      const newSession = await createSession(session.scenario_id, seed);
      setSession(newSession);
      connectWS(newSession.id);
      closeAARModal();
    } catch (e: any) {
      alert(`Restart failed: ${e.message}`);
    }
  };

  const gameMins = Math.floor(session.game_second / 60);
  const gameSecs = session.game_second % 60;

  return (
    <div className="fixed inset-0 bg-black/90 backdrop-blur-lg z-50 flex items-center justify-center p-4 animate-fade-in">
      <div className="glass-panel-elevated w-full max-w-xl rounded-3xl border border-white/10 shadow-2xl p-8 space-y-6 text-center">
        {/* Outcome Badge */}
        <div className="flex flex-col items-center gap-3">
          {isWon ? (
            <>
              <div className="w-16 h-16 rounded-full bg-emerald-500/20 border-2 border-emerald-400 flex items-center justify-center shadow-lg glow-truth">
                <Trophy className="w-8 h-8 text-emerald-400" />
              </div>
              <h2 className="text-2xl font-black tracking-wide text-emerald-400">
                SỰ THẬT ĐÃ ĐƯỢC BẢO VỆ!
              </h2>
              <p className="text-xs text-gray-300 max-w-md leading-relaxed">
                Bạn đã giải mã chính xác bản chất sự cố và thu thập đủ chứng cứ pháp lý vững chắc trước khi tin đồn độc hại thống trị dư luận.
              </p>
            </>
          ) : (
            <>
              <div className="w-16 h-16 rounded-full bg-red-500/20 border-2 border-red-500 flex items-center justify-center shadow-lg glow-danger">
                <Skull className="w-8 h-8 text-red-500" />
              </div>
              <h2 className="text-2xl font-black tracking-wide text-red-500">
                LỜI NÓI DỐI ĐÃ TRỞ THÀNH SỰ THẬT!
              </h2>
              <p className="text-xs text-gray-300 max-w-md leading-relaxed">
                Tin đồn thất thiệt đã vượt mốc 75% xã hội chấp nhận. Hoang mang bao trùm thành phố và sự thật đã bị vùi lấp.
              </p>
            </>
          )}
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-3 p-4 bg-black/50 border border-white/5 rounded-2xl text-left font-mono text-xs">
          <div className="p-3 bg-surface rounded-xl space-y-1">
            <div className="flex items-center gap-1.5 text-gray-400 text-[11px]">
              <Clock className="w-3.5 h-3.5 text-truth-glow" /> Thời gian phá án:
            </div>
            <div className="text-sm font-bold text-gray-100">
              {gameMins} phút {gameSecs} giây
            </div>
          </div>

          <div className="p-3 bg-surface rounded-xl space-y-1">
            <div className="flex items-center gap-1.5 text-gray-400 text-[11px]">
              <AlertTriangle className="w-3.5 h-3.5 text-danger" /> Tỉ lệ tin giả cuối cùng:
            </div>
            <div className="text-sm font-bold text-gray-100">
              {Math.round(session.primary_false_narrative_ratio * 100)}%
            </div>
          </div>

          <div className="p-3 bg-surface rounded-xl space-y-1">
            <div className="flex items-center gap-1.5 text-gray-400 text-[11px]">
              <FileCheck2 className="w-3.5 h-3.5 text-truth" /> Độ mạnh bằng chứng:
            </div>
            <div className="text-sm font-bold text-gray-100">
              {Math.round(session.evidence_strength * 100)}%
            </div>
          </div>

          <div className="p-3 bg-surface rounded-xl space-y-1">
            <div className="flex items-center gap-1.5 text-gray-400 text-[11px]">
              <Award className="w-3.5 h-3.5 text-credibility-glow" /> Điểm uy tín điều tra:
            </div>
            <div className="text-sm font-bold text-gray-100">
              {Math.round(session.player_credibility * 100)}%
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex items-center justify-center gap-4 pt-2">
          <button
            onClick={() => handleRestart(true)}
            className="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-surface-elevated hover:bg-surface-border text-gray-200 text-xs font-semibold border border-white/10 transition-all shadow-md"
          >
            <RotateCcw className="w-4 h-4" />
            Chơi lại Seed này
          </button>

          <button
            onClick={() => handleRestart(false)}
            className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-truth hover:bg-truth-glow text-white text-xs font-bold transition-all shadow-lg glow-truth"
          >
            <Sparkles className="w-4 h-4" />
            Bắt đầu màn chơi mới
          </button>
        </div>
      </div>
    </div>
  );
}
