package main

import (
	"fmt"
	"os"
	"path/filepath"

	"city-of-lies/backend/internal/game/scenario"
)

func main() {
	targetDir := "data/scenarios/riverside-factory"
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	cleanPath, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Printf("Error resolving scenario path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking scenario bundle at: %s\n", cleanPath)

	bundle, err := scenario.LoadScenario(cleanPath)
	if err != nil {
		fmt.Printf("FAIL: Failed to load scenario: %v\n", err)
		os.Exit(1)
	}

	if err := scenario.ValidateScenarioBundle(bundle); err != nil {
		fmt.Printf("FAIL: Scenario validation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("PASS: Scenario bundle is valid!")
	fmt.Printf("  - Scenario: %s (v%d)\n", bundle.Metadata.Title, bundle.Metadata.Version)
	fmt.Printf("  - Locations: %d\n", len(bundle.Locations))
	fmt.Printf("  - Agents: %d\n", len(bundle.Agents))
	fmt.Printf("  - Claims: %d\n", len(bundle.Claims))
	fmt.Printf("  - Initial Observations: %d\n", len(bundle.Observations))
	fmt.Printf("  - Relationships: %d\n", len(bundle.Relationships))
	fmt.Printf("  - Evidence Items: %d\n", len(bundle.Evidence))
	fmt.Printf("  - Public Events: %d\n", len(bundle.PublicEvents))
}
