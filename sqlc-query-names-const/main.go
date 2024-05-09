package main

import (
	"bufio"

	"github.com/spf13/cobra"

	"fmt"
	"os"
	"regexp"
	"strings"
)

// Reads the content of the SQLC query file.
func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Parses the query names from the file content.
func parseQueryNames(content string) ([]string, error) {
	var queryNames []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	queryNameRegex := regexp.MustCompile(`--\s*name:\s*(\w+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := queryNameRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			queryNames = append(queryNames, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return queryNames, nil
}

// Generates Go code with constants for the parsed query names.
func generateGoCode(queryNames []string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("package %s\n\n", packageName))
	builder.WriteString("// Query names as constants\nconst (\n")
	constantPrefix = strings.ToUpper(constantPrefix)
	for _, name := range queryNames {
		builder.WriteString(fmt.Sprintf("    %s%s = \"%s\"\n", constantPrefix, name, name))
	}
	builder.WriteString(")\n")
	return builder.String()
}

// Writes the generated Go code to a file.
func writeGoCodeToFile(code, fileName string) error {
	return os.WriteFile(fileName, []byte(code), 0644)
}

var (
	inputFile      string
	outputFile     string
	packageName    string
	constantPrefix string
)
var rootCmd = &cobra.Command{
	Use:   "sqlc-query-names-const",
	Short: "Generate Go constants from sqlc query names",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := readFileContent(inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		queryNames, err := parseQueryNames(content)
		if err != nil {
			fmt.Printf("Error parsing query names: %v\n", err)
			os.Exit(1)
		}

		goCode := generateGoCode(queryNames)
		if err := writeGoCodeToFile(goCode, outputFile); err != nil {
			fmt.Printf("Error writing Go code to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Constants generated successfully in %s\n", outputFile)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to the sqlc query file e.g. db/queries.sql")
	rootCmd.MarkFlagRequired("input")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file name e.g. querynames.go")
	rootCmd.MarkFlagRequired("input")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "main", "package name for the generated Go file")
	rootCmd.Flags().StringVarP(&constantPrefix, "prefix", "c", "SQL", "prefix for the generated constants")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
