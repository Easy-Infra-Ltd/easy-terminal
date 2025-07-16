# Terminal Task List Component

A thread-safe task list component for terminal applications built with [lipgloss](https://github.com/charmbracelet/lipgloss) and [bubbletea](https://github.com/charmbracelet/bubbletea).

## Features

- **Thread-safe**: Safe for use with goroutines
- **Visual indicators**: Different icons for pending (●), active (spinner), success (✓), and failed (✗) tasks
- **Progress tracking**: Support for progress bars or text-based status updates
- **Responsive**: Automatic text truncation for different terminal widths
- **Styled**: Sensible color scheme with customizable styles
- **Optional title**: Can display with or without a title/header

## Installation

```bash
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbletea
```

## Usage

### Basic Usage

```go
package main

import (
    "easy-terminal/tasklist"
)

func main() {
    // Create a new task list with title and max width
    taskList := tasklist.NewTaskList("Build Process", 80)
    
    // Add tasks
    task1 := taskList.AddTask("Initialize project")
    task2 := taskList.AddTask("Install dependencies")
    
    // Update task status
    task1.UpdateStatus(tasklist.TaskStatusSuccess)
    task1.SetMessage("Project initialized successfully")
    
    // Show progress
    task2.UpdateStatus(tasklist.TaskStatusActive)
    task2.SetProgressMode()
    task2.UpdateProgress(65, 100)
    
    // Render the task list
    fmt.Print(taskList.View())
}
```

### Task States

- `TaskStatusPending`: Task not yet started (● gray circle)
- `TaskStatusActive`: Task currently running (animated spinner)
- `TaskStatusSuccess`: Task completed successfully (✓ green checkmark)
- `TaskStatusFailed`: Task failed (✗ red cross)

### Display Modes

Each task can display in two modes:

1. **Text Mode**: Shows a status message
   ```go
   task.SetMessage("Processing files...")
   ```

2. **Progress Mode**: Shows a progress bar with percentage
   ```go
   task.SetProgressMode()
   task.UpdateProgress(current, total)
   ```

### Thread Safety

All methods are thread-safe and can be called from different goroutines:

```go
go func() {
    for i := 0; i <= 100; i++ {
        task.UpdateProgress(i, 100)
        time.Sleep(100 * time.Millisecond)
    }
    task.UpdateStatus(tasklist.TaskStatusSuccess)
    task.SetMessage("Task completed!")
}()
```

### Bubbletea Integration

The component is designed to work as a sub-component in bubbletea applications:

```go
type model struct {
    taskList *tasklist.TaskList
}

func (m model) View() string {
    return m.taskList.View()
}
```

### API Reference

#### TaskList Methods

- `NewTaskList(title string, maxWidth int) *TaskList` - Creates a new task list
- `AddTask(name string) *Task` - Adds a new task and returns a reference
- `GetTask(id string) *Task` - Retrieves a task by ID
- `GetTasks() []*Task` - Returns all tasks
- `SetMaxWidth(width int)` - Updates the maximum width for text truncation
- `View() string` - Renders the task list as a string

#### Task Methods

- `UpdateStatus(status TaskStatus)` - Updates the task status
- `UpdateProgress(current, total int)` - Updates progress and switches to progress mode
- `SetMessage(message string)` - Sets status message and switches to text mode
- `SetProgressMode()` - Switches to progress bar display
- `SetTextMode()` - Switches to text message display

## Examples

Run the examples:

```bash
# Simple console output
go run simple_example.go

# Interactive bubbletea demo (requires TTY)
go run example.go
```

## Color Scheme

- **Title**: Blue (bold)
- **Pending**: Gray
- **Active**: Yellow (spinning)
- **Success**: Green
- **Failed**: Red
- **Progress bars**: Gray
- **Success messages**: Green
- **Error messages**: Red
- **Other messages**: Gray