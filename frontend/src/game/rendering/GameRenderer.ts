import * as THREE from "three";
import { CITY_SCENE_CONFIG } from "../config/scene";

export interface RenderingRuntime {
  scene: THREE.Scene;
  camera: THREE.PerspectiveCamera;
  renderer: THREE.WebGLRenderer;
  resize: () => void;
  dispose: () => void;
}

export function createRenderingRuntime(container: HTMLDivElement): RenderingRuntime {
  const scene = new THREE.Scene();
  scene.background = new THREE.Color(CITY_SCENE_CONFIG.background);
  scene.fog = new THREE.FogExp2(CITY_SCENE_CONFIG.fogColor, CITY_SCENE_CONFIG.fogDensity);

  const camera = new THREE.PerspectiveCamera(
    CITY_SCENE_CONFIG.camera.fov,
    Math.max(container.clientWidth, 1) / Math.max(container.clientHeight, 1),
    CITY_SCENE_CONFIG.camera.near,
    CITY_SCENE_CONFIG.camera.far
  );
  camera.position.set(...CITY_SCENE_CONFIG.camera.position);
  camera.lookAt(...CITY_SCENE_CONFIG.camera.target);

  const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, CITY_SCENE_CONFIG.maxDpr));

  const resize = () => {
    const width = Math.max(container.clientWidth, 1);
    const height = Math.max(container.clientHeight, 1);
    camera.aspect = width / height;
    camera.updateProjectionMatrix();
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, CITY_SCENE_CONFIG.maxDpr));
    renderer.setSize(width, height, false);
  };

  resize();
  container.appendChild(renderer.domElement);

  return {
    scene,
    camera,
    renderer,
    resize,
    dispose: () => {
      renderer.dispose();
      renderer.forceContextLoss();
      if (renderer.domElement.parentElement === container) {
        container.removeChild(renderer.domElement);
      }
    },
  };
}
