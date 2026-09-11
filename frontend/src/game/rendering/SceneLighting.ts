import * as THREE from "three";

export interface SceneLightingRuntime {
  update: (elapsed: number) => void;
  dispose: () => void;
}

export function installSceneLighting(scene: THREE.Scene): SceneLightingRuntime {
  const root = new THREE.Group();
  root.name = "scene-lighting";

  const ambientLight = new THREE.AmbientLight(0xdce7f5, 0.6);
  root.add(ambientLight);

  const dirLight = new THREE.DirectionalLight(0xfff5e6, 1.2);
  dirLight.position.set(30, 60, 40);
  dirLight.castShadow = true;
  dirLight.shadow.mapSize.set(2048, 2048);
  dirLight.shadow.camera.near = 10;
  dirLight.shadow.camera.far = 150;
  dirLight.shadow.camera.left = -50;
  dirLight.shadow.camera.right = 50;
  dirLight.shadow.camera.top = 50;
  dirLight.shadow.camera.bottom = -50;
  root.add(dirLight);

  const streetLight = new THREE.PointLight(0x4fd1c5, 1.5, 40);
  streetLight.position.set(0, 10, 5);
  root.add(streetLight);

  const factoryStrobeRed = new THREE.PointLight(0xff2222, 2.5, 35);
  factoryStrobeRed.position.set(-24, 8, -16);
  root.add(factoryStrobeRed);

  const factoryStrobeBlue = new THREE.PointLight(0x2266ff, 2.5, 35);
  factoryStrobeBlue.position.set(-24, 8, -16);
  root.add(factoryStrobeBlue);

  const clinicStrobeRed = new THREE.PointLight(0xff2222, 2.0, 30);
  clinicStrobeRed.position.set(26, 10, -20);
  root.add(clinicStrobeRed);

  const clinicStrobeBlue = new THREE.PointLight(0x2266ff, 2.0, 30);
  clinicStrobeBlue.position.set(26, 10, -20);
  root.add(clinicStrobeBlue);

  scene.add(root);

  return {
    update: (elapsed) => {
      const strobeOn = Math.floor(elapsed * 5) % 2 === 0;
      factoryStrobeRed.intensity = strobeOn ? 3.5 : 0;
      factoryStrobeBlue.intensity = strobeOn ? 0 : 3.5;
      clinicStrobeRed.intensity = strobeOn ? 0 : 2.5;
      clinicStrobeBlue.intensity = strobeOn ? 2.5 : 0;
    },
    dispose: () => {
      scene.remove(root);
    },
  };
}
