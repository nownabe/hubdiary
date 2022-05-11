package main

import (
	"fmt"
	"os"
	"os/exec"
)

type editor struct {
	editor string
}

func (e *editor) Edit(text string) (string, error) {
	tmpFile, err := os.CreateTemp("", "*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer removeTempFile(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(text)); err != nil {
		return "", fmt.Errorf("failed to write text to tempfile: %w", err)
	}

	cmd := exec.Command(e.editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run editor: %w", err)
	}

	body, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return string(body), nil
}

func removeTempFile(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}
}
