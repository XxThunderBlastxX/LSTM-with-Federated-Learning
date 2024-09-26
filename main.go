package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WeightData struct {
	Weights  []interface{} `json:"weights"`
	ClientID string        `json:"clientId"`
}

var (
	weightArray = make(map[string][]interface{})
	mu          sync.Mutex
)

func main() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(handleWebSocket))

	log.Fatal(app.Listen(":1233"))
}

func handleWebSocket(c *websocket.Conn) {
	for {
		var data WeightData
		if err := c.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			break
		}

		mu.Lock()
		weightArray[data.ClientID] = data.Weights
		mu.Unlock()

		fmt.Printf("Received data from %s\n", data.ClientID)

		reply := calculateMeanWeights()

		if err := c.WriteJSON(reply); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func calculateMeanWeights() []interface{} {
	mu.Lock()
	defer mu.Unlock()

	if len(weightArray) == 0 {
		return nil
	}

	var allWeights [][]interface{}
	for _, weights := range weightArray {
		allWeights = append(allWeights, weights)
	}

	if len(allWeights) == 1 {
		return allWeights[0]
	}

	// Calculate mean weights
	meanWeights := make([]interface{}, len(allWeights[0]))
	for i := range meanWeights {
		switch w := allWeights[0][i].(type) {
		case []interface{}:
			meanWeights[i] = averageSlice(extractSlice(allWeights, i))
		case [][]interface{}:
			meanWeights[i] = average2DSlice(extract2DSlice(allWeights, i))
		default:
			log.Printf("Unexpected weight type at index %d: %T", i, w)
		}
	}

	fmt.Println("\nWeights received from clients")
	for clientID, weights := range weightArray {
		fmt.Printf("%s: %v (first weight)\n", clientID, weights[0])
	}

	return meanWeights
}

func extractSlice(weights [][]interface{}, index int) [][]float64 {
	result := make([][]float64, len(weights))
	for i, w := range weights {
		if slice, ok := w[index].([]interface{}); ok {
			result[i] = make([]float64, len(slice))
			for j, v := range slice {
				if f, ok := v.(float64); ok {
					result[i][j] = f
				}
			}
		}
	}
	return result
}

func extract2DSlice(weights [][]interface{}, index int) [][][]float64 {
	result := make([][][]float64, len(weights))
	for i, w := range weights {
		if slice2D, ok := w[index].([][]interface{}); ok {
			result[i] = make([][]float64, len(slice2D))
			for j, slice := range slice2D {
				result[i][j] = make([]float64, len(slice))
				for k, v := range slice {
					if f, ok := v.(float64); ok {
						result[i][j][k] = f
					}
				}
			}
		}
	}
	return result
}

func averageSlice(slices [][]float64) []float64 {
	if len(slices) == 0 || len(slices[0]) == 0 {
		return nil
	}
	result := make([]float64, len(slices[0]))
	for i := range result {
		sum := 0.0
		for _, s := range slices {
			sum += s[i]
		}
		result[i] = sum / float64(len(slices))
	}
	return result
}

func average2DSlice(slices [][][]float64) [][]float64 {
	if len(slices) == 0 || len(slices[0]) == 0 || len(slices[0][0]) == 0 {
		return nil
	}
	result := make([][]float64, len(slices[0]))
	for i := range result {
		result[i] = make([]float64, len(slices[0][i]))
		for j := range result[i] {
			sum := 0.0
			for _, s := range slices {
				sum += s[i][j]
			}
			result[i][j] = sum / float64(len(slices))
		}
	}
	return result
}
