"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { submitTruth } from "@/services/api";
import { SubmitTruthRequest } from "@/types/game";
import {
  X,
  Flag,
  CheckCircle2,
  AlertCircle,
  FileCheck2,
  Loader2,
} from "lucide-react";

export function SubmitTruthModal() {
  const session = useGameStore((s) => s.session);
  const refreshState = useGameStore((s) => s.refreshState);
  const { isSubmitTruthModalOpen, closeSubmitTruthModal, openAARModal } = useUIStore();

  const [eventType, setEventType] = useState("fire");
  const [causeCategory, setCauseCategory] = useState("electrical_fault");
  const [majorExplosion, setMajorExplosion] = useState(false);
  const [fatalities, setFatalities] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<string | null>(null);

  if (!isSubmitTruthModalOpen || !session) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setFeedback(null);

    const payload: SubmitTruthRequest = {
      event_type: eventType,
      cause_category: causeCategory,
      major_explosion: majorExplosion,
      fatalities: Number(fatalities),
    };

    try {
      const result = await submitTruth(session.id, payload);
      await refreshState();

      if (result.won || result.terminal) {
        closeSubmitTruthModal();
        openAARModal(result);
      } else {
        setFeedback(result.feedback_message || "Báo cáo chưa đủ căn cứ xác thực.");
      }
    } catch (err: any) {
      alert(`Nộp kết luận thất bại: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-md z-50 flex items-center justify-center p-4">
      <div className="glass-panel-elevated w-full max-w-lg rounded-2xl border border-white/10 shadow-2xl overflow-hidden animate-fade-in">
        {/* Header */}
        <div className="px-6 py-4 bg-surface/90 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-emerald-500/20 border border-emerald-500/40 flex items-center justify-center">
              <Flag className="w-5 h-5 text-emerald-400" />
            </div>
            <div>
              <h3 className="font-bold text-white text-base">
                NỘP KẾT LUẬN SỰ THẬT CUỐI CÙNG
              </h3>
              <p className="text-xs text-gray-400">
                Xác lập Ground Truth để kết thúc cuộc điều tra
              </p>
            </div>
          </div>
          <button
            onClick={closeSubmitTruthModal}
            className="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {/* Rules info */}
          <div className="p-3 bg-black/40 border border-white/5 rounded-xl text-xs text-gray-400 space-y-1">
            <div className="flex items-center justify-between font-mono">
              <span>Độ mạnh bằng chứng hiện tại:</span>
              <span className="text-truth-glow font-bold">
                {Math.round(session.evidence_strength * 100)}% (Yêu cầu &ge; 70%)
              </span>
            </div>
            <p className="text-[11px] text-gray-500">
              *Bạn chỉ chiến thắng khi điền đúng sự thật khách quan VÀ đã thu thập đủ &ge;70% độ tin cậy bằng chứng.
            </p>
          </div>

          {feedback && (
            <div className="p-3 bg-amber-950/40 border border-amber-500/40 rounded-xl text-xs text-amber-300 flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-amber-400 flex-shrink-0 mt-0.5" />
              <span>{feedback}</span>
            </div>
          )}

          {/* 1. Event Type */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200">
              Bản chất sự cố (Event Type):
            </label>
            <select
              value={eventType}
              onChange={(e) => setEventType(e.target.value)}
              className="w-full bg-black/50 border border-white/10 rounded-xl p-2.5 text-xs text-white focus:outline-none focus:border-truth"
            >
              <option value="fire">Cháy thông thường (Fire)</option>
              <option value="chemical_explosion">Nổ hóa chất (Chemical Explosion)</option>
              <option value="gas_explosion">Nổ khí gas (Gas Explosion)</option>
              <option value="unknown">Không rõ nguyên nhân</option>
            </select>
          </div>

          {/* 2. Cause Category */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200">
              Nguồn gốc kỹ thuật (Cause Category):
            </label>
            <select
              value={causeCategory}
              onChange={(e) => setCauseCategory(e.target.value)}
              className="w-full bg-black/50 border border-white/10 rounded-xl p-2.5 text-xs text-white focus:outline-none focus:border-truth"
            >
              <option value="electrical_fault">Chập điện / Hỏng hóc kỹ thuật (Electrical Fault)</option>
              <option value="chemical_reaction">Phản ứng hóa học độc hại (Chemical Reaction)</option>
              <option value="arson">Phóng hỏa có chủ đích (Arson)</option>
              <option value="unknown">Chưa xác định</option>
            </select>
          </div>

          {/* 3. Major Explosion */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200 block">
              Có vụ nổ lớn phá hủy cấu trúc không?
            </label>
            <div className="flex items-center gap-6">
              <label className="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
                <input
                  type="radio"
                  name="explosion"
                  checked={majorExplosion === false}
                  onChange={() => setMajorExplosion(false)}
                  className="accent-truth"
                />
                <span>Không có nổ lớn (Chỉ có tiếng nổ chập aptomat)</span>
              </label>
              <label className="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
                <input
                  type="radio"
                  name="explosion"
                  checked={majorExplosion === true}
                  onChange={() => setMajorExplosion(true)}
                  className="accent-truth"
                />
                <span>Có nổ lớn</span>
              </label>
            </div>
          </div>

          {/* 4. Fatalities */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200 block">
              Số người tử vong (Fatalities):
            </label>
            <input
              type="number"
              min={0}
              max={20}
              value={fatalities}
              onChange={(e) => setFatalities(Number(e.target.value))}
              className="w-full bg-black/50 border border-white/10 rounded-xl p-2.5 text-xs text-white focus:outline-none focus:border-truth font-mono"
            />
          </div>

          {/* Footer */}
          <div className="flex items-center justify-end gap-3 pt-3 border-t border-white/10">
            <button
              type="button"
              onClick={closeSubmitTruthModal}
              className="px-4 py-2 rounded-xl text-xs text-gray-400 hover:text-white"
            >
              Hủy
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold rounded-xl flex items-center gap-2 transition-all glow-truth disabled:opacity-50 shadow-lg"
            >
              {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Flag className="w-4 h-4" />}
              Gửi báo cáo sự thật
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
