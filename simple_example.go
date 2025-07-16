package main

import (
	"fmt"
	"time"

	"easy-terminal/tasklist"
)

func main() {
	// Create a new task list
	taskList := tasklist.NewTaskList("Build Process", 80)
	
	// Add some example tasks
	task1 := taskList.AddTask("Initialize project")
	task2 := taskList.AddTask("Install dependencies") 
	task3 := taskList.AddTask("Compile source code")
	task4 := taskList.AddTask("Run tests")
	task5 := taskList.AddTask("Package application")
	
	// Simulate some task states
	task1.UpdateStatus(tasklist.TaskStatusSuccess)
	task1.SetMessage("Project initialized successfully")
	
	task2.UpdateStatus(tasklist.TaskStatusActive)
	task2.SetProgressMode()
	task2.UpdateProgress(65, 100)
	
	task3.UpdateStatus(tasklist.TaskStatusPending)
	
	task4.UpdateStatus(tasklist.TaskStatusFailed)
	task4.SetMessage("Test suite failed: 3 failures")
	
	task5.UpdateStatus(tasklist.TaskStatusPending)
	
	// Print the task list
	fmt.Println("Initial state:")
	fmt.Print(taskList.View())
	
	// Simulate progress updates
	fmt.Println("\nAfter some progress:")
	for i := 0; i < 10; i++ {
		task2.UpdateProgress(65+i*3, 100)
		time.Sleep(100 * time.Millisecond)
	}
	
	task2.UpdateStatus(tasklist.TaskStatusSuccess)
	task2.SetMessage("All dependencies installed")
	
	task3.UpdateStatus(tasklist.TaskStatusActive)
	task3.SetProgressMode()
	task3.UpdateProgress(50, 100)
	
	fmt.Print(taskList.View())
	
	// Test truncation
	fmt.Println("\nWith truncation (smaller width):")
	taskList.SetMaxWidth(40)
	longTask := taskList.AddTask("This is a very long task name that should be truncated")
	longTask.SetMessage("This is also a very long message that should be truncated")
	
	fmt.Print(taskList.View())
}