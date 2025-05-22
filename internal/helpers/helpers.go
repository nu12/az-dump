package helpers

import "strings"

func GetResourceGroupNameFromFileName(fileName string) string {
	return strings.Split(fileName, ".json")[0]
}

func ComaListContains(commaList, item string) bool {
	slice := strings.Split(commaList, ",")
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
