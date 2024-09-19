package main

// Simple but Powerful YAML Content Merger
import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var ignoreList = pflag.StringSlice("ignore", []string{}, "keys to ignore when merging")

// Init initializes the logger and parses command-line flags.
func init() {
	pflag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
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

		result, err = Merge(result, parsedData)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to merge YAML content from file: %s", file)
		}
	}

	mergedYAML, err := yaml.Marshal(result)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal merged result to YAML")
	}

	// Using strings.Builder to efficiently concatenate the final output
	var output strings.Builder
	output.WriteString(string(mergedYAML))
	fmt.Println(output.String())
}

// Merge recursively merges two YAML structures.
func Merge(a, b interface{}) (interface{}, error) {
	log.Debug().Msgf("Merging: %v (%T) with %v (%T)", a, a, b, b)
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
			if !found {
				typedA[key] = rightVal
			} else {
				mergedVal, err := Merge(leftVal, rightVal)
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
