"use client";

import React, { useEffect, useRef } from "react";
import * as THREE from "three";
import { useGameStore } from "@/store/gameStore";
import { useUIStore } from "@/store/uiStore";

// The 12 Riverside district locations with coordinates
export const DISTRICT_LOCATIONS = [
  { id: "factory_gate", name: "Cổng Nhà Máy", x: -24, z: -16, color: 0x4a5568, category: "factory" },
  { id: "storage_building", name: "Nhà Kho DB-4", x: -32, z: -24, color: 0xe53e3e, category: "factory" },
  { id: "security_room", name: "Phòng An Ninh", x: -28, z: -8, color: 0x2b6cb0, category: "factory" },
  { id: "factory_yard", name: "Sân Nhà Máy", x: -20, z: -26, color: 0x718096, category: "factory" },
  { id: "factory_office", name: "Văn Phòng Quản Lý", x: -14, z: -18, color: 0xd69e2e, category: "factory" },
  { id: "hospital", name: "Bệnh Viện Đa Khoa", x: 26, z: -20, color: 0x319795, category: "official" },
  { id: "fire_station", name: "Trạm Cứu Hỏa", x: -6, z: -6, color: 0xc53030, category: "official" },
  { id: "news_office", name: "Tòa Soạn & Blogger", x: 18, z: 12, color: 0x805ad5, category: "public" },
  { id: "public_square", name: "Quảng Trường", x: 2, z: 6, color: 0xdd6b20, category: "public" },
  { id: "market", name: "Chợ Dân Sinh", x: -10, z: 18, color: 0x38a169, category: "public" },
  { id: "cafe", name: "Quán Cà Phê Vỉa Hè", x: 10, z: -4, color: 0xd69e2e, category: "public" },
  { id: "residential_street", name: "Khu Dân Cư", x: -22, z: 14, color: 0x4a5568, category: "residential" },
];

export function DistrictMap() {
  const mountRef = useRef<HTMLDivElement>(null);
  const agents = useGameStore((s) => s.agents);
  const recentPulses = useGameStore((s) => s.recentPulses);
  const { openInterviewModal, openInspectModal } = useUIStore();

  useEffect(() => {
    if (!mountRef.current) return;
    const container = mountRef.current;
    const width = container.clientWidth;
    const height = container.clientHeight;

    // 1. Scene setup
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x0a0c10);
    scene.fog = new THREE.FogExp2(0x0a0c10, 0.012);

    // 2. Camera (Isometric Orthographic / Angled Perspective)
    const camera = new THREE.PerspectiveCamera(40, width / height, 0.1, 500);
    camera.position.set(0, 65, 75);
    camera.lookAt(0, 0, 0);

    // 3. Renderer
    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
    renderer.setSize(width, height);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;
    container.appendChild(renderer.domElement);

    // 4. Lighting
    const ambientLight = new THREE.AmbientLight(0xdce7f5, 0.6);
    scene.add(ambientLight);

    const dirLight = new THREE.DirectionalLight(0xfff5e6, 1.2);
    dirLight.position.set(30, 60, 40);
    dirLight.castShadow = true;
    dirLight.shadow.mapSize.width = 2048;
    dirLight.shadow.mapSize.height = 2048;
    dirLight.shadow.camera.near = 10;
    dirLight.shadow.camera.far = 150;
    dirLight.shadow.camera.left = -50;
    dirLight.shadow.camera.right = 50;
    dirLight.shadow.camera.top = 50;
    dirLight.shadow.camera.bottom = -50;
    scene.add(dirLight);

    const streetLight = new THREE.PointLight(0x4fd1c5, 1.5, 40);
    streetLight.position.set(0, 10, 5);
    scene.add(streetLight);

    // 5. Ground Grid & Roads
    const groundGeo = new THREE.PlaneGeometry(140, 140);
    const groundMat = new THREE.MeshStandardMaterial({
      color: 0x0e131b,
      roughness: 0.9,
      metalness: 0.1,
    });
    const ground = new THREE.Mesh(groundGeo, groundMat);
    ground.rotation.x = -Math.PI / 2;
    ground.receiveShadow = true;
    scene.add(ground);

    const grid = new THREE.GridHelper(140, 70, 0x1f293d, 0x141b27);
    grid.position.y = 0.02;
    scene.add(grid);

    // 6. Spawn 12 Location Buildings & Labels
    const interactiveObjects: THREE.Object3D[] = [];
    const locationMap = new Map<string, THREE.Vector3>();

    DISTRICT_LOCATIONS.forEach((loc) => {
      const pos = new THREE.Vector3(loc.x, 0, loc.z);
      locationMap.set(loc.id, pos);

      // Building Box
      const bHeight = loc.category === "factory" ? 7 : loc.category === "official" ? 9 : 5;
      const bGeo = new THREE.BoxGeometry(9, bHeight, 7);
      const bMat = new THREE.MeshStandardMaterial({
        color: loc.color,
        roughness: 0.7,
        metalness: 0.2,
      });
      const building = new THREE.Mesh(bGeo, bMat);
      building.position.set(loc.x, bHeight / 2, loc.z);
      building.castShadow = true;
      building.receiveShadow = true;
      building.userData = { type: "location", locationId: loc.id, name: loc.name };
      scene.add(building);
      interactiveObjects.push(building);

      // Roof Marker / Glow Beacon
      const beaconGeo = new THREE.CylinderGeometry(0.3, 0.3, 1.2, 8);
      const beaconMat = new THREE.MeshBasicMaterial({ color: loc.color });
      const beacon = new THREE.Mesh(beaconGeo, beaconMat);
      beacon.position.set(loc.x, bHeight + 0.6, loc.z);
      scene.add(beacon);

      // Base footprint ring
      const ringGeo = new THREE.RingGeometry(5.5, 6.2, 32);
      const ringMat = new THREE.MeshBasicMaterial({
        color: loc.color,
        side: THREE.DoubleSide,
        transparent: true,
        opacity: 0.3,
      });
      const ring = new THREE.Mesh(ringGeo, ringMat);
      ring.rotation.x = -Math.PI / 2;
      ring.position.set(loc.x, 0.05, loc.z);
      scene.add(ring);
    });

    // 7. Spawn Agent Pins
    const agentMeshes = new Map<string, THREE.Group>();

    agents.forEach((ag, index) => {
      const locPos = locationMap.get(ag.location_id) || new THREE.Vector3(0, 0, 0);

      const agentGroup = new THREE.Group();
      // Offset agents around their location building so they don't overlap
      const angle = (index * (Math.PI * 2)) / 5;
      const offsetRadius = 6.0;
      agentGroup.position.set(
        locPos.x + Math.cos(angle) * offsetRadius,
        0,
        locPos.z + Math.sin(angle) * offsetRadius
      );

      // Agent Cylinder Body
      const bodyGeo = new THREE.CylinderGeometry(0.5, 0.5, 2.2, 12);
      const bodyMat = new THREE.MeshStandardMaterial({
        color: ag.interviewed ? 0x4fd1c5 : 0xf6ad55,
        roughness: 0.5,
      });
      const body = new THREE.Mesh(bodyGeo, bodyMat);
      body.position.y = 1.1;
      body.castShadow = true;
      agentGroup.add(body);

      // Agent Head Sphere
      const headGeo = new THREE.SphereGeometry(0.55, 12, 12);
      const headMat = new THREE.MeshStandardMaterial({ color: 0xffffff });
      const head = new THREE.Mesh(headGeo, headMat);
      head.position.y = 2.6;
      head.castShadow = true;
      agentGroup.add(head);

      agentGroup.userData = { type: "agent", agentId: ag.id, name: ag.name, role: ag.role };
      scene.add(agentGroup);
      interactiveObjects.push(body);
      agentMeshes.set(ag.id, agentGroup);
    });

    // 8. Raycasting for Clicks & Hover
    const raycaster = new THREE.Raycaster();
    const mouse = new THREE.Vector2();

    const handleClick = (event: MouseEvent) => {
      const rect = renderer.domElement.getBoundingClientRect();
      mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
      mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;

      raycaster.setFromCamera(mouse, camera);
      const intersects = raycaster.intersectObjects(interactiveObjects, true);

      if (intersects.length > 0) {
        let obj: THREE.Object3D | null = intersects[0].object;
        while (obj && !obj.userData.type && obj.parent) {
          obj = obj.parent;
        }

        if (obj && obj.userData) {
          if (obj.userData.type === "agent") {
            openInterviewModal(obj.userData.agentId);
          } else if (obj.userData.type === "location") {
            openInspectModal(obj.userData.locationId);
          }
        }
      }
    };

    renderer.domElement.addEventListener("click", handleClick);

    // 9. Camera Drag Controls
    let isDragging = false;
    let prevMouse = { x: 0, y: 0 };

    const handleMouseDown = (e: MouseEvent) => {
      if (e.button === 0 || e.button === 2) {
        isDragging = true;
        prevMouse = { x: e.clientX, y: e.clientY };
      }
    };
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDragging) return;
      const dx = e.clientX - prevMouse.x;
      const dy = e.clientY - prevMouse.y;
      camera.position.x -= dx * 0.15;
      camera.position.z -= dy * 0.15;
      prevMouse = { x: e.clientX, y: e.clientY };
    };
    const handleMouseUp = () => {
      isDragging = false;
    };
    const handleWheel = (e: WheelEvent) => {
      camera.position.y = Math.max(25, Math.min(120, camera.position.y + e.deltaY * 0.05));
    };

    renderer.domElement.addEventListener("mousedown", handleMouseDown);
    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);
    renderer.domElement.addEventListener("wheel", handleWheel);

    // 10. Animation Loop & Visual Pulses
    let animId: number;
    let clock = new THREE.Clock();

    const animate = () => {
      animId = requestAnimationFrame(animate);
      const elapsed = clock.getElapsedTime();

      // Subtle float animation for agent heads
      agentMeshes.forEach((mesh, _) => {
        mesh.position.y = Math.sin(elapsed * 2 + mesh.position.x) * 0.15;
      });

      renderer.render(scene, camera);
    };
    animate();

    const handleResize = () => {
      if (!container) return;
      const w = container.clientWidth;
      const h = container.clientHeight;
      camera.aspect = w / h;
      camera.updateProjectionMatrix();
      renderer.setSize(w, h);
    };
    window.addEventListener("resize", handleResize);

    return () => {
      cancelAnimationFrame(animId);
      renderer.domElement.removeEventListener("click", handleClick);
      renderer.domElement.removeEventListener("mousedown", handleMouseDown);
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
      renderer.domElement.removeEventListener("wheel", handleWheel);
      window.removeEventListener("resize", handleResize);
      if (container && renderer.domElement) {
        container.removeChild(renderer.domElement);
      }
      renderer.dispose();
    };
  }, [agents, openInterviewModal, openInspectModal]);

  return (
    <div className="relative w-full h-full">
      <div ref={mountRef} className="w-full h-full cursor-grab active:cursor-grabbing" />

      {/* Map Controls Helper */}
      <div className="absolute top-20 left-6 z-10 glass-panel px-3.5 py-2 rounded-lg text-[11px] text-gray-400 pointer-events-none flex items-center gap-3 border border-white/10">
        <span>Kéo chuột để di chuyển camera</span>
        <span>•</span>
        <span>Lăn chuột để phóng to/thu nhỏ</span>
        <span>•</span>
        <span className="text-truth-glow">Bấm vào người hoặc công trình để tương tác</span>
      </div>

      {/* Active Rumor Broadcast Banner */}
      {recentPulses.length > 0 && (
        <div className="absolute top-20 right-6 z-10 glass-panel-elevated px-4 py-2.5 rounded-xl border border-amber-500/30 flex items-center gap-3 animate-fade-in shadow-xl">
          <div className="w-2.5 h-2.5 rounded-full bg-amber-400 animate-ping" />
          <div className="text-xs">
            <span className="font-semibold text-amber-300">
              {recentPulses[recentPulses.length - 1].speakerName}
            </span>
            <span className="text-gray-400"> vừa kể cho </span>
            <span className="font-semibold text-gray-200">
              {recentPulses[recentPulses.length - 1].listenerName}
            </span>
          </div>
        </div>
      )}
    </div>
  );
}
