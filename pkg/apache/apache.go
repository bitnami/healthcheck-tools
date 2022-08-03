// Package apache provides functions for reading the Apache configuration files
package apache

import (
	"os"
	"path"
	"regexp"
)

// OpenApacheConfigurationFile opens a single apache configuration file and returns a string with the content
func OpenApacheConfigurationFile(confPath string) (string, error) {
	res := ""
	currentBuf, err := os.ReadFile(confPath)
	if err != nil {
		return res, err
	}
	res = string(currentBuf)
	return res, err
}

// GetIncludes parses a string and obtains the path to all the included Apache files
func GetIncludes(text, apacheRoot string) []string {
	includeRe := regexp.MustCompilePOSIX("^[[:space:]]*Include [\"]?([^\n\"]+)[\"]?")
	includeMatches := includeRe.FindAllStringSubmatch(text, -1)
	res := []string{}
	for _, element := range includeMatches {
		newFile := element[1]
		if !path.IsAbs(newFile) {
			newFile = path.Join(apacheRoot, newFile)
		}
		res = append(res, newFile)
	}
	return res
}

// OpenAllApacheConfigurationFiles opens an apache configuration file (and all the included ones) and returns their content as a map of <path to apache file>:<content of apache file>
func OpenAllApacheConfigurationFiles(confPath, apacheRoot string) (map[string]string, error) {
	remainingConfFiles := []string{confPath}
	resBuffers := make(map[string]string)
	for len(remainingConfFiles) > 0 {
		currentConfFile := remainingConfFiles[0]
		remainingConfFiles = remainingConfFiles[1:]
		if _, ok := resBuffers[currentConfFile]; ok {
			continue
		}
		if _, err := os.Stat(currentConfFile); err == nil {
			bufferString, err := OpenApacheConfigurationFile(currentConfFile)
			if err != nil {
				return nil, err
			}
			resBuffers[currentConfFile] = bufferString
			includedFilePaths := GetIncludes(bufferString, apacheRoot)
			remainingConfFiles = append(remainingConfFiles, includedFilePaths...)
		} else {
			return nil, err
		}
	}
	return resBuffers, nil
}
