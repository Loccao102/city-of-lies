import * as THREE from "three";
import { CITY_SCENE_CONFIG } from "../config/scene";
import { DISTRICT_LOCATIONS } from "./data/districts";

export interface CityWorldRuntime {
  interactiveObjects: THREE.Object3D[];
  locationMap: Map<string, THREE.Vector3>;
  dispose: () => void;
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

export function createCityWorld(scene: THREE.Scene): CityWorldRuntime {
  const root = new THREE.Group();
  root.name = "city-world";
  const interactiveObjects: THREE.Object3D[] = [];
  const locationMap = new Map<string, THREE.Vector3>();

  const groundGeometry = new THREE.PlaneGeometry(CITY_SCENE_CONFIG.groundSize, CITY_SCENE_CONFIG.groundSize);
  const groundMaterial = new THREE.MeshStandardMaterial({
    color: 0x0e131b,
    roughness: 0.9,
    metalness: 0.1,
  });
  const ground = new THREE.Mesh(groundGeometry, groundMaterial);
  ground.rotation.x = -Math.PI / 2;
  ground.receiveShadow = true;
  root.add(ground);

  const grid = new THREE.GridHelper(
    CITY_SCENE_CONFIG.groundSize,
    CITY_SCENE_CONFIG.gridDivisions,
    0x1f293d,
    0x141b27
  );
  grid.position.y = 0.02;
  root.add(grid);

  DISTRICT_LOCATIONS.forEach((location) => {
    const position = new THREE.Vector3(location.x, 0, location.z);
    locationMap.set(location.id, position);

    const buildingHeight = location.category === "factory" ? 7 : location.category === "official" ? 9 : 5;
    const building = new THREE.Mesh(
      new THREE.BoxGeometry(9, buildingHeight, 7),
      new THREE.MeshStandardMaterial({
        color: location.color,
        roughness: 0.7,
        metalness: 0.2,
      })
    );
    building.position.set(location.x, buildingHeight / 2, location.z);
    building.castShadow = true;
    building.receiveShadow = true;
    building.userData = {
      type: "location",
      locationId: location.id,
      name: location.name,
    };
    root.add(building);
    interactiveObjects.push(building);

    const beacon = new THREE.Mesh(
      new THREE.CylinderGeometry(0.3, 0.3, 1.2, 8),
      new THREE.MeshBasicMaterial({ color: location.color })
    );
    beacon.position.set(location.x, buildingHeight + 0.6, location.z);
    root.add(beacon);

    const ring = new THREE.Mesh(
      new THREE.RingGeometry(5.5, 6.2, 32),
      new THREE.MeshBasicMaterial({
        color: location.color,
        side: THREE.DoubleSide,
        transparent: true,
        opacity: 0.3,
      })
    );
    ring.rotation.x = -Math.PI / 2;
    ring.position.set(location.x, 0.05, location.z);
    root.add(ring);
  });

  scene.add(root);

  return {
    interactiveObjects,
    locationMap,
    dispose: () => {
      scene.remove(root);
      disposeObjectTree(root);
    },
  };
}
