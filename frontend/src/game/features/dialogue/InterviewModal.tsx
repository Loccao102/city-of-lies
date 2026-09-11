"use client";

import React, { useState, useEffect } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { getAgentDetail, interviewAgent } from "@/services/api";
import { AgentDetail, DialogueResponse } from "@/types/game";
import {
  X,
  MessageSquare,
  Sparkles,
  Send,
  AlertCircle,
  FileCheck2,
  Loader2,
  User,
} from "lucide-react";

export function InterviewModal() {
  const session = useGameStore((s) => s.session);
  const refreshState = useGameStore((s) => s.refreshState);
  const { isInterviewModalOpen, selectedAgentId, closeInterviewModal } = useUIStore();

  const [agent, setAgent] = useState<AgentDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [query, setQuery] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [history, setHistory] = useState<{ q: string; r: DialogueResponse }[]>([]);

  useEffect(() => {
    if (!isInterviewModalOpen || !selectedAgentId || !session) return;
    setLoading(true);
    setHistory([]);
    getAgentDetail(session.id, selectedAgentId)
      .then((data) => setAgent(data))
      .catch((e) => console.error("Failed to load agent:", e))
      .finally(() => setLoading(false));
  }, [isInterviewModalOpen, selectedAgentId, session]);

  if (!isInterviewModalOpen || !agent) return null;

  const quickQuestions = [
    "Lúc xảy ra sự cố, anh/chị trực tiếp chứng kiến những gì?",
    "Anh/chị có thấy khói hay nghe tiếng nổ từ đâu không?",
    "Có ai bị thương nặng hay có ca tử vong nào không?",
    "Anh/chị có tài liệu hay bằng chứng gì về sự cố này không?",
  ];

  const handleAsk = async (questionText: string) => {
    if (!session || !questionText.trim() || submitting) return;
    setSubmitting(true);
    try {
      const resp = await interviewAgent(session.id, agent.id, questionText.trim());
      setHistory((prev) => [...prev, { q: questionText.trim(), r: resp }]);
      setQuery("");
      // Refresh notebook in case evidence was unlocked
      if (resp.revealed_evidence_ids?.length > 0) {
        await refreshState();
      }
    } catch (err: any) {
      alert(`Phỏng vấn thất bại: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="glass-panel-elevated w-full max-w-2xl rounded-2xl border border-white/10 shadow-2xl overflow-hidden flex flex-col max-h-[85vh] animate-fade-in">
        {/* Header */}
        <div className="px-6 py-4 bg-surface/90 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-truth/20 border border-truth/40 flex items-center justify-center">
              <User className="w-5 h-5 text-truth-glow" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h3 className="font-bold text-gray-100 text-base">{agent.name}</h3>
                <span className="text-[11px] px-2 py-0.5 rounded-full bg-white/10 font-mono text-gray-300">
                  {agent.role}
                </span>
              </div>
              <p className="text-xs text-gray-400 mt-0.5">{agent.personality}</p>
            </div>
          </div>
          <button
            onClick={closeInterviewModal}
            className="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Dialogue History */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4">
          {history.length === 0 ? (
            <div className="py-8 text-center text-gray-400 text-xs flex flex-col items-center gap-2">
              <MessageSquare className="w-8 h-8 text-gray-600" />
              <span>Hãy đặt câu hỏi để tìm kiếm manh mối hoặc bằng chứng từ nhân vật.</span>
              <span className="text-[11px] text-gray-500">
                (Lưu ý: Mỗi câu hỏi sẽ tiêu tốn 30 giây thời gian điều tra)
              </span>
            </div>
          ) : (
            history.map((item, idx) => (
              <div key={idx} className="space-y-3">
                {/* Player Query */}
                <div className="flex justify-end">
                  <div className="max-w-[80%] bg-truth/20 border border-truth/40 rounded-xl px-4 py-2.5 text-xs text-gray-100 shadow-sm">
                    {item.q}
                  </div>
                </div>

                {/* NPC Response */}
                <div className="flex justify-start">
                  <div className="max-w-[85%] bg-surface border border-white/10 rounded-xl p-4 text-xs text-gray-200 shadow-sm space-y-2">
                    <p className="leading-relaxed">{item.r.utterance}</p>

                    <div className="flex flex-wrap items-center gap-3 pt-1 text-[11px] text-gray-400 border-t border-white/5 font-mono">
                      <span>Cảm xúc: {item.r.emotion}</span>
                      <span>•</span>
                      <span>Độ chắc chắn: {Math.round(item.r.certainty * 100)}%</span>
                    </div>

                    {/* Unlocked Evidence notification */}
                    {item.r.revealed_evidence_ids && item.r.revealed_evidence_ids.length > 0 && (
                      <div className="mt-2 p-2.5 rounded-lg bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 flex items-center gap-2 text-xs">
                        <FileCheck2 className="w-4 h-4 text-emerald-400 flex-shrink-0" />
                        <span>
                          <strong>Mở khóa bằng chứng mới!</strong> Đã ghi nhận vào sổ tay điều tra.
                        </span>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Suggested Quick Questions */}
        <div className="px-6 py-3 bg-black/40 border-t border-white/5 flex flex-wrap gap-2">
          {quickQuestions.map((qText, i) => (
            <button
              key={i}
              onClick={() => handleAsk(qText)}
              disabled={submitting}
              className="text-[11px] px-3 py-1.5 rounded-md bg-surface-elevated hover:bg-surface-border text-gray-300 border border-white/10 transition-colors text-left"
            >
              {qText}
            </button>
          ))}
        </div>

        {/* Input Dock */}
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleAsk(query);
          }}
          className="p-4 bg-surface border-t border-white/10 flex items-center gap-3"
        >
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Nhập câu hỏi bằng tiếng Việt..."
            disabled={submitting}
            className="flex-1 bg-black/50 border border-white/10 rounded-xl px-4 py-2.5 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-truth"
          />
          <button
            type="submit"
            disabled={submitting || !query.trim()}
            className="px-4 py-2.5 bg-truth hover:bg-truth-glow text-white text-xs font-semibold rounded-xl flex items-center gap-2 transition-all disabled:opacity-40"
          >
            {submitting ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Send className="w-4 h-4" />
            )}
            Hỏi
          </button>
        </form>
      </div>
    </div>
  );
}
