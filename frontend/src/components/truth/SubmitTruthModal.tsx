"use client";

import React, { useState, useEffect } from "react";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";
import { submitTruth } from "@/services/api";
import { SubmitTruthRequest } from "@/types/game";
import {
  X,
  Flag,
  AlertCircle,
  Loader2,
} from "lucide-react";

interface ScenarioConfig {
  eventTypes: { value: string; label: string }[];
  causeCategories: { value: string; label: string }[];
  hazardLabel: string;
  noHazardLabel: string;
  hasHazardLabel: string;
}

const SCENARIO_CONFIGS: Record<string, ScenarioConfig> = {
  "riverside-factory": {
    eventTypes: [
      { value: "fire", label: "Cháy cục bộ (Fire)" },
      { value: "chemical_explosion", label: "Nổ hóa chất (Chemical Explosion)" },
      { value: "gas_explosion", label: "Nổ khí gas (Gas Explosion)" },
      { value: "unknown", label: "Chưa xác định" },
    ],
    causeCategories: [
      { value: "electrical_fault", label: "Sự cố chập điện tủ phân phối DB-4" },
      { value: "chemical_reaction", label: "Phản ứng hóa chất nguy hiểm" },
      { value: "arson", label: "Cố ý phóng hỏa" },
      { value: "unknown", label: "Không rõ" },
    ],
    hazardLabel: "Có vụ nổ lớn phá hủy cấu trúc không?",
    noHazardLabel: "Không có nổ lớn (Chỉ có tiếng nổ chập aptomat)",
    hasHazardLabel: "Có nổ lớn phá hủy nhà máy",
  },
  "metro-hospital-outbreak": {
    eventTypes: [
      { value: "food_poisoning", label: "Ngộ độc thực phẩm tập thể" },
      { value: "viral_outbreak", label: "Bùng phát virus nguy hiểm" },
      { value: "chemical_poisoning", label: "Đầu độc hóa chất" },
      { value: "unknown", label: "Chưa xác định" },
    ],
    causeCategories: [
      { value: "bacterial_toxin", label: "Độc tố histamine do bảo quản cá ngừ sai quy chuẩn" },
      { value: "lab_leak", label: "Rò rỉ mầm bệnh từ phòng thí nghiệm vi sinh" },
      { value: "deliberate_sabotage", label: "Kẻ xấu cố ý đầu độc thức ăn" },
      { value: "unknown", label: "Không rõ" },
    ],
    hazardLabel: "Có sự cố an toàn sinh học cấp cao / thảm họa lây nhiễm không?",
    noHazardLabel: "Không có rò rỉ virus sinh học",
    hasHazardLabel: "Có rò rỉ virus mầm bệnh nguy hiểm",
  },
  "midtown-bank-run": {
    eventTypes: [
      { value: "tech_outage", label: "Gián đoạn kỹ thuật hệ thống Core Banking" },
      { value: "bank_insolvency", label: "Ngân hàng mất thanh khoản / phá sản" },
      { value: "cyber_heist", label: "Vụ trộm mạng quy mô lớn" },
      { value: "unknown", label: "Chưa xác định" },
    ],
    causeCategories: [
      { value: "core_banking_deadlock", label: "Xung đột deadlock cơ sở dữ liệu khi nâng cấp bản vá v4.12" },
      { value: "embezzlement", label: "Ban lãnh đạo tẩu tán tài sản" },
      { value: "hacker_ransomware", label: "Mã độc tống tiền khóa dữ liệu" },
      { value: "unknown", label: "Không rõ" },
    ],
    hazardLabel: "Có sự cố mất trắng tiền gửi / kho quỹ bị cướp phá không?",
    noHazardLabel: "Không có mất mát tài sản (Hầm kho an toàn 100%)",
    hasHazardLabel: "Có mất tiền / ban lãnh đạo tháo chạy",
  },
  "subway-line3-standstill": {
    eventTypes: [
      { value: "signal_failure", label: "Sự cố chập tín hiệu kích hoạt dừng tàu an toàn" },
      { value: "terrorist_attack", label: "Tấn công khủng bố bom khí độc" },
      { value: "tunnel_collapse", label: "Sập sụt lún kết cấu hầm ngầm" },
      { value: "unknown", label: "Chưa xác định" },
    ],
    causeCategories: [
      { value: "short_circuit_cable", label: "Đoản mạch cáp tín hiệu hộp 14 do ẩm mốc" },
      { value: "chemical_gas_device", label: "Kích nổ thiết bị hơi ngạt độc hại" },
      { value: "train_collision", label: "Va chạm phương tiện khác" },
      { value: "unknown", label: "Không rõ" },
    ],
    hazardLabel: "Có vụ nổ phá hủy toa tàu hay sập hầm ngầm không?",
    noHazardLabel: "Không có nổ bom hay sập hầm (Vỏ hầm nguyên vẹn)",
    hasHazardLabel: "Có tấn công nổ bom / sập hầm ngầm",
  },
  "city-water-panic": {
    eventTypes: [
      { value: "pipe_rupture", label: "Sự cố van áp lực làm sục cặn oxit sắt gây đục tạm thời" },
      { value: "chemical_poisoning", label: "Nguồn nước bị đầu độc bằng xyanua công nghiệp" },
      { value: "algal_bloom", label: "Bùng phát tảo độc ô nhiễm nguồn nước" },
      { value: "unknown", label: "Chưa xác định" },
    ],
    causeCategories: [
      { value: "mineral_sediment_surge", label: "Dòng xoáy sục lớp bùn oxit sắt tự nhiên trong thành ống cũ" },
      { value: "deliberate_sabotage", label: "Kẻ xấu lén đổ hóa chất độc hại vào bể chứa" },
      { value: "industrial_effluent", label: "Nhà máy xả trộm nước thải chưa qua xử lý" },
      { value: "unknown", label: "Không rõ" },
    ],
    hazardLabel: "Có thảm họa nhiễm độc diện rộng / nước nhiễm độc chết người không?",
    noHazardLabel: "Không có độc tố (Xét nghiệm âm tính 100%)",
    hasHazardLabel: "Có chất độc xyanua chết người",
  },
};

export function SubmitTruthModal() {
  const session = useGameStore((s) => s.session);
  const refreshState = useGameStore((s) => s.refreshState);
  const { isSubmitTruthModalOpen, closeSubmitTruthModal, openAARModal } = useUIStore();

  const cfg = SCENARIO_CONFIGS[session?.scenario_id || "riverside-factory"] || SCENARIO_CONFIGS["riverside-factory"];

  const [eventType, setEventType] = useState(cfg.eventTypes[0]?.value || "fire");
  const [causeCategory, setCauseCategory] = useState(cfg.causeCategories[0]?.value || "electrical_fault");
  const [majorExplosion, setMajorExplosion] = useState(false);
  const [fatalities, setFatalities] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<string | null>(null);

  useEffect(() => {
    if (session?.scenario_id && SCENARIO_CONFIGS[session.scenario_id]) {
      const c = SCENARIO_CONFIGS[session.scenario_id];
      setEventType(c.eventTypes[0]?.value || "fire");
      setCauseCategory(c.causeCategories[0]?.value || "electrical_fault");
      setMajorExplosion(false);
      setFatalities(0);
      setFeedback(null);
    }
  }, [session?.scenario_id, isSubmitTruthModalOpen]);

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
                Xác lập Ground Truth cho {session.scenario_title || "vụ án"}
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
              {cfg.eventTypes.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          {/* 2. Cause Category */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200">
              Nguyên nhân chính (Cause Category):
            </label>
            <select
              value={causeCategory}
              onChange={(e) => setCauseCategory(e.target.value)}
              className="w-full bg-black/50 border border-white/10 rounded-xl p-2.5 text-xs text-white focus:outline-none focus:border-truth"
            >
              {cfg.causeCategories.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          {/* 3. Major Hazard / Explosion */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-200 block">
              {cfg.hazardLabel}
            </label>
            <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6">
              <label className="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
                <input
                  type="radio"
                  name="explosion"
                  checked={majorExplosion === false}
                  onChange={() => setMajorExplosion(false)}
                  className="accent-truth"
                />
                <span>{cfg.noHazardLabel}</span>
              </label>
              <label className="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
                <input
                  type="radio"
                  name="explosion"
                  checked={majorExplosion === true}
                  onChange={() => setMajorExplosion(true)}
                  className="accent-truth"
                />
                <span>{cfg.hasHazardLabel}</span>
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
              max={200}
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
