package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/scenario"
	"city-of-lies/backend/internal/game/simulation"
)

type RunResult struct {
	Seed               int64                `json:"seed"`
	Status             domain.SessionStatus `json:"status"`
	EndGameSecond      int                  `json:"end_game_second"`
	FinalRatio         float64              `json:"final_ratio"`
	TotalConversations int                  `json:"total_conversations"`
	DurationMs         int64                `json:"duration_ms"`
}

type SoakSummary struct {
	TotalRuns              int           `json:"total_runs"`
	DefeatCount            int           `json:"defeat_count"`
	DefeatRatePercent      float64       `json:"defeat_rate_percent"`
	WonCount               int           `json:"won_count"`
	TimeoutCount           int           `json:"timeout_count"`
	MinDefeatSecond        int           `json:"min_defeat_second"`
	MaxDefeatSecond        int           `json:"max_defeat_second"`
	MeanDefeatSecond       float64       `json:"mean_defeat_second"`
	MedianDefeatSecond     float64       `json:"median_defeat_second"`
	MeanConversations      float64       `json:"mean_conversations"`
	TotalExecutionDuration time.Duration `json:"total_duration"`
	Verdict                string        `json:"verdict"`
}

func main() {
	runsFlag := flag.Int("runs", 100, "Number of simulation seeds to run")
	workersFlag := flag.Int("workers", 10, "Number of parallel worker goroutines")
	maxTicksFlag := flag.Int("max-ticks", 2400, "Maximum simulation ticks per run (2400 = 40 min)")
	scenarioFlag := flag.String("scenario", "", "Path to scenario bundle directory")
	jsonFlag := flag.Bool("json", false, "Output results in JSON format")
	flag.Parse()

	scenarioPath := *scenarioFlag
	if scenarioPath == "" {
		scenarioPath = filepath.Join("data", "scenarios", "riverside-factory")
		if _, err := os.Stat(scenarioPath); os.IsNotExist(err) {
			scenarioPath = filepath.Join("..", "data", "scenarios", "riverside-factory")
		}
	}

	bundle, err := scenario.LoadScenario(scenarioPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading scenario bundle: %v\n", err)
		os.Exit(1)
	}

	cfg, _ := config.Load()
	cfg.Simulation.MaxNewConversationsPerTick = 2
	cfg.Simulation.PairCooldownSeconds = 30

	if !*jsonFlag {
		fmt.Println("================================================================================")
		fmt.Println("  CITY OF LIES — AUTOMATED SIMULATION SOAK TEST (PHASE 11)")
		fmt.Println("================================================================================")
		fmt.Printf("Scenario:       %s (%s)\n", bundle.Metadata.Title, bundle.Metadata.ID)
		fmt.Printf("Total Seeds:    %d\n", *runsFlag)
		fmt.Printf("Concurrency:    %d workers\n", *workersFlag)
		fmt.Printf("Max Ticks/Run:  %d seconds (~%d mins)\n", *maxTicksFlag, *maxTicksFlag/60)
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("Progress: ")
	}

	startTime := time.Now()
	jobs := make(chan int64, *runsFlag)
	resultsChan := make(chan RunResult, *runsFlag)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < *workersFlag; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seed := range jobs {
				runStart := time.Now()
				engine := simulation.NewSimulationEngine(bundle, seed, cfg)

				totalConvs := 0
				engine.OnEventEmitted = func(evt domain.WorldEvent) {
					if evt.EventType == domain.WorldEventTypeRumorVisualized {
						totalConvs++
					}
				}

				for tick := 1; tick <= *maxTicksFlag; tick++ {
					engine.Tick(1)
					if !engine.Session.IsRunning() {
						break
					}
				}

				resultsChan <- RunResult{
					Seed:               seed,
					Status:             engine.Session.Status,
					EndGameSecond:      int(engine.Session.GameSecond),
					FinalRatio:         engine.Session.FalseNarrativeRatio,
					TotalConversations: totalConvs,
					DurationMs:         time.Since(runStart).Milliseconds(),
				}
			}
		}()
	}

	// Feed seeds
	for i := 1; i <= *runsFlag; i++ {
		jobs <- int64(i)
	}
	close(jobs)

	// Collect results
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []RunResult
	completed := 0
	for res := range resultsChan {
		results = append(results, res)
		completed++
		if !*jsonFlag && (*runsFlag <= 20 || completed%(*runsFlag/10) == 0 || completed == *runsFlag) {
			pct := float64(completed) / float64(*runsFlag) * 100
			fmt.Printf("%.0f%%.. ", pct)
		}
	}
	if !*jsonFlag {
		fmt.Println("Done!")
		fmt.Println()
	}

	totalDuration := time.Since(startTime)

	// Calculate statistics
	defeats := 0
	wons := 0
	timeouts := 0
	var defeatSeconds []int
	totalConvs := 0

	for _, r := range results {
		totalConvs += r.TotalConversations
		if r.Status == domain.SessionStatusLostFalseBelief {
			defeats++
			defeatSeconds = append(defeatSeconds, r.EndGameSecond)
		} else if r.Status == domain.SessionStatusWon {
			wons++
		} else {
			timeouts++
		}
	}

	sort.Ints(defeatSeconds)

	minSec, maxSec := 0, 0
	meanSec, medianSec := 0.0, 0.0
	if len(defeatSeconds) > 0 {
		minSec = defeatSeconds[0]
		maxSec = defeatSeconds[len(defeatSeconds)-1]
		sum := 0
		for _, s := range defeatSeconds {
			sum += s
		}
		meanSec = float64(sum) / float64(len(defeatSeconds))
		mid := len(defeatSeconds) / 2
		if len(defeatSeconds)%2 == 0 {
			medianSec = float64(defeatSeconds[mid-1]+defeatSeconds[mid]) / 2.0
		} else {
			medianSec = float64(defeatSeconds[mid])
		}
	}

	defeatRate := float64(defeats) / float64(*runsFlag) * 100
	meanConvs := float64(totalConvs) / float64(*runsFlag)

	verdict := "PASS"
	if defeatRate < 90.0 {
		verdict = "FAIL (Defeat rate < 90%)"
	}

	summary := SoakSummary{
		TotalRuns:              *runsFlag,
		DefeatCount:            defeats,
		DefeatRatePercent:      math.Round(defeatRate*100) / 100,
		WonCount:               wons,
		TimeoutCount:           timeouts,
		MinDefeatSecond:        minSec,
		MaxDefeatSecond:        maxSec,
		MeanDefeatSecond:       math.Round(meanSec*10) / 10,
		MedianDefeatSecond:     math.Round(medianSec*10) / 10,
		MeanConversations:      math.Round(meanConvs*10) / 10,
		TotalExecutionDuration: totalDuration,
		Verdict:                verdict,
	}

	if *jsonFlag {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(summary)
		return
	}

	fmt.Println("================================================================================")
	fmt.Println("  SOAK TEST METRICS & CALIBRATION RESULTS")
	fmt.Println("================================================================================")
	fmt.Printf("Total Runs Simulated:     %d seeds\n", summary.TotalRuns)
	fmt.Printf("Social Defeats:           %d / %d (%.1f%%)  [Target: > 90.0%%]\n", summary.DefeatCount, summary.TotalRuns, summary.DefeatRatePercent)
	fmt.Printf("Survivals / Timeouts:     %d\n", summary.TimeoutCount)
	fmt.Printf("Player Victories (unplayed): %d\n", summary.WonCount)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Min Time to Defeat:       %d s (%02d:%02d)\n", summary.MinDefeatSecond, summary.MinDefeatSecond/60, summary.MinDefeatSecond%60)
	fmt.Printf("Max Time to Defeat:       %d s (%02d:%02d)\n", summary.MaxDefeatSecond, summary.MaxDefeatSecond/60, summary.MaxDefeatSecond%60)
	fmt.Printf("Mean Time to Defeat:      %.1f s (~%02d:%02d)\n", summary.MeanDefeatSecond, int(summary.MeanDefeatSecond)/60, int(summary.MeanDefeatSecond)%60)
	fmt.Printf("Median Time to Defeat:    %.1f s (~%02d:%02d)\n", summary.MedianDefeatSecond, int(summary.MedianDefeatSecond)/60, int(summary.MedianDefeatSecond)%60)
	fmt.Printf("Mean NPC Conversations:   %.1f convs/game\n", summary.MeanConversations)
	fmt.Printf("Execution Clock Time:     %v (%.2f ms/game)\n", summary.TotalExecutionDuration, float64(summary.TotalExecutionDuration.Milliseconds())/float64(summary.TotalRuns))
	fmt.Println("================================================================================")
	if verdict == "PASS" {
		fmt.Printf("FINAL VERDICT: [PASS] - Simulation is mathematically stable and correctly balanced!\n")
	} else {
		fmt.Printf("FINAL VERDICT: [%s]\n", verdict)
		os.Exit(1)
	}
	fmt.Println("================================================================================")
}
