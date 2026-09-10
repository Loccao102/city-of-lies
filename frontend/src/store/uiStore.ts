import { create } from "zustand";
import { TruthSubmissionResult } from "@/types/game";

export type NotebookTab = "evidence" | "claims" | "people" | "timeline";

interface UIStore {
  selectedAgentId: string | null;
  selectedLocationId: string | null;

  isInterviewModalOpen: boolean;
  isInspectModalOpen: boolean;
  isNotebookModalOpen: boolean;
  notebookTab: NotebookTab;
  isCorrectionModalOpen: boolean;
  isSubmitTruthModalOpen: boolean;
  isAARModalOpen: boolean;
  aarResult: TruthSubmissionResult | null;

  setSelectedAgentId: (id: string | null) => void;
  setSelectedLocationId: (id: string | null) => void;
  openInterviewModal: (agentId: string) => void;
  closeInterviewModal: () => void;
  openInspectModal: (locationId: string) => void;
  closeInspectModal: () => void;
  openNotebookModal: (tab?: NotebookTab) => void;
  closeNotebookModal: () => void;
  setNotebookTab: (tab: NotebookTab) => void;
  openCorrectionModal: () => void;
  closeCorrectionModal: () => void;
  openSubmitTruthModal: () => void;
  closeSubmitTruthModal: () => void;
  openAARModal: (result?: TruthSubmissionResult) => void;
  closeAARModal: () => void;
}

export const useUIStore = create<UIStore>((set) => ({
  selectedAgentId: null,
  selectedLocationId: null,

  isInterviewModalOpen: false,
  isInspectModalOpen: false,
  isNotebookModalOpen: false,
  notebookTab: "evidence",
  isCorrectionModalOpen: false,
  isSubmitTruthModalOpen: false,
  isAARModalOpen: false,
  aarResult: null,

  setSelectedAgentId: (id) => set({ selectedAgentId: id }),
  setSelectedLocationId: (id) => set({ selectedLocationId: id }),

  openInterviewModal: (agentId) =>
    set({ selectedAgentId: agentId, isInterviewModalOpen: true }),
  closeInterviewModal: () => set({ isInterviewModalOpen: false }),

  openInspectModal: (locationId) =>
    set({ selectedLocationId: locationId, isInspectModalOpen: true }),
  closeInspectModal: () => set({ isInspectModalOpen: false }),

  openNotebookModal: (tab = "evidence") =>
    set({ isNotebookModalOpen: true, notebookTab: tab }),
  closeNotebookModal: () => set({ isNotebookModalOpen: false }),
  setNotebookTab: (tab) => set({ notebookTab: tab }),

  openCorrectionModal: () => set({ isCorrectionModalOpen: true }),
  closeCorrectionModal: () => set({ isCorrectionModalOpen: false }),

  openSubmitTruthModal: () => set({ isSubmitTruthModalOpen: true }),
  closeSubmitTruthModal: () => set({ isSubmitTruthModalOpen: false }),

  openAARModal: (result) =>
    set({ isAARModalOpen: true, aarResult: result || null }),
  closeAARModal: () => set({ isAARModalOpen: false }),
}));
