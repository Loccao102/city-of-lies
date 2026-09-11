import * as THREE from "three";
import type { AgentSummary } from "@/types/game";

function getRoleColor(role: string): number {
  const normalized = (role || "").toLowerCase();
  if (normalized.includes("cứu hỏa") || normalized.includes("lính cứu hỏa") || normalized.includes("fire")) return 0xe53e3e;
  if (normalized.includes("bác sĩ") || normalized.includes("y tá") || normalized.includes("y tế") || normalized.includes("doctor")) return 0x319795;
  if (normalized.includes("phóng viên") || normalized.includes("blogger") || normalized.includes("influencer") || normalized.includes("báo")) return 0x9f7aea;
  if (normalized.includes("bảo vệ") || normalized.includes("an ninh") || normalized.includes("công an") || normalized.includes("security")) return 0x3182ce;
  if (normalized.includes("công nhân") || normalized.includes("quản lý") || normalized.includes("worker") || normalized.includes("kho")) return 0xd69e2e;
  return 0x4fd1c5;
}

function disposeObjectTree(root: THREE.Object3D) {
  root.traverse((child) => {
    const disposable = child as THREE.Object3D & {
      geometry?: THREE.BufferGeometry;
      material?: THREE.Material | THREE.Material[];
    };
    disposable.geometry?.dispose();
    const material = disposable.material;
    if (Array.isArray(material)) material.forEach((entry) => entry.dispose());
    else material?.dispose();
  });
}

export class AgentSystem {
  private readonly root = new THREE.Group();
  private interactiveObjects: THREE.Object3D[] = [];

  constructor(
    private readonly scene: THREE.Scene,
    private readonly locationMap: Map<string, THREE.Vector3>
  ) {
    this.root.name = "npc-agents";
    this.scene.add(this.root);
  }

  setAgents(agents: AgentSummary[]) {
    const previousChildren = [...this.root.children];
    previousChildren.forEach((child) => {
      this.root.remove(child);
      disposeObjectTree(child);
    });
    this.interactiveObjects = [];

    const slotByLocation = new Map<string, number>();

    agents.forEach((agent) => {
      const locationPosition = this.locationMap.get(agent.location_id) ?? new THREE.Vector3();
      const slot = slotByLocation.get(agent.location_id) ?? 0;
      slotByLocation.set(agent.location_id, slot + 1);
      const angle = (slot * Math.PI * 2) / 6;
      const radius = 6;

      const group = new THREE.Group();
      group.position.set(
        locationPosition.x + Math.cos(angle) * radius,
        0,
        locationPosition.z + Math.sin(angle) * radius
      );
      group.userData = {
        type: "agent",
        agentId: agent.id,
        name: agent.name,
        role: agent.role,
      };

      const body = new THREE.Mesh(
        new THREE.CylinderGeometry(0.5, 0.5, 2.2, 12),
        new THREE.MeshStandardMaterial({
          color: agent.interviewed ? 0x4fd1c5 : getRoleColor(agent.role),
          roughness: 0.4,
          metalness: 0.1,
        })
      );
      body.position.y = 1.1;
      body.castShadow = true;
      group.add(body);

      const head = new THREE.Mesh(
        new THREE.SphereGeometry(0.55, 12, 12),
        new THREE.MeshStandardMaterial({ color: 0xffffff })
      );
      head.position.y = 2.6;
      head.castShadow = true;
      group.add(head);

      this.root.add(group);
      this.interactiveObjects.push(body);
    });
  }

  update(elapsed: number) {
    this.root.children.forEach((group) => {
      group.position.y = Math.sin(elapsed * 2 + group.position.x) * 0.15;
    });
  }

  getInteractiveObjects() {
    return this.interactiveObjects;
  }

  dispose() {
    this.scene.remove(this.root);
    disposeObjectTree(this.root);
    this.interactiveObjects = [];
  }
}
