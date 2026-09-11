import * as THREE from "three";
import { sound } from "@/services/sound";

interface InteractionCallbacks {
  onAgentSelected: (agentId: string) => void;
  onLocationSelected: (locationId: string) => void;
}

export class InteractionSystem {
  private readonly raycaster = new THREE.Raycaster();
  private readonly pointer = new THREE.Vector2();

  constructor(
    private readonly renderer: THREE.WebGLRenderer,
    private readonly camera: THREE.PerspectiveCamera,
    private readonly getInteractiveObjects: () => THREE.Object3D[],
    private readonly callbacks: InteractionCallbacks
  ) {
    this.renderer.domElement.addEventListener("click", this.handleClick);
  }

  private handleClick = (event: MouseEvent) => {
    const rect = this.renderer.domElement.getBoundingClientRect();
    this.pointer.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
    this.pointer.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

    this.raycaster.setFromCamera(this.pointer, this.camera);
    const [hit] = this.raycaster.intersectObjects(this.getInteractiveObjects(), true);
    if (!hit) return;

    let object: THREE.Object3D | null = hit.object;
    while (object && !object.userData.type && object.parent) object = object.parent;
    if (!object?.userData?.type) return;

    sound.playClick();
    if (object.userData.type === "agent" && object.userData.agentId) {
      this.callbacks.onAgentSelected(object.userData.agentId);
    } else if (object.userData.type === "location" && object.userData.locationId) {
      this.callbacks.onLocationSelected(object.userData.locationId);
    }
  };

  dispose() {
    this.renderer.domElement.removeEventListener("click", this.handleClick);
  }
}
