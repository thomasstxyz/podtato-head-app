package version

import (
	"os"
)

const (
	defaultServiceVersion string = "v0.1.0"
	defaultPartNumber     string = "01"
)

func ServiceVersion() string {
	dynamicVersion, found := os.LookupEnv("PODTATO_VERSION")
	if !found || dynamicVersion == "" {
		return defaultServiceVersion
	}
	return dynamicVersion
}

func PartNumber(part string) string {
	var key string
	switch part {
	case "hat":
		key = "PODTATO_HAT_PART_NUMBER"
	case "left-arm":
		key = "PODTATO_LEFT_ARM_PART_NUMBER"
	case "right-arm":
		key = "PODTATO_RIGHT_ARM_PART_NUMBER"
	case "left-leg":
		key = "PODTATO_LEFT_LEG_PART_NUMBER"
	case "right-leg":
		key = "PODTATO_RIGHT_LEG_PART_NUMBER"
	}
	dynamicPartNumber, found := os.LookupEnv(key)
	if dynamicPartNumber == "" {
		dynamicPartNumber, found = os.LookupEnv("PODTATO_PART_NUMBER")
	}
	if !found || dynamicPartNumber == "" {
		return defaultPartNumber
	}
	return dynamicPartNumber
}
