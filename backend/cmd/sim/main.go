package main

import (
	"fmt"
	"os"
	"path/filepath"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/scenario"
	"city-of-lies/backend/internal/game/simulation"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("City of Lies - Headless Simulation Engine")
	fmt.Println("==================================================")

	scenarioPath := filepath.Join("data", "scenarios", "riverside-factory")
	if _, err := os.Stat(scenarioPath); os.IsNotExist(err) {
		scenarioPath = filepath.Join("..", "data", "scenarios", "riverside-factory")
	}

	bundle, err := scenario.LoadScenario(scenarioPath)
	if err != nil {
		fmt.Printf("Error loading scenario: %v\n", err)
		os.Exit(1)
	}

	cfg, _ := config.Load()
	cfg.Simulation.MaxNewConversationsPerTick = 2
	cfg.Simulation.PairCooldownSeconds = 30

	seed := int64(424242)
	engine := simulation.NewSimulationEngine(bundle, seed, cfg)

	fmt.Printf("Loaded Scenario: %s\n", bundle.Metadata.Title)
	fmt.Printf("Seed: %d | Starting False Belief: %.2f%%\n\n", seed, engine.Session.FalseNarrativeRatio*100)

	totalConversations := 0
	engine.OnEventEmitted = func(evt domain.WorldEvent) {
		if evt.EventType == domain.WorldEventTypeRumorVisualized {
			totalConversations++
		}
	}

	// Run up to 2200 ticks (36 game minutes)
	for tick := 1; tick <= 2200; tick++ {
		engine.Tick(1)

		if tick%30 == 0 || !engine.Session.IsRunning() {
			fmt.Printf("[Second %03d] False Narrative: %5.1f%% | Conversations: %d | Status: %s\n",
				engine.Session.GameSecond,
				engine.Session.FalseNarrativeRatio*100,
				totalConversations,
				engine.Session.Status,
			)
		}

		if !engine.Session.IsRunning() {
			break
		}
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Printf("Simulation Concluded: %s\n", engine.Session.Status)
	fmt.Printf("Final False Narrative Ratio: %.2f%%\n", engine.Session.FalseNarrativeRatio*100)
	fmt.Printf("Total NPC Conversations: %d\n", totalConversations)
	fmt.Println("--------------------------------------------------")
}
