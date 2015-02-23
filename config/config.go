package config

// The config package parses configuration files
// for bao creation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gitbao/gitbao/model"
)

const (
	portRegexString    = "^\\s*PORT\\s*(\\d*).*$"
	envVarRegexString  = "^\\s*(\\S*?)\\s*=\\s*(\\S*).*$"
	commentRegexString = "^\\s*#.*$"
)

func Parse(config string) (response model.Config, err error) {

	portRegex := regexp.MustCompile(portRegexString)
	envVarRegex := regexp.MustCompile(envVarRegexString)

	lines := strings.Split(config, "\n")
	for _, line := range lines {
		var comment bool
		comment, err = regexp.MatchString(commentRegexString, line)
		if err != nil {
			return
		}
		if comment == true {
			continue
		}

		var port bool
		port, err = regexp.MatchString(portRegexString, line)
		if err != nil {
			return
		}
		if port == true {
			matches := portRegex.FindStringSubmatch(line)
			var portInt int
			portInt, err = strconv.Atoi(matches[1])
			if err != nil {
				return
			}
			response.Port = int64(portInt)
		}

		var envVar bool
		envVar, err = regexp.MatchString(envVarRegexString, line)
		if err != nil {
			return
		}
		if envVar == true {
			matches := envVarRegex.FindStringSubmatch(line)
			envVar := model.EnvVar{
				Key:   matches[1],
				Value: matches[2],
			}
			response.EnvVars = append(response.EnvVars, envVar)
		}
	}
	return
}
