package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Represents a todo item
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Manages a list of tasks
type TodoList struct {
	Tasks  []Task `json:"tasks"`
	nextID int
}

// Adds a new task to the list
func (tl *TodoList) AddTask(title string) {
	task := Task{
		ID:        tl.nextID,
		Title:     title,
		Completed: false,
	}
	tl.Tasks = append(tl.Tasks, task)
	tl.nextID++
	fmt.Printf("Added task: %s (ID: %d)\n", title, task.ID)
}

// Prints all tasks in the list
func (tl *TodoList) ListTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("ID | Status | Task")
	fmt.Println("----------------------")
	for _, task := range tl.Tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("%2d | [%s]    | %s\n", task.ID, status, task.Title)
	}
}

// Marks a task as completed
func (tl *TodoList) CompleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks[i].Completed = true
			fmt.Printf("Marked task %d as completed: %s\n", id, task.Title)
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// Removes a task from the list
func (tl *TodoList) DeleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			// Remove the task by slicing it out
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			fmt.Printf("Deleted task %d: %s\n", id, task.Title)
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// SaveToFile saves the todo list to a JSON file
func (tl *TodoList) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// Loads the todo list from a JSON file
func (tl *TodoList) LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// If the file doesn't exist, start with an empty list
		if os.IsNotExist(err) {
			tl.Tasks = []Task{}
			tl.nextID = 1
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, tl)
	if err != nil {
		return err
	}

	// Find the highest ID to set nextID correctly
	maxID := 0
	for _, task := range tl.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	tl.nextID = maxID + 1

	return nil
}

func printUsage() {
	fmt.Println("Todo CLI - A simple task manager")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  todo [command] [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  add <task description>    Add a new task")
	fmt.Println("  list                      List all tasks")
	fmt.Println("  complete <task-id>        Mark a task as completed")
	fmt.Println("  delete <task-id>          Delete a task")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  todo add \"Buy groceries\"")
	fmt.Println("  todo list")
	fmt.Println("  todo complete 2")
	fmt.Println("  todo delete 3")
}

func main() {
	// Define command-line flags
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Set up todo list and data file
	todoList := TodoList{}
	filename := "todo.json"

	// Load existing tasks from file
	err := todoList.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
	}

	// Check if a command was provided
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Handle commands
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() < 1 {
			fmt.Println("Error: Task description required")
			return
		}
		// Collect all arguments as the task description
		taskDesc := strings.Join(os.Args[2:], " ")
		todoList.AddTask(taskDesc)
		todoList.SaveToFile(filename)

	case "list":
		listCmd.Parse(os.Args[2:])
		todoList.ListTasks()

	case "complete":
		completeCmd.Parse(os.Args[2:])
		if completeCmd.NArg() != 1 {
			fmt.Println("Error: Task ID required")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID '%s'\n", os.Args[2])
			return
		}
		err = todoList.CompleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		todoList.SaveToFile(filename)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if deleteCmd.NArg() != 1 {
			fmt.Println("Error: Task ID required")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: Invalid task ID '%s'\n", os.Args[2])
			return
		}
		err = todoList.DeleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		todoList.SaveToFile(filename)

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
	}
}
