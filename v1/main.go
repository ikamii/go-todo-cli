package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Respresents a todo item
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

// Adds a task to the list
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
	fmt.Println("-------------------")
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
	return fmt.Errorf("Task with ID %d not found", id)
}

// Removes a task from the list
func (tl *TodoList) DeleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			fmt.Printf("Deleted task %d: %s\n", id, task.Title)
			return nil
		}
	}
	return fmt.Errorf("Task with ID %d not found", id)
}

// Saves the todo list to a JSON file
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

// Displays the available commands
func PrintHelp() {
	fmt.Println("Todo List Application")
	fmt.Println("---------------------")
	fmt.Println("Commands:")
	fmt.Println("  add <task>       - Add a new task")
	fmt.Println("  list             - List all tasks")
	fmt.Println("  complete <id>    - Mark a task as completed")
	fmt.Println("  delete <id>      - Delete a task")
	fmt.Println("  help             - Show this help message")
	fmt.Println("  exit             - Exit the application")
}

func main() {
	todoList := TodoList{}
	filename := "todo.json"

	// Load existing tasks from file
	err := todoList.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Todo List Application!")
	fmt.Println("Type 'help' for a list of commands.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		command := strings.ToLower(parts[0])

		switch command {
		case "add":
			if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
				fmt.Println("Error: Task description cannot be empty")
				continue
			}
			todoList.AddTask(parts[1])
			todoList.SaveToFile(filename)

		case "list":
			todoList.ListTasks()

		case "complete":
			if len(parts) < 2 {
				fmt.Println("Error: Please specify a task ID")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Error: Invalid task ID '%s'\n", parts[1])
				continue
			}
			err = todoList.CompleteTask(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			todoList.SaveToFile(filename)

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Error: Please specify a task ID")
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Error: Invalid task ID '%s'\n", parts[1])
				continue
			}
			err = todoList.DeleteTask(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			todoList.SaveToFile(filename)

		case "help":
			PrintHelp()

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
