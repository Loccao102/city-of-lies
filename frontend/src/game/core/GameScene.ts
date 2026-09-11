import * as THREE from "three";
import type { AgentSummary } from "@/types/game";
import { AgentSystem } from "../entities/npc/AgentSystem";
import { createRenderingRuntime } from "../rendering/GameRenderer";
import { installSceneLighting } from "../rendering/SceneLighting";
import { MapCameraSystem } from "../systems/camera/MapCameraSystem";
import { InteractionSystem } from "../systems/interaction/InteractionSystem";
import { createCityWorld } from "../world/CityWorld";

export interface GameSceneCallbacks {
  onAgentSelected: (agentId: string) => void;
  onLocationSelected: (locationId: string) => void;
}

/**
 * Composition root for the browser 3D runtime.
 * It wires renderer, world, entities and systems but owns no domain rules.
 */
export class GameScene {
  private readonly rendering;
  private readonly lighting;
  private readonly world;
  private readonly agents;
  private readonly cameraSystem;
  private readonly interactionSystem;
  private readonly clock = new THREE.Clock();
  private animationFrame: number | null = null;

  constructor(container: HTMLDivElement, callbacks: GameSceneCallbacks) {
    this.rendering = createRenderingRuntime(container);
    this.lighting = installSceneLighting(this.rendering.scene);
    this.world = createCityWorld(this.rendering.scene);
    this.agents = new AgentSystem(this.rendering.scene, this.world.locationMap);
    this.cameraSystem = new MapCameraSystem(
      container,
      this.rendering.camera,
      this.rendering.renderer,
      this.rendering.resize
    );
    this.interactionSystem = new InteractionSystem(
      this.rendering.renderer,
      this.rendering.camera,
      () => [...this.world.interactiveObjects, ...this.agents.getInteractiveObjects()],
      callbacks
    );
    this.animate();
  }

  setAgents(agents: AgentSummary[]) {
    this.agents.setAgents(agents);
  }

  private animate = () => {
    this.animationFrame = requestAnimationFrame(this.animate);
    const elapsed = this.clock.getElapsedTime();
    this.lighting.update(elapsed);
    this.agents.update(elapsed);
    this.rendering.renderer.render(this.rendering.scene, this.rendering.camera);
  };

  dispose() {
    if (this.animationFrame !== null) cancelAnimationFrame(this.animationFrame);
    this.interactionSystem.dispose();
    this.cameraSystem.dispose();
    this.agents.dispose();
    this.world.dispose();
    this.lighting.dispose();
    this.rendering.dispose();
  }
}
