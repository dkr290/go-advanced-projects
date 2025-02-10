package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	torch "github.com/wangkuiyi/gotorch"
)

// loadModel loads a PyTorch model from a .pth file
func LoadModel(modelPath string) torch.Tensor {
	// Define the model architecture (must match the saved model)
	m := torch.Load(modelPath)

	return m
}

// Tokenization Function (Calls Python Tokenizer Service)
// the url is something like "http://localhost:5001/tokenize"
func TokenizeText(text string, tokenUrl string) ([]int64, error) {
	payload := map[string]string{"text": text}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		tokenUrl,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, _ := io.ReadAll(resp.Body)
	var response map[string][]int64
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}

	return response["tokens"], nil
}
