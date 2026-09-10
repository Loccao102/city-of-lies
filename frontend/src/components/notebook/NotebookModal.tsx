"use client";

import React from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore, NotebookTab } from "@/store/uiStore";
import {
  X,
  BookOpen,
  FileCheck2,
  ListTodo,
  Users,
  Clock,
  ShieldCheck,
  CheckCircle2,
  XCircle,
  HelpCircle,
  MapPin,
} from "lucide-react";

export function NotebookModal() {
  const notebook = useGameStore((s) => s.notebook);
  const agents = useGameStore((s) => s.agents);
  const {
    isNotebookModalOpen,
    notebookTab,
    closeNotebookModal,
    setNotebookTab,
    openInterviewModal,
  } = useUIStore();

  if (!isNotebookModalOpen) return null;

  const tabs: { id: NotebookTab; label: string; icon: any }[] = [
    { id: "evidence", label: "Vật chứng & Tài liệu", icon: FileCheck2 },
    { id: "claims", label: "Mệnh đề đã biết", icon: ListTodo },
    { id: "people", label: "Danh bạ nhân vật (20)", icon: Users },
    { id: "timeline", label: "Dòng thời gian sự kiện", icon: Clock },
  ];

  return (
    <div className="fixed inset-0 bg-black/75 backdrop-blur-md z-50 flex items-center justify-center p-4">
      <div className="glass-panel-elevated w-full max-w-4xl rounded-2xl border border-white/10 shadow-2xl overflow-hidden flex flex-col h-[85vh] animate-fade-in">
        {/* Header */}
        <div className="px-6 py-4 bg-surface/90 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-truth/20 border border-truth/40 flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-truth-glow" />
            </div>
            <div>
              <h2 className="font-extrabold text-white text-base tracking-wide">
                SỔ TAY ĐIỀU TRA
              </h2>
              <p className="text-xs text-gray-400">
                Ghi chép chứng cứ, tuyên bố và hồ sơ nhân chứng vụ Riverside
              </p>
            </div>
          </div>
          <button
            onClick={closeNotebookModal}
            className="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation Tabs */}
        <div className="flex items-center px-6 bg-black/40 border-b border-white/10 gap-2 pt-2">
          {tabs.map((t) => {
            const Icon = t.icon;
            const active = notebookTab === t.id;
            return (
              <button
                key={t.id}
                onClick={() => setNotebookTab(t.id)}
                className={`flex items-center gap-2 px-4 py-2.5 text-xs font-semibold rounded-t-lg transition-all border-b-2 ${
                  active
                    ? "bg-surface text-truth-glow border-truth"
                    : "text-gray-400 hover:text-gray-200 border-transparent hover:bg-white/5"
                }`}
              >
                <Icon className="w-4 h-4" />
                {t.label}
              </button>
            );
          })}
        </div>

        {/* Tab Body */}
        <div className="flex-1 overflow-y-auto p-6">
          {/* 1. EVIDENCE TAB */}
          {notebookTab === "evidence" && (
            <div className="space-y-4">
              {!notebook?.discovered_evidence || notebook.discovered_evidence.length === 0 ? (
                <div className="text-center py-16 text-gray-500 text-xs flex flex-col items-center gap-2">
                  <FileCheck2 className="w-8 h-8 text-gray-600" />
                  <span>Chưa thu thập được bằng chứng vật lý nào.</span>
                  <span>Hãy khám nghiệm hiện trường hoặc phỏng vấn người có thẩm quyền.</span>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {notebook.discovered_evidence.map((ev) => {
                    let relTag = "Trung bình";
                    let relColor = "text-yellow-400 bg-yellow-950/40 border-yellow-800/40";
                    if (ev.reliability >= 0.85) {
                      relTag = "Độ tin cậy cao";
                      relColor = "text-emerald-400 bg-emerald-950/40 border-emerald-800/40";
                    } else if (ev.reliability < 0.6) {
                      relTag = "Độ tin cậy thấp";
                      relColor = "text-gray-400 bg-gray-900 border-gray-700";
                    }

                    return (
                      <div
                        key={ev.id}
                        className="p-4 rounded-xl bg-surface border border-white/10 flex flex-col justify-between gap-3 shadow-sm hover:border-white/20 transition-all"
                      >
                        <div className="space-y-2">
                          <div className="flex items-start justify-between gap-2">
                            <h4 className="font-bold text-gray-100 text-xs flex items-center gap-2">
                              <ShieldCheck className="w-4 h-4 text-truth-glow flex-shrink-0" />
                              {ev.name}
                            </h4>
                            <span
                              className={`text-[10px] font-mono px-2 py-0.5 rounded-full border ${relColor}`}
                            >
                              {relTag} ({Math.round(ev.reliability * 100)}%)
                            </span>
                          </div>
                          <p className="text-gray-300 text-xs leading-relaxed">
                            {ev.description}
                          </p>
                        </div>
                        <div className="flex items-center justify-between text-[11px] text-gray-500 font-mono pt-2 border-t border-white/5">
                          <span>Nguồn: {ev.source_category}</span>
                          <span>Vị trí: {ev.location_id}</span>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          )}

          {/* 2. CLAIMS TAB */}
          {notebookTab === "claims" && (
            <div className="space-y-3">
              {!notebook?.known_claims || notebook.known_claims.length === 0 ? (
                <div className="text-center py-16 text-gray-500 text-xs">
                  Chưa ghi nhận mệnh đề nào từ các cuộc trao đổi.
                </div>
              ) : (
                notebook.known_claims.map((claim) => {
                  let statusIcon = <HelpCircle className="w-4 h-4 text-gray-400" />;
                  let statusColor = "text-gray-400 bg-gray-800/40 border-gray-700";

                  if (claim.status === "Supported") {
                    statusIcon = <CheckCircle2 className="w-4 h-4 text-emerald-400" />;
                    statusColor = "text-emerald-300 bg-emerald-950/40 border-emerald-800/40";
                  } else if (claim.status === "Contradicted") {
                    statusIcon = <XCircle className="w-4 h-4 text-red-400" />;
                    statusColor = "text-red-300 bg-red-950/40 border-red-800/40";
                  }

                  return (
                    <div
                      key={claim.id}
                      className="p-3.5 rounded-xl bg-surface border border-white/10 flex items-center justify-between gap-4 text-xs"
                    >
                      <span className="text-gray-200 leading-relaxed">
                        {claim.display_text}
                      </span>
                      <div
                        className={`flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-[11px] font-mono flex-shrink-0 ${statusColor}`}
                      >
                        {statusIcon}
                        <span>{claim.status}</span>
                      </div>
                    </div>
                  );
                })
              )}
            </div>
          )}

          {/* 3. PEOPLE TAB */}
          {notebookTab === "people" && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {agents.map((ag) => (
                <div
                  key={ag.id}
                  className="p-3.5 rounded-xl bg-surface border border-white/10 flex items-center justify-between gap-3 text-xs hover:border-white/20 transition-all"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="font-bold text-gray-200">{ag.name}</span>
                      {ag.interviewed ? (
                        <span className="text-[10px] px-2 py-0.5 rounded bg-truth/20 text-truth-glow border border-truth/40">
                          Đã phỏng vấn
                        </span>
                      ) : (
                        <span className="text-[10px] px-2 py-0.5 rounded bg-white/5 text-gray-500">
                          Chưa gặp
                        </span>
                      )}
                    </div>
                    <p className="text-gray-400 text-[11px]">{ag.role}</p>
                    <div className="flex items-center gap-1 text-[10px] text-gray-500 font-mono">
                      <MapPin className="w-3 h-3 text-gray-400" />
                      <span>{ag.location_id}</span>
                    </div>
                  </div>

                  <button
                    onClick={() => {
                      closeNotebookModal();
                      openInterviewModal(ag.id);
                    }}
                    className="px-3 py-1.5 rounded-lg bg-surface-elevated hover:bg-surface-border text-gray-300 hover:text-white border border-white/10 text-xs transition-colors flex-shrink-0"
                  >
                    Hỏi chuyện
                  </button>
                </div>
              ))}
            </div>
          )}

          {/* 4. TIMELINE TAB */}
          {notebookTab === "timeline" && (
            <div className="space-y-3">
              {!notebook?.timeline_events || notebook.timeline_events.length === 0 ? (
                <div className="text-center py-16 text-gray-500 text-xs">
                  Chưa có mốc thời gian nào được ghi nhận.
                </div>
              ) : (
                notebook.timeline_events.map((tEv, i) => (
                  <div
                    key={i}
                    className="flex items-start gap-4 p-3 rounded-xl bg-surface border border-white/5 text-xs"
                  >
                    <div className="px-2.5 py-1 rounded bg-black/50 border border-white/10 font-mono text-[11px] text-truth-glow flex-shrink-0">
                      {Math.floor(tEv.game_second / 60)}m {tEv.game_second % 60}s
                    </div>
                    <div className="space-y-1">
                      <h5 className="font-semibold text-gray-200">{tEv.headline}</h5>
                      <p className="text-gray-400 text-[11px] leading-relaxed">
                        {tEv.description}
                      </p>
                    </div>
                  </div>
                ))
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
