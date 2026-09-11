"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { inspectLocation } from "@/services/api";
import { DISTRICT_LOCATIONS } from "@/components/map/DistrictMap";
import { EvidenceItem } from "@/types/game";
import {
  X,
  Search,
  FileCheck2,
  AlertCircle,
  MapPin,
  Loader2,
  Sparkles,
} from "lucide-react";

export function InspectModal() {
  const session = useGameStore((s) => s.session);
  const refreshState = useGameStore((s) => s.refreshState);
  const { isInspectModalOpen, selectedLocationId, closeInspectModal } = useUIStore();

  const [loading, setLoading] = useState(false);
  const [resultMsg, setResultMsg] = useState<string | null>(null);
  const [discovered, setDiscovered] = useState<EvidenceItem[]>([]);

  if (!isInspectModalOpen || !selectedLocationId || !session) return null;

  const loc = DISTRICT_LOCATIONS.find((l) => l.id === selectedLocationId) || {
    id: selectedLocationId,
    name: selectedLocationId,
    category: "public",
  };

  const handleInspect = async () => {
    setLoading(true);
    setResultMsg(null);
    try {
      const resp = await inspectLocation(session.id, selectedLocationId);
      setResultMsg(resp.message);
      setDiscovered(resp.discovered_evidence || []);
      await refreshState();
    } catch (err: any) {
      alert(`Khám nghiệm thất bại: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="glass-panel-elevated w-full max-w-lg rounded-2xl border border-white/10 shadow-2xl overflow-hidden animate-fade-in">
        {/* Header */}
        <div className="px-6 py-4 bg-surface/90 border-b border-white/10 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-blue-500/20 border border-blue-400/40 flex items-center justify-center">
              <MapPin className="w-5 h-5 text-blue-400" />
            </div>
            <div>
              <h3 className="font-bold text-gray-100 text-base">{loc.name}</h3>
              <p className="text-xs text-gray-400">Khám nghiệm hiện trường vật lý</p>
            </div>
          </div>
          <button
            onClick={closeInspectModal}
            className="p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-white/10 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-5">
          <p className="text-xs text-gray-300 leading-relaxed">
            Bạn đang có mặt tại <strong>{loc.name}</strong>. Khám nghiệm kỹ lưỡng khu vực này có thể
            giúp bạn thu thập vật chứng, dữ liệu camera an ninh, hoặc hồ sơ tài liệu chính thức.
          </p>

          <div className="p-3 bg-black/40 border border-white/5 rounded-xl text-xs text-gray-400 flex items-center justify-between font-mono">
            <span>Chi phí thời gian:</span>
            <span className="text-amber-400 font-semibold">+45 giây điều tra</span>
          </div>

          {resultMsg && (
            <div className="p-4 rounded-xl bg-surface border border-white/10 space-y-3">
              <p className="text-xs font-semibold text-gray-200">{resultMsg}</p>
              {discovered.length > 0 && (
                <div className="space-y-2">
                  {discovered.map((ev) => (
                    <div
                      key={ev.id}
                      className="p-3 bg-emerald-950/40 border border-emerald-500/40 rounded-lg text-xs space-y-1"
                    >
                      <div className="flex items-center gap-2 font-bold text-emerald-300">
                        <FileCheck2 className="w-4 h-4 text-emerald-400" />
                        {ev.name}
                      </div>
                      <p className="text-gray-300 text-[11px]">{ev.description}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          <div className="flex items-center justify-end gap-3 pt-2">
            <button
              onClick={closeInspectModal}
              className="px-4 py-2 rounded-xl text-xs text-gray-400 hover:text-white hover:bg-white/5 transition-colors"
            >
              Đóng
            </button>
            <button
              onClick={handleInspect}
              disabled={loading}
              className="px-5 py-2.5 bg-truth hover:bg-truth-glow text-white text-xs font-semibold rounded-xl flex items-center gap-2 transition-all disabled:opacity-50 glow-truth"
            >
              {loading ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Search className="w-4 h-4" />
              )}
              Kiểm tra khu vực
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
