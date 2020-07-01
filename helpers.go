package main

import (
	"strconv"
	"strings"
)

func ListToInts(list string, delimiter string) ([]int, error) {
	var err error

	listSplit := strings.Split(list, delimiter)
	result := make([]int, len(listSplit))

	for index, value := range listSplit {
		if result[index], err = strconv.Atoi(value); err != nil {
			return result, err
		}
	}

	return result, nil
}
