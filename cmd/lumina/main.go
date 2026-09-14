package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "health":
		checkHealth()
	case "ask":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a prompt. Usage: lumina ask \"Your question\"")
			os.Exit(1)
		}
		askAI(os.Args[2])
	case "prompt":
		if len(os.Args) < 4 {
			fmt.Println("Error: Please provide project ID and prompt. Usage: lumina prompt set \"<project>\" \"<prompt>\"")
			os.Exit(1)
		}
		if os.Args[2] == "set" {
			setPrompt(os.Args[3], os.Args[4])
		} else {
			fmt.Printf("Unknown prompt command: %s\n", os.Args[2])
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Lumina-Plane CLI")
	fmt.Println("Usage: lumina <command> [args]")
	fmt.Println("\nCommands:")
	fmt.Println("  health            Check platform and database health")
	fmt.Println("  ask \"<prompt>\"     Query the AI Gateway")
	fmt.Println("  prompt set \"<id>\" \"<template>\" Set a project prompt template")
}

func checkHealth() {
	url := "http://localhost:8000/health"
	fmt.Printf("Checking platform health at %s...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: Could not connect to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Platform Status: %s\n", string(body))
	} else {
		fmt.Printf("❌ Platform Unhealthy. Status: %d, Response: %s\n", resp.StatusCode, string(body))
	}
}

func askAI(prompt string) {
	url := "http://localhost:8000/ask"
	fmt.Printf("Querying AI Gateway: \"%s\"\n", prompt)

	reqBody := map[string]string{
		"project_id": "default-project",
		"prompt":     prompt,
		"model":      "llama3-8b-8192",
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error: AI Gateway connection failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("\n🤖 AI Response:\n%s\n", string(body))
	} else {
		fmt.Printf("❌ AI Error. Status: %d, Response: %s\n", resp.StatusCode, string(body))
	}
}

func setPrompt(projectID, template string) {
	url := "http://localhost:8000/prompt"
	fmt.Printf("Setting prompt template for project %s...\n", projectID)

	reqBody := map[string]interface{}{
		"project_id": projectID,
		"template":   template,
		"version":    1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error: Prompt API connection failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Prompt template set successfully!")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Failed to set prompt. Status: %d, Response: %s\n", resp.StatusCode, string(body))
	}
}
