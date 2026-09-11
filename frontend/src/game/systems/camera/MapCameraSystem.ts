import * as THREE from "three";
import { CITY_SCENE_CONFIG } from "../../config/scene";

export class MapCameraSystem {
  private dragging = false;
  private previousPointer = { x: 0, y: 0 };
  private readonly target = new THREE.Vector3(...CITY_SCENE_CONFIG.camera.target);

  constructor(
    private readonly container: HTMLDivElement,
    private readonly camera: THREE.PerspectiveCamera,
    private readonly renderer: THREE.WebGLRenderer,
    private readonly resizeRenderer: () => void
  ) {
    const canvas = this.renderer.domElement;
    canvas.addEventListener("mousedown", this.handleMouseDown);
    canvas.addEventListener("wheel", this.handleWheel, { passive: false });
    canvas.addEventListener("contextmenu", this.handleContextMenu);
    window.addEventListener("mousemove", this.handleMouseMove);
    window.addEventListener("mouseup", this.handleMouseUp);
    window.addEventListener("resize", this.handleResize);
  }

  private handleMouseDown = (event: MouseEvent) => {
    if (event.button !== 0 && event.button !== 2) return;
    this.dragging = true;
    this.previousPointer = { x: event.clientX, y: event.clientY };
  };

  private handleMouseMove = (event: MouseEvent) => {
    if (!this.dragging) return;
    const dx = event.clientX - this.previousPointer.x;
    const dy = event.clientY - this.previousPointer.y;
    const panX = -dx * 0.15;
    const panZ = -dy * 0.15;

    this.camera.position.x += panX;
    this.camera.position.z += panZ;
    this.target.x += panX;
    this.target.z += panZ;
    this.camera.lookAt(this.target);
    this.previousPointer = { x: event.clientX, y: event.clientY };
  };

  private handleMouseUp = () => {
    this.dragging = false;
  };

  private handleWheel = (event: WheelEvent) => {
    event.preventDefault();
    const offset = this.camera.position.clone().sub(this.target);
    const currentDistance = offset.length();
    const nextDistance = THREE.MathUtils.clamp(
      currentDistance * (1 + event.deltaY * 0.001),
      CITY_SCENE_CONFIG.camera.minDistance,
      CITY_SCENE_CONFIG.camera.maxDistance
    );
    offset.setLength(nextDistance);
    this.camera.position.copy(this.target).add(offset);
    this.camera.lookAt(this.target);
  };

  private handleContextMenu = (event: MouseEvent) => event.preventDefault();

  private handleResize = () => {
    if (!this.container.isConnected) return;
    this.resizeRenderer();
  };

  dispose() {
    const canvas = this.renderer.domElement;
    canvas.removeEventListener("mousedown", this.handleMouseDown);
    canvas.removeEventListener("wheel", this.handleWheel);
    canvas.removeEventListener("contextmenu", this.handleContextMenu);
    window.removeEventListener("mousemove", this.handleMouseMove);
    window.removeEventListener("mouseup", this.handleMouseUp);
    window.removeEventListener("resize", this.handleResize);
  }
}
