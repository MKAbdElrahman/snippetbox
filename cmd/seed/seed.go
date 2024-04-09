package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/model"
)

func main() {

	endpoint := "http://localhost:3000/api/snippet/create"
	data1 := model.NewSnippetParams{Title: "Snippet 1", Content: "The content of snippet 1", Expires: 3}
	data2 := model.NewSnippetParams{Title: "Snippet 2", Content: "Content for snippet 2", Expires: 4}

	err := sendPostRequest(endpoint, data1)
	if err != nil {
		fmt.Println("Error sending request 1:", err)
		return
	}

	err = sendPostRequest(endpoint, data2)
	if err != nil {
		fmt.Println("Error sending request 2:", err)
		return
	}

	fmt.Println("Successfully seeded database!")

}

func sendPostRequest(url string, data model.NewSnippetParams) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
