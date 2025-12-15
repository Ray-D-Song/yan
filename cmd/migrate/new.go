package migrate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const migrateDir = "internal/embedfs/sql/migrate"

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new migration file",
	Long:  `Create a new migration file interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Ask for migration name
		fmt.Print("Enter migration name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading migration name: %v\n", err)
			os.Exit(1)
		}
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Println("Migration name cannot be empty")
			os.Exit(1)
		}

		// Ask for migration description
		fmt.Print("Enter migration description: ")
		description, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading migration description: %v\n", err)
			os.Exit(1)
		}
		description = strings.TrimSpace(description)

		if err := createMigration(name, description); err != nil {
			fmt.Printf("Error creating migration: %v\n", err)
			os.Exit(1)
		}
	},
}

func createMigration(name, description string) error {
	// Ensure migration directory exists
	if err := os.MkdirAll(migrateDir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// Get next migration number
	nextNum, err := getNextMigrationNumber()
	if err != nil {
		return fmt.Errorf("failed to get next migration number: %w", err)
	}

	// Format migration number with leading zeros (001, 002, etc.)
	migrationNum := fmt.Sprintf("%03d", nextNum)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Default description if empty
	if description == "" {
		description = "Add your migration description here"
	}

	// Create up migration file
	upFile := filepath.Join(migrateDir, fmt.Sprintf("%s_%s.up.sql", migrationNum, name))
	upContent := fmt.Sprintf(`-- Migration: %s
-- Created at: %s
-- Description: %s

-- Write your UP migration here

`, name, currentTime, description)

	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downFile := filepath.Join(migrateDir, fmt.Sprintf("%s_%s.down.sql", migrationNum, name))
	downContent := fmt.Sprintf(`-- Migration: %s
-- Created at: %s
-- Description: %s

-- Write your DOWN migration here (rollback)

`, name, currentTime, description)

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("\n✓ Created migration files:\n")
	fmt.Printf("  - %s\n", upFile)
	fmt.Printf("  - %s\n", downFile)
	fmt.Printf("\nMigration number: %s\n", migrationNum)
	fmt.Printf("Description: %s\n", description)

	return nil
}

func getNextMigrationNumber() (int, error) {
	files, err := os.ReadDir(migrateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}

	maxNum := 0
	// Regular expression to match migration files: 001_name.up.sql or 001_name.down.sql
	re := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(file.Name())
		if len(matches) > 1 {
			num, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			if num > maxNum {
				maxNum = num
			}
		}
	}

	return maxNum + 1, nil
}
