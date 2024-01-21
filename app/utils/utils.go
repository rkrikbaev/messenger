package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// // Предположим, что data - это словарь, полученный от ParseSelectedFields
// data := map[string]float64{
// 	"ayraq_dev1.mass.Value": 0,
// 	"ayraq_dev2.level.Value": 55.25,
// 	"ayraq_dev3.level.Value": 251.1,
// }

// keyMap := map[string]string{
// 	"ayraq_dev1.mass.Value": "407001",
// 	"ayraq_dev2.level.Value": "407002",
// 	"ayraq_dev3.level.Value": "407003",
// }

type Data struct {
    Value     float64 `json:"Value"`
    Quality   int     `json:"Quality"`
    Timestamp string  `json:"Timestamp"`
}

// Функция для маппинга ключей
func MapKeys(data map[string]float64, keyMap map[string]string) map[string]float64 {
    mappedData := make(map[string]float64)
    for oldKey, newKey := range keyMap {
        if value, ok := data[oldKey]; ok {
            mappedData[newKey] = value
        }
    }
    return mappedData
}

// ParseSelectedFields парсит определенные поля из JSON.
func ParseFields(data []byte) (map[string]float64, error) {
    var rawData map[string]Data
    err := json.Unmarshal(data, &rawData)
    if err != nil {
        return nil, err
    }

    parsedData := make(map[string]float64)
    for key, value := range rawData {
        prefixedKey := fmt.Sprintf("%s.value", key)
        parsedData[prefixedKey] = value.Value
    }

    return parsedData, nil
}

// ReplaceKeys заменяет части ключей в словаре согласно предоставленной карте замен.
func ReplaceKeys(data map[string]float64, replacements map[string]string) map[string]float64 {
    replacedData := make(map[string]float64)
    for oldKey, value := range data {
        newKey := oldKey
        for oldPart, newPart := range replacements {
            newKey = strings.Replace(newKey, oldPart, newPart, -1)
        }
        replacedData[newKey] = value
    }
    return replacedData
}

func GetValidPrefix(fileName string, prefixes []string) (string, bool) {
    // Получение префикса из имени файла
    prefix := strings.Split(fileName, "_")[0]

    for _, validPrefix := range prefixes {
        if prefix == validPrefix {
            return prefix, true
        }
    }
    return "", false
}