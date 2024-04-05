package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int          `json:"id,omitempty"`
	Description string       `json:"description"`
	Due         time.Time    `json:"due,omitempty"`
	Entry       time.Time    `json:"entry,omitempty"`
	Imask       float64      `json:"imask,omitempty"`
	Modified    time.Time    `json:"modified,omitempty"`
	Parent      string       `json:"parent,omitempty"`
	Project     string       `json:"project,omitempty"`
	Recur       string       `json:"recur,omitempty"`
	Reviewed    time.Time    `json:"reviewed,omitempty"`
	Rtype       string       `json:"rtype,omitempty"`
	Status      string       `json:"status,omitempty"`
	Until       time.Time    `json:"until,omitempty"`
	UUID        string       `json:"uuid,omitempty"`
	Wait        time.Time    `json:"wait,omitempty"`
	Annotations []Annotation `json:"annotations,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Urgency     float64      `json:"urgency,omitempty"`
}

type Annotation struct {
	Entry       time.Time `json:"entry,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		Due         string `json:"due,omitempty"`
		Entry       string `json:"entry,omitempty"`
		Modified    string `json:"modified,omitempty"`
		Reviewed    string `json:"reviewed,omitempty"`
		Until       string `json:"until,omitempty"`
		Wait        string `json:"wait,omitempty"`
		Annotations []struct {
			Entry string `json:"entry,omitempty"`
		} `json:"annotations"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	layout := "20060102T150405Z"
	if len(aux.Due) > 0 {
		dueTime, err := time.Parse(layout, aux.Due)
		if err != nil {
			return err
		}
		t.Due = dueTime
	}
	if len(aux.Entry) > 0 {
		entryTime, err := time.Parse(layout, aux.Entry)
		if err != nil {
			return err
		}
		t.Entry = entryTime
	}
	if len(aux.Modified) > 0 {
		modifiedTime, err := time.Parse(layout, aux.Modified)
		if err != nil {
			return err
		}
		t.Modified = modifiedTime
	}
	if len(aux.Reviewed) > 0 {
		reviewedTime, err := time.Parse(layout, aux.Reviewed)
		if err != nil {
			return err
		}
		t.Reviewed = reviewedTime
	}
	if len(aux.Until) > 0 {
		untilTime, err := time.Parse(layout, aux.Until)
		if err != nil {
			return err
		}
		t.Until = untilTime
	}
	if len(aux.Wait) > 0 {
		waitTime, err := time.Parse(layout, aux.Wait)
		if err != nil {
			return err
		}
		t.Wait = waitTime
	}
	for _, annotation := range aux.Annotations {
		if len(annotation.Entry) > 0 {
			entryTime, err := time.Parse(layout, annotation.Entry)
			if err != nil {
				return err
			}
			t.Annotations = append(t.Annotations, Annotation{
				Entry: entryTime,
			})
		}
	}
	return nil
}

func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task
	layout := "20060102T150405Z"
	var annotations []struct {
		Entry string `json:"entry,omitempty"`
	}
	for _, annotation := range t.Annotations {
		annotations = append(annotations, struct {
			Entry string `json:"entry,omitempty"`
		}{
			Entry: annotation.Entry.Format(layout),
		})
	}
	return json.Marshal(&struct {
		Due         *string `json:"due,omitempty"`
		Entry       *string `json:"entry,omitempty"`
		Modified    *string `json:"modified,omitempty"`
		Reviewed    *string `json:"reviewed,omitempty"`
		Until       *string `json:"until,omitempty"`
		Wait        *string `json:"wait,omitempty"`
		Annotations []struct {
			Entry string `json:"entry,omitempty"`
		} `json:"annotations,omitempty"`
		*Alias
	}{
		Due:         formatTime(t.Due, layout),
		Entry:       formatTime(t.Entry, layout),
		Modified:    formatTime(t.Modified, layout),
		Reviewed:    formatTime(t.Reviewed, layout),
		Until:       formatTime(t.Until, layout),
		Wait:        formatTime(t.Wait, layout),
		Annotations: annotations,
		Alias:       (*Alias)(&t),
	})
}

func formatTime(t time.Time, layout string) *string {
	if t.IsZero() {
		return nil
	}
	formatted := t.Format(layout)
	return &formatted
}

func mustParseTask(line string) Task {
	var task Task

	err := json.Unmarshal([]byte(line), &task)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	return task
}

func main() {
	// let's read from the standard input
	file := os.Stdin
	defer file.Close()

	// scann the input
	scanner := bufio.NewScanner(file)

	// Read the first Task
	if !scanner.Scan() {
		fmt.Println("Error reading standard input:")
		os.Exit(1)
	}
	line := scanner.Text()
	originalTask := mustParseTask(line)

	// Read the first Task
	if !scanner.Scan() {
		fmt.Println("Error reading standard input:")
		os.Exit(1)
	}
	line = scanner.Text()
	modifiedTask := mustParseTask(line)

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading standard input:", err)
		os.Exit(1)
	}

	result, err := json.Marshal(modifiedTask)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}

	// Exit with success
	fmt.Println(string(result))

	if originalTask.Status == "pending" && (modifiedTask.Status == "completed" || modifiedTask.Status == "deleted") {
		// Move the note related to this task into the archived folder (creating the folder if necessary)
		fmt.Printf("Task archived from %s to %s\n", originalTask.Status, modifiedTask.Status)
	}

	os.Exit(0)
}
