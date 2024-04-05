package httplogger

import (
	"encoding/json"
	"io"
	"net/http"
)

type Data struct {
    Value     float64 `json:"Value"`
    Quality   int     `json:"Quality"`
    Timestamp string  `json:"Timestamp"`
}

func FetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

func parseJSON(data []byte) (map[string]Data, error) {
    var result map[string]Data
    err := json.Unmarshal(data, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

