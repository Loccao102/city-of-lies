"use client";

import React, { useState } from "react";
import { useGameStore } from "@/store/gameStore";
import { createSession } from "@/services/api";
import { TopBar } from "@/components/hud/TopBar";
import { EventFeed } from "@/components/hud/EventFeed";
import { DistrictMap } from "@/components/map/DistrictMap";
import { InterviewModal } from "@/components/interview/InterviewModal";
import { InspectModal } from "@/components/inspection/InspectModal";
import { NotebookModal } from "@/components/notebook/NotebookModal";
import { CorrectionModal } from "@/components/correction/CorrectionModal";
import { SubmitTruthModal } from "@/components/truth/SubmitTruthModal";
import { AfterActionReport } from "@/components/aar/AfterActionReport";
import {
  Sparkles,
  ShieldAlert,
  Flame,
  Search,
  MessageSquare,
  Send,
  Flag,
  Loader2,
  Dice5,
} from "lucide-react";

export default function Home() {
  const session = useGameStore((s) => s.session);
  const setSession = useGameStore((s) => s.setSession);
  const connectWS = useGameStore((s) => s.connectWS);
  const refreshState = useGameStore((s) => s.refreshState);

  const [loading, setLoading] = useState(false);
  const [seedInput, setSeedInput] = useState<string>("424242");

  const handleStartGame = async () => {
    setLoading(true);
    try {
      const seedNum = seedInput ? parseInt(seedInput, 10) : undefined;
      const newSession = await createSession("riverside-factory", seedNum);
      setSession(newSession);
      connectWS(newSession.id);
      await refreshState();
    } catch (err: any) {
      alert(`Khởi tạo game thất bại: ${err.message}. Hãy đảm bảo Backend Go đang chạy tại port 8080.`);
    } finally {
      setLoading(false);
    }
  };

  // If session is active, render Game UI
  if (session) {
    return (
      <main className="relative w-screen h-screen overflow-hidden bg-background">
        <TopBar />
        <div className="w-full h-full pt-16">
          <DistrictMap />
        </div>
        <EventFeed />

        {/* Modals */}
        <InterviewModal />
        <InspectModal />
        <NotebookModal />
        <CorrectionModal />
        <SubmitTruthModal />
        <AfterActionReport />
      </main>
    );
  }

  // Otherwise, render Atmospheric Landing Screen
  return (
    <main className="relative min-h-screen flex flex-col items-center justify-center p-6 bg-gradient-to-b from-background via-surface to-background">
      {/* Background Decor */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(49,151,149,0.08)_0%,transparent_70%)] pointer-events-none" />

      <div className="relative z-10 w-full max-w-4xl space-y-12 text-center">
        {/* Header Badge */}
        <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-white/5 border border-white/10 text-xs font-mono text-truth-glow">
          <Sparkles className="w-3.5 h-3.5 text-truth-glow" />
          <span>SIMULATION ENGINE AUTHORITATIVE • ZERO-LLM PLAYABLE</span>
        </div>

        {/* Hero Title */}
        <div className="space-y-4">
          <h1 className="text-5xl md:text-7xl font-black tracking-tight text-white">
            CITY OF <span className="text-transparent bg-clip-text bg-gradient-to-r from-red-500 via-amber-500 to-truth-glow">LIES</span>
          </h1>
          <p className="text-lg md:text-xl text-gray-300 max-w-2xl mx-auto font-light leading-relaxed">
            Tìm ra sự thật trước khi lời dối trá trở thành chân lý của cả thành phố.
          </p>
        </div>

        {/* Scenario Card */}
        <div className="glass-panel-elevated max-w-2xl mx-auto rounded-3xl border border-white/10 p-8 shadow-2xl text-left space-y-6">
          <div className="flex items-center justify-between border-b border-white/10 pb-4">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-2xl bg-red-500/20 border border-red-500/40 flex items-center justify-center">
                <Flame className="w-5 h-5 text-red-400" />
              </div>
              <div>
                <h3 className="font-bold text-white text-base">Hồ Sơ: Sự Cố Nhà Máy Riverside</h3>
                <span className="text-xs text-gray-400 font-mono">Bối cảnh: 17:58 — 20 Nhân vật tự trị</span>
              </div>
            </div>
            <span className="text-[11px] font-mono px-3 py-1 rounded-full bg-red-950/60 text-red-300 border border-red-800/40">
              Nguy cơ: Nổ hóa chất & Giấu xác
            </span>
          </div>

          <p className="text-xs text-gray-300 leading-relaxed">
            Một vụ chập tủ điện DB-4 gây khói tại nhà kho Riverside. Trong khi chỉ có một vài người bị ngạt khói nhẹ và không có ai tử vong, tin đồn độc hại đang bùng phát từ lời kể truyền tai, bài đăng mạng xã hội của blogger và livestream kích động. 
            Nếu không có sự can thiệp kịp thời, <strong>75% cư dân</strong> sẽ tin vào lời nói dối và bạn sẽ thua cuộc.
          </p>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs font-mono">
            <div className="p-3 bg-black/40 rounded-xl border border-white/5 space-y-1">
              <span className="text-gray-500 text-[10px] block">THỜI GIAN ĐẦU HÀNG</span>
              <span className="text-red-400 font-bold">12 - 18 Phút</span>
            </div>
            <div className="p-3 bg-black/40 rounded-xl border border-white/5 space-y-1">
              <span className="text-gray-500 text-[10px] block">NGƯỠNG THUA</span>
              <span className="text-amber-400 font-bold">75% Tin giả</span>
            </div>
            <div className="p-3 bg-black/40 rounded-xl border border-white/5 space-y-1">
              <span className="text-gray-500 text-[10px] block">ĐỘ MẠNH CHỨNG CỨ</span>
              <span className="text-truth font-bold">&ge; 70% Để thắng</span>
            </div>
            <div className="p-3 bg-black/40 rounded-xl border border-white/5 space-y-1">
              <span className="text-gray-500 text-[10px] block">NHÂN VẬT AI</span>
              <span className="text-gray-200 font-bold">20 Tự trị</span>
            </div>
          </div>

          {/* Seed Input & Start Action */}
          <div className="flex flex-col sm:flex-row items-center gap-4 pt-2">
            <div className="flex items-center gap-2 px-3 py-2 bg-black/50 border border-white/10 rounded-xl text-xs w-full sm:w-auto">
              <Dice5 className="w-4 h-4 text-gray-400" />
              <span className="text-gray-400">Seed:</span>
              <input
                type="text"
                value={seedInput}
                onChange={(e) => setSeedInput(e.target.value)}
                placeholder="424242"
                className="w-20 bg-transparent text-white focus:outline-none font-mono font-semibold"
              />
            </div>

            <button
              onClick={handleStartGame}
              disabled={loading}
              className="w-full sm:flex-1 py-3.5 bg-truth hover:bg-truth-glow text-white text-sm font-bold rounded-2xl flex items-center justify-center gap-2 transition-all shadow-xl glow-truth disabled:opacity-50"
            >
              {loading ? (
                <Loader2 className="w-5 h-5 animate-spin" />
              ) : (
                <Sparkles className="w-5 h-5" />
              )}
              BẮT ĐẦU ĐIỀU TRA NGAY
            </button>
          </div>
        </div>

        {/* 3 Steps Guide */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-3xl mx-auto text-left">
          <div className="p-5 rounded-2xl bg-surface/60 border border-white/5 space-y-2">
            <div className="w-8 h-8 rounded-xl bg-blue-500/20 text-blue-400 flex items-center justify-center font-bold text-xs font-mono">
              01
            </div>
            <h4 className="font-semibold text-gray-200 text-sm">Thu thập bằng chứng</h4>
            <p className="text-xs text-gray-400 leading-relaxed">
              Phỏng vấn nhân chứng trực tiếp, kiểm tra hiện trường kho cháy, camera an ninh và báo cáo cứu hỏa.
            </p>
          </div>

          <div className="p-5 rounded-2xl bg-surface/60 border border-white/5 space-y-2">
            <div className="w-8 h-8 rounded-xl bg-amber-500/20 text-amber-400 flex items-center justify-center font-bold text-xs font-mono">
              02
            </div>
            <h4 className="font-semibold text-gray-200 text-sm">Xuất bản đính chính</h4>
            <p className="text-xs text-gray-400 leading-relaxed">
              Dùng tài liệu y tế và biên bản kỹ thuật phản bác tin đồn nổ hóa chất & chết người trước khi lan truyền rộng.
            </p>
          </div>

          <div className="p-5 rounded-2xl bg-surface/60 border border-white/5 space-y-2">
            <div className="w-8 h-8 rounded-xl bg-emerald-500/20 text-emerald-400 flex items-center justify-center font-bold text-xs font-mono">
              03
            </div>
            <h4 className="font-semibold text-gray-200 text-sm">Xác lập Sự thật</h4>
            <p className="text-xs text-gray-400 leading-relaxed">
              Gửi kết luận cuối cùng với &ge;70% chứng cứ xác đáng để kết thúc vụ án và bảo vệ công lý cho thành phố.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
