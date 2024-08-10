package utils

import "fmt"

func CheckEnvNotEmpty(env map[string]string, errTemplate string) {
	for envName, envValue := range env {
		if envValue == "" {
			errLog := fmt.Sprintf(errTemplate, envName)
			panic(errLog)
		}
	}
}
