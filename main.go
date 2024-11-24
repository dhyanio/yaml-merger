package main

// Simple but Powerful YAML Content Merger
import (
	"fmt"
	"os"
	"strings"

	"github.com/dhyanio/gogger"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var (
	ignoreList    = pflag.StringSlice("ignore", []string{}, "keys to ignore when merging")
	outputFile    = pflag.String("output", "", "output file to write the merged YAML (optional)")
	dryRun        = pflag.Bool("dry-run", false, "perform a dry run, showing what changes would be made without applying them")
	mergeStrategy = pflag.String("merge-strategy", "merge", "specify merge strategy: 'merge' (deep merge) or 'override' (override conflicting keys)")
)

// Init initializes the logger and parses command-line flags.
func init() {
	pflag.Parse()
}

func main() {
	log, err := gogger.NewLogger("logfile.log", gogger.INFO)
	if err != nil {
		panic(err)
	}
	if len(pflag.Args()) < 1 {
		log.Fatal().Msgf("Usage: %s <file1.yaml> [<file2.yaml> ...]", os.Args[0])
	}

	var result interface{}
	for _, file := range pflag.Args() {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read file: %s", file)
		}

		var parsedData interface{}
		if err := yaml.Unmarshal(content, &parsedData); err != nil {
			log.Fatal().Err(err).Msgf("Failed to parse YAML from file: %s", file)
		}

		result, err = Merge(result, parsedData, *mergeStrategy, log)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to merge YAML content from file: %s", file)
		}
	}

	mergedYAML, err := yaml.Marshal(result)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal merged result to YAML")
	}

	if *dryRun {
		// Dry run: Output the summary and show the changes without applying them
		fmt.Println("Dry Run: The following merged YAML would be generated (no file changes):")
		fmt.Println(string(mergedYAML))
	} else {
		// Perform actual output to file or stdout
		if *outputFile != "" {
			err := os.WriteFile(*outputFile, mergedYAML, 0644)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed to write merged YAML to file: %s", *outputFile)
			}
			fmt.Printf("Merged YAML written to %s\n", *outputFile)
		} else {
			// Using strings.Builder to efficiently concatenate the final output
			var output strings.Builder
			output.WriteString(string(mergedYAML))
			fmt.Println(output.String())
		}
	}
}

// Merge recursively merges or overrides two YAML structures based on the strategy.
func Merge(a, b interface{}, strategy string, log *gogger.Logger) (interface{}, error) {
	log.Debug().Msgf("Merging with strategy '%s': %v (%T) with %v (%T)", strategy, a, a, b, b)
	switch typedA := a.(type) {
	case []interface{}:
		typedB, ok := b.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected list on right side, got %T", b)
		}
		return append(typedA, typedB...), nil
	case map[interface{}]interface{}:
		typedB, ok := b.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map on right side, got %T", b)
		}
		for key, rightVal := range typedB {
			if shouldIgnore(key) {
				continue
			}
			leftVal, found := typedA[key]
			if !found || strategy == "override" {
				typedA[key] = rightVal
			} else {
				mergedVal, err := Merge(leftVal, rightVal, strategy, log)
				if err != nil {
					return nil, err
				}
				typedA[key] = mergedVal
			}
		}
		return typedA, nil
	default:
		return b, nil
	}
}

// shouldIgnore checks whether a given key is in the ignore list.
func shouldIgnore(key interface{}) bool {
	str, ok := key.(string)
	if !ok {
		return false
	}
	for _, ignoreKey := range *ignoreList {
		if str == ignoreKey {
			return true
		}
	}
	return false
}
