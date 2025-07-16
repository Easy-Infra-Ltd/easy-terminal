package tasklist

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusActive
	TaskStatusSuccess
	TaskStatusFailed
)

type TaskDisplayMode int

const (
	TaskDisplayText TaskDisplayMode = iota
	TaskDisplayProgress
)

type Task struct {
	ID          string
	Name        string
	Status      TaskStatus
	Message     string
	DisplayMode TaskDisplayMode
	Progress    int
	Total       int
	mu          sync.RWMutex
}

type TaskList struct {
	Title        string
	tasks        []*Task
	mu           sync.RWMutex
	maxWidth     int
	spinnerPos   int
	lastSpin     time.Time
	titleStyle   lipgloss.Style
	taskStyle    lipgloss.Style
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	warningStyle lipgloss.Style
	dimStyle     lipgloss.Style
}

func NewTaskList(title string, maxWidth int) *TaskList {
	return &TaskList{
		Title:    title,
		maxWidth: maxWidth,
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			MarginBottom(1),
		taskStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")),
		successStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")),
		warningStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")),
		dimStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")),
	}
}

func (tl *TaskList) AddTask(name string) *Task {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	task := &Task{
		ID:          uuid.New().String(),
		Name:        name,
		Status:      TaskStatusPending,
		DisplayMode: TaskDisplayText,
		Total:       100,
	}

	tl.tasks = append(tl.tasks, task)
	return task
}

func (t *Task) UpdateStatus(status TaskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = status
}

func (t *Task) UpdateProgress(current, total int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Progress = current
	t.Total = total
	t.DisplayMode = TaskDisplayProgress
}

func (t *Task) SetMessage(message string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Message = message
	t.DisplayMode = TaskDisplayText
}

func (t *Task) SetProgressMode() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.DisplayMode = TaskDisplayProgress
}

func (t *Task) SetTextMode() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.DisplayMode = TaskDisplayText
}

func (tl *TaskList) getSpinner() string {
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	now := time.Now()
	if now.Sub(tl.lastSpin) > time.Millisecond*100 {
		tl.spinnerPos = (tl.spinnerPos + 1) % len(spinnerChars)
		tl.lastSpin = now
	}

	return spinnerChars[tl.spinnerPos]
}

func (tl *TaskList) getStatusIcon(status TaskStatus) (string, lipgloss.Style) {
	switch status {
	case TaskStatusPending:
		return "●", tl.dimStyle
	case TaskStatusActive:
		return tl.getSpinner(), tl.warningStyle
	case TaskStatusSuccess:
		return "✓", tl.successStyle
	case TaskStatusFailed:
		return "✗", tl.errorStyle
	default:
		return "●", tl.dimStyle
	}
}

func (tl *TaskList) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

func (tl *TaskList) renderProgressBar(current, total int, width int) string {
	if width <= 0 {
		return ""
	}

	if total == 0 {
		return strings.Repeat("░", width) + " 0%"
	}

	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}

	filled := int(percentage * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percent := fmt.Sprintf(" %d%%", int(percentage*100))

	return bar + percent
}

func (tl *TaskList) View() string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	var output strings.Builder

	if tl.Title != "" {
		output.WriteString(tl.titleStyle.Render(tl.Title))
		output.WriteString("\n")
	}

	for _, task := range tl.tasks {
		task.mu.RLock()

		icon, iconStyle := tl.getStatusIcon(task.Status)
		iconStr := iconStyle.Render(icon)

		namePadding := 4
		nameWidth := tl.maxWidth - namePadding
		if nameWidth < 10 {
			nameWidth = 10
		}

		taskName := tl.truncateText(task.Name, nameWidth)

		line := fmt.Sprintf("%s %s", iconStr, taskName)

		if task.DisplayMode == TaskDisplayProgress {
			minWidth := 10
			progressWidth := tl.maxWidth - len(taskName) - minWidth
			if progressWidth > 0 {
				progressBar := tl.renderProgressBar(task.Progress, task.Total, progressWidth)
				line = fmt.Sprintf("%s %s %s", iconStr, taskName, tl.dimStyle.Render(progressBar))
			}
		} else if task.Message != "" {
			messagePadding := 6
			messageWidth := tl.maxWidth - len(taskName) - messagePadding
			if messageWidth > 0 {
				message := tl.truncateText(task.Message, messageWidth)
				var messageStyle lipgloss.Style
				switch task.Status {
				case TaskStatusSuccess:
					messageStyle = tl.successStyle
				case TaskStatusFailed:
					messageStyle = tl.errorStyle
				default:
					messageStyle = tl.dimStyle
				}
				line = fmt.Sprintf("%s %s %s", iconStr, taskName, messageStyle.Render(message))
			}
		}

		output.WriteString(line)
		output.WriteString("\n")

		task.mu.RUnlock()
	}

	return output.String()
}

func (tl *TaskList) GetTasks() []*Task {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	tasks := make([]*Task, len(tl.tasks))
	copy(tasks, tl.tasks)
	return tasks
}

func (tl *TaskList) GetTask(id string) *Task {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	for _, task := range tl.tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (tl *TaskList) SetMaxWidth(width int) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	tl.maxWidth = width
}

