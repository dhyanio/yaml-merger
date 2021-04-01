package main

// Dhyanio's Simple but powerful Yaml content merger
// Written in Beautiful Golang.
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var (
	ignoreList = pflag.StringSlice("ignore", []string{}, "keys to ignore when merging")
)

// Init - Before main function
func init() {
	pflag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msgf("usage: %s <file> [<file>...]", os.Args[0])
	}
	var res interface{}
	for _, f := range pflag.Args() {
		bs, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read file")
		}
		var part interface{}
		err = yaml.Unmarshal(bs, &part)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse file")
		}
		res, err = Merge(res, part)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to merge file")
		}
	}
	bs, err := yaml.Marshal(res)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal result")
	}
	fmt.Println(string(bs))
}

// Merger func for Yaml merging
func Merge(a, b interface{}) (_ interface{}, err error) {
	log.Debug().Msgf("merge %v (%T) %v (%T)", a, a, b, b)
	switch typedA := a.(type) {
	case []interface{}:
		typedB, ok := b.([]interface{})
		if !ok {
			return nil, errors.New("wrong type on right side")
		}
		return append(typedA, typedB...), nil
	case map[interface{}]interface{}:
		typedB, ok := b.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("wrong type on right side")
		}
		for key, rightVal := range typedB {
			if shouldIgnore(key) {
				continue
			}
			leftVal, ok := typedA[key]
			if !ok {
				typedA[key] = rightVal
			} else {
				typedA[key], err = Merge(leftVal, rightVal)
				if err != nil {
					return nil, err
				}
			}
		}
		return typedA, nil
	default:
		return b, nil
	}
	return nil, errors.New("unexpected end")
}

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
