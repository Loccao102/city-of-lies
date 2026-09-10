"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { Activity, ChevronDown, ChevronUp, Radio, AlertCircle, CheckCircle2 } from "lucide-react";

export function EventFeed() {
  const events = useGameStore((s) => s.events);
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="fixed bottom-6 left-6 w-96 glass-panel rounded-xl border border-white/10 z-20 shadow-2xl overflow-hidden transition-all duration-300">
      {/* Header */}
      <div
        className="flex items-center justify-between px-4 py-2.5 bg-black/40 border-b border-white/10 cursor-pointer"
        onClick={() => setCollapsed(!collapsed)}
      >
        <div className="flex items-center gap-2">
          <Activity className="w-4 h-4 text-truth-glow animate-pulse" />
          <span className="text-xs font-semibold text-gray-200 uppercase tracking-wider">
            Nhật ký sự kiện thành phố
          </span>
          <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-white/10 font-mono text-gray-400">
            {events.length}
          </span>
        </div>
        <button className="text-gray-400 hover:text-white">
          {collapsed ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
        </button>
      </div>

      {/* Feed List */}
      {!collapsed && (
        <div className="max-h-60 overflow-y-auto p-3 flex flex-col gap-2.5">
          {events.length === 0 ? (
            <div className="text-xs text-gray-500 text-center py-6">
              Chưa có diễn biến mới nào được ghi nhận...
            </div>
          ) : (
            events.slice(0, 20).map((evt, idx) => {
              let icon = <Radio className="w-3.5 h-3.5 text-blue-400" />;
              let badgeBg = "bg-blue-950/40 text-blue-300 border-blue-800/40";
              let title = evt.event_type;
              let desc = "";

              if (evt.event_type === "rumor.visualized") {
                icon = <AlertCircle className="w-3.5 h-3.5 text-amber-400" />;
                badgeBg = "bg-amber-950/40 text-amber-300 border-amber-800/40";
                title = "Tin đồn lan truyền";
                desc = `${evt.payload.speaker_name} trao đổi với ${evt.payload.listener_name}`;
              } else if (evt.event_type === "rumor.shared") {
                icon = <Radio className="w-3.5 h-3.5 text-purple-400" />;
                badgeBg = "bg-purple-950/40 text-purple-300 border-purple-800/40";
                title = evt.payload.headline || "Diễn biến công cộng";
                desc = evt.payload.description || "";
              } else if (evt.event_type === "correction.published") {
                icon = <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400" />;
                badgeBg = "bg-emerald-950/40 text-emerald-300 border-emerald-800/40";
                title = "Đính chính xuất bản";
                desc = evt.payload.message || "";
              } else if (evt.event_type === "conversation.agent_message") {
                title = `${evt.payload.name} trả lời`;
                desc = evt.payload.utterance || "";
              }

              return (
                <div
                  key={evt.id || idx}
                  className="p-2 rounded-lg bg-surface/80 border border-white/5 flex flex-col gap-1 text-xs hover:border-white/20 transition-all"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-1.5 font-medium">
                      {icon}
                      <span className="text-gray-200">{title}</span>
                    </div>
                    {evt.game_second !== undefined && (
                      <span className="text-[10px] font-mono text-gray-500">
                        {Math.floor(evt.game_second / 60)}m {evt.game_second % 60}s
                      </span>
                    )}
                  </div>
                  {desc && <p className="text-gray-400 text-[11px] leading-relaxed pl-5">{desc}</p>}
                </div>
              );
            })
          )}
        </div>
      )}
    </div>
  );
}
