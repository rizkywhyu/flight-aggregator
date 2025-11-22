package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("🧪 Running Flight Aggregator Unit Tests")
	fmt.Println("=====================================")

	// Test packages to run
	testPackages := []string{
		"./internal/models",
		"./internal/utils",
		"./internal/usecase",
		"./internal/service", 
		"./internal/controller",
		"./internal/providers",
	}

	var failedTests []string
	totalTests := 0
	passedTests := 0

	for _, pkg := range testPackages {
		fmt.Printf("\n📦 Testing package: %s\n", pkg)
		fmt.Println(strings.Repeat("-", 50))

		cmd := exec.Command("go", "test", "-v", pkg)
		output, err := cmd.CombinedOutput()
		
		fmt.Print(string(output))

		if err != nil {
			failedTests = append(failedTests, pkg)
			fmt.Printf("❌ Tests failed in %s\n", pkg)
		} else {
			fmt.Printf("✅ All tests passed in %s\n", pkg)
		}

		// Count tests (basic parsing)
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "RUN") {
				totalTests++
			}
			if strings.Contains(line, "PASS") && strings.Contains(line, "Test") {
				passedTests++
			}
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total Packages: %d\n", len(testPackages))
	fmt.Printf("Passed Packages: %d\n", len(testPackages)-len(failedTests))
	fmt.Printf("Failed Packages: %d\n", len(failedTests))
	fmt.Printf("Estimated Total Tests: %d\n", totalTests)
	fmt.Printf("Estimated Passed Tests: %d\n", passedTests)

	if len(failedTests) > 0 {
		fmt.Println("\n❌ Failed packages:")
		for _, pkg := range failedTests {
			fmt.Printf("  - %s\n", pkg)
		}
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All tests passed successfully!")
	}
}