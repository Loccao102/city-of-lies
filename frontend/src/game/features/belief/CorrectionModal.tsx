"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { publishCorrection } from "@/services/api";
import {
  X,
  Send,
  AlertTriangle,
  FileCheck2,
  CheckCircle2,
  Loader2,
  TrendingDown,
} from "lucide-react";

export function CorrectionModal() {
  const session = useGameStore((s) => s.session);
  const notebook = useGameStore((s) => s.notebook);
  const refreshState = useGameStore((s) => s.refreshState);
  const { isCorrectionModalOpen, closeCorrectionModal } = useUIStore();

  const [challengedClaimId, setChallengedClaimId] = useState("claim_multiple_deaths");
  const [selectedEvIds, setSelectedEvIds] = useState<string[]>([]);
  const [message, setMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<any | null>(null);

  if (!isCorrectionModalOpen || !session) return null;

  const targetRumors = [
    { id: "claim_multiple_deaths", text: "Nhiều công nhân đã chết (Tin đồn tử vong)" },
    { id: "claim_chemical_explosion", text: "Nhà máy xảy ra nổ hóa chất (Tin đồn nguồn nổ)" },
    { id: "claim_company_coverup", text: "Công ty đang che giấu số người chết (Nghi vấn bưng bít)" },
  ];

  const handleToggleEvidence = (evId: string) => {
    setSelectedEvIds((prev) =>
      prev.includes(evId) ? prev.filter((id) => id !== evId) : [...prev, evId]
    );
  };

  const handlePublish = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!message.trim() || submitting) return;
    setSubmitting(true);
    setResult(null);
    try {
      const pub = await publishCorrection(session.id, {
        challenged_claim_id: challengedClaimId,
        evidence_ids: selectedEvIds,
        message: message.trim(),
      });
      setResult(pub);
      await refreshState();
    } catch (err: any) {
      alert(`Đính chính thất bại: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/75 backdrop-blur-md z-50 flex items-center justify-center p-4">
      <div className="glass-panel-elevated w-full max-w-xl rounded-2xl border border-white/10 shadow-2xl overflow-hidden animate-fade-in">
        {/* Header */}
        <div className="px-6 py-4 bg-surface/90 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-amber-500/20 border border-amber-500/40 flex items-center justify-center">
              <Send className="w-5 h-5 text-amber-400" />
            </div>
            <div>
              <h3 className="font-bold text-white text-base">
                PHÁT THÔNG BÁO ĐÍNH CHÍNH
              </h3>
              <p className="text-xs text-gray-400">
                Xuất bản báo cáo kiểm chứng để chặn đà lan truyền của tin đồn
              </p>
            </div>
          </div>
          <button
            onClick={closeCorrectionModal}
            className="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {result ? (
          <div className="p-6 space-y-5">
            <div className="p-4 rounded-xl bg-emerald-950/40 border border-emerald-500/40 space-y-3 text-xs">
              <div className="flex items-center gap-2 font-bold text-emerald-300 text-sm">
                <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                Đính chính đã được xuất bản tới toàn quận!
              </div>
              <p className="text-gray-200 leading-relaxed">&ldquo;{result.message}&rdquo;</p>
              <div className="grid grid-cols-2 gap-3 pt-2 font-mono text-[11px] border-t border-emerald-500/20">
                <div>
                  <span className="text-gray-400">Sức mạnh phản biện: </span>
                  <span className="text-emerald-400 font-bold">
                    {Math.round(result.strength * 100)}%
                  </span>
                </div>
                <div>
                  <span className="text-gray-400">Thay đổi uy tín: </span>
                  <span className="text-credibility-glow font-bold">
                    {Math.round(result.player_credibility_before * 100)}% &rarr;{" "}
                    {Math.round(result.player_credibility_after * 100)}%
                  </span>
                </div>
              </div>
            </div>

            <div className="flex justify-end">
              <button
                onClick={() => {
                  setResult(null);
                  closeCorrectionModal();
                }}
                className="px-5 py-2.5 bg-truth hover:bg-truth-glow text-white text-xs font-semibold rounded-xl transition-all"
              >
                Xác nhận & Tiếp tục điều tra
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handlePublish} className="p-6 space-y-5">
            {/* 1. Target Rumor Selection */}
            <div className="space-y-2">
              <label className="text-xs font-semibold text-gray-200 block">
                1. Chọn tin đồn muốn phản bác:
              </label>
              <div className="space-y-1.5">
                {targetRumors.map((r) => (
                  <label
                    key={r.id}
                    className={`flex items-center gap-3 p-2.5 rounded-xl border text-xs cursor-pointer transition-all ${
                      challengedClaimId === r.id
                        ? "bg-amber-950/40 border-amber-500/60 text-amber-200"
                        : "bg-surface border-white/5 text-gray-300 hover:border-white/20"
                    }`}
                  >
                    <input
                      type="radio"
                      name="rumor"
                      value={r.id}
                      checked={challengedClaimId === r.id}
                      onChange={() => setChallengedClaimId(r.id)}
                      className="accent-amber-500"
                    />
                    <span>{r.text}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* 2. Attach Evidence */}
            <div className="space-y-2">
              <label className="text-xs font-semibold text-gray-200 block">
                2. Đính kèm vật chứng xác thực (từ sổ tay):
              </label>
              {!notebook?.discovered_evidence || notebook.discovered_evidence.length === 0 ? (
                <div className="p-3 bg-black/40 rounded-xl border border-white/5 text-xs text-gray-500">
                  Bạn chưa có bằng chứng nào để đính kèm. Đính chính không có bằng chứng sẽ làm giảm uy tín của bạn!
                </div>
              ) : (
                <div className="max-h-36 overflow-y-auto space-y-1.5 p-1">
                  {notebook.discovered_evidence.map((ev) => (
                    <label
                      key={ev.id}
                      className={`flex items-center gap-3 p-2 rounded-lg border text-xs cursor-pointer transition-all ${
                        selectedEvIds.includes(ev.id)
                          ? "bg-truth/20 border-truth text-truth-glow"
                          : "bg-surface border-white/5 text-gray-300 hover:border-white/20"
                      }`}
                    >
                      <input
                        type="checkbox"
                        checked={selectedEvIds.includes(ev.id)}
                        onChange={() => handleToggleEvidence(ev.id)}
                        className="accent-truth"
                      />
                      <FileCheck2 className="w-3.5 h-3.5 flex-shrink-0" />
                      <span className="truncate">{ev.name}</span>
                    </label>
                  ))}
                </div>
              )}
            </div>

            {/* 3. Statement Message */}
            <div className="space-y-2">
              <label className="text-xs font-semibold text-gray-200 block">
                3. Tuyên bố kiểm chứng chính thức:
              </label>
              <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Ví dụ: Báo cáo y tế chính thức xác nhận chỉ có 4 ca ngạt khói nhẹ, hoàn toàn không có ca tử vong..."
                rows={3}
                required
                className="w-full bg-black/50 border border-white/10 rounded-xl p-3 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-amber-500"
              />
            </div>

            {/* Footer */}
            <div className="flex items-center justify-between pt-2 border-t border-white/10">
              <span className="text-[11px] text-gray-500 font-mono">Chi phí: +45 giây</span>
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={closeCorrectionModal}
                  className="px-4 py-2 rounded-xl text-xs text-gray-400 hover:text-white"
                >
                  Hủy
                </button>
                <button
                  type="submit"
                  disabled={submitting || !message.trim()}
                  className="px-5 py-2.5 bg-amber-600 hover:bg-amber-500 text-white text-xs font-semibold rounded-xl flex items-center gap-2 transition-all disabled:opacity-50 shadow-md"
                >
                  {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
                  Xuất bản ngay
                </button>
              </div>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
