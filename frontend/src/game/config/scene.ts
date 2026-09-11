export const CITY_SCENE_CONFIG = {
  background: 0x0a0c10,
  fogColor: 0x0a0c10,
  fogDensity: 0.012,
  groundSize: 140,
  gridDivisions: 70,
  maxDpr: 2,
  camera: {
    fov: 40,
    near: 0.1,
    far: 500,
    position: [0, 65, 75] as [number, number, number],
    target: [0, 0, 0] as [number, number, number],
    minDistance: 35,
    maxDistance: 145,
  },
} as const;
