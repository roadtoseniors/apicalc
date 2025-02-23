package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	http.Client
	Host string
	Port int
}

// запрашивам таску у оркестратора.
func (client *Client) GetTask() *task.Task {
	requesturl := fmt.Sprintf("http://%s:%d/internal/task", client.Host, client.Port)

	req, err := http.NewRequest(http.MethodGet, requesturl, nil)
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	compreq, err := client.Do(req.WithContext(ctx))
	if err != nil {.
		time.Sleep(500 * time.Millisecond)
		return nil
	}
	defer compreq.Body.Close()

	if compreq.StatusCode != http.StatusOK {
		return nil
	}

	answer := struct {
		Task task.Task `json:"task"`
	}{}

	err = json.NewDecoder(compreq.Body).Decode(&answer)
	if err != nil {
		return nil
	}

	return &answer.Task
}

// отправлям результат выполнения задачи оркестратору.
func (client *Client) SendResult(result result.Result) {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "    ")

	err := encoder.Encode(result)
	if err != nil {
		return
	}

	requesturl := fmt.Sprintf("http://%s:%d/internal/task", client.Host, client.Port)

	reqhttp, err := http.NewRequest(http.MethodPost, requesturl, &buf)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	compreq, err := client.Do(reqhttp.WithContext(ctx))
	if err != nil {
		return
	}
	defer compreq.Body.Close()
}
