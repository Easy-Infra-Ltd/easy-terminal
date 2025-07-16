package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"easy-terminal/tasklist"
)

type tickMsg time.Time

type model struct {
	taskList *tasklist.TaskList
	width    int
	height   int
}

func initialModel() model {
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
	
	return model{
		taskList: taskList,
		width:    80,
		height:   20,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.taskList.SetMaxWidth(m.width - 4)
		return m, nil
		
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			// Toggle task 2 between active and success
			task2 := m.taskList.GetTask(m.task2ID)
			if task2 != nil {
				if task2.Status == tasklist.TaskStatusActive {
					task2.UpdateStatus(tasklist.TaskStatusSuccess)
					task2.SetMessage("All dependencies installed")
				} else {
					task2.UpdateStatus(tasklist.TaskStatusActive)
					task2.SetProgressMode()
					task2.UpdateProgress(65, 100)
				}
			}
			return m, nil
		case "2":
			// Start task 3
			task3 := m.taskList.GetTask("task_3")
			if task3 != nil {
				task3.UpdateStatus(tasklist.TaskStatusActive)
				task3.SetMessage("Compiling...")
			}
			return m, nil
		case "3":
			// Update task 3 progress
			task3 := m.taskList.GetTask("task_3")
			if task3 != nil && task3.Status == tasklist.TaskStatusActive {
				task3.SetProgressMode()
				task3.UpdateProgress(task3.Progress+10, 100)
				if task3.Progress >= 100 {
					task3.UpdateStatus(tasklist.TaskStatusSuccess)
					task3.SetMessage("Compilation complete")
				}
			}
			return m, nil
		}
		
	case tickMsg:
		// Update progress for active tasks
		for _, task := range m.taskList.GetTasks() {
			if task.Status == tasklist.TaskStatusActive && task.DisplayMode == tasklist.TaskDisplayProgress {
				if task.Progress < task.Total {
					task.UpdateProgress(task.Progress+1, task.Total)
				}
			}
		}
		return m, tickCmd()
	}
	
	return m, nil
}

func (m model) View() string {
	header := "Task List Component Demo\n"
	header += "Press '1' to toggle task 2, '2' to start task 3, '3' to advance task 3, 'q' to quit\n\n"
	
	return header + m.taskList.View()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}