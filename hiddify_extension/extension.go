package hiddify_extension

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hiddify/hiddify-core/config"
	"github.com/sagernet/sing-box/option"

	"github.com/fatih/color"
	ex "github.com/hiddify/hiddify-core/extension"
	ui "github.com/hiddify/hiddify-core/extension/ui"
)

// Color definitions for console output
var (
	red    = color.New(color.FgRed).Add(color.Bold)
	green  = color.New(color.FgGreen).Add(color.Underline)
	yellow = color.New(color.FgYellow)
)

// HiddifyExtension01Data holds the data specific to HiddifyExtension01
type HiddifyExtension01Data struct {
	Count int `json:"count"` // Number of counts for the extension
}

// Field name constants for easy reference, use similar name to the json key
const (
	CountKey = "count"
)

// HiddifyExtension01 represents the core functionality of the extension
type HiddifyExtension01 struct {
	ex.Base[HiddifyExtension01Data]                    // Embedding base extension functionality
	cancel                        context.CancelFunc // Function to cancel background tasks
	console                       string             // Stores console output
}

// GetUI returns the UI form for the extension
func (e *HiddifyExtension01) GetUI() ui.Form {
	// Create a form depending on whether there is a background task or not
	if e.cancel != nil {
		return ui.Form{
			Title:       "Hiddify-Extension-01",
			Description: "Awesome Extension hiddify_extension_01 created by vAHiD55555",
			Buttons:     []string{ui.Button_Cancel}, // Cancel button only when task is ongoing
			Fields: []ui.FormField{
				{
					Type:  ui.FieldConsole,
					Key:   "console",
					Label: "Console",
					Value: e.console, // Display console output
					Lines: 20,
				},
			},
		}
	}
	// Inital page
	return ui.Form{
		Title:       "Hiddify-Extension-01",
		Description: "Awesome Extension hiddify_extension_01 created by vAHiD55555",
		Buttons:     []string{ui.Button_Cancel, ui.Button_Submit},
		Fields: []ui.FormField{
			{
				Type:        ui.FieldInput,
				Key:         CountKey,
				Label:       "Count",
				Placeholder: "This will be the count",
				Required:    true,
				Value:       fmt.Sprintf("%d", e.Base.Data.Count), // Default value from stored data
				Validator:   ui.ValidatorDigitsOnly,               // Only allow digits
			},
			{
				Type:  ui.FieldConsole,
				Key:   "console",
				Label: "Console",
				Value: e.console, // Display current console output
				Lines: 20,
			},
		},
	}
}

// setFormData validates and sets the form data from input
func (e *HiddifyExtension01) setFormData(data map[string]string) error {
	// Check if CountKey exists in the provided data
	if val, ok := data[CountKey]; ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			// Validate that the count is greater than 5
			if intValue < 5 {
				return fmt.Errorf("please use a number greater than 5")
			} else {
				e.Base.Data.Count = intValue // Set valid count value
			}
		} else {
			return err // Return parsing error
		}
	}
	return nil // Return nil if data is set successfully
}

// backgroundTask runs a task in the background, updating the console at intervals
func (e *HiddifyExtension01) backgroundTask(ctx context.Context) {
	for count := 1; count <= e.Base.Data.Count; count++ {
		select {
		case <-ctx.Done(): // If context is done (cancel is pressed), exit the task
			e.cancel = nil
			e.addAndUpdateConsole(red.Sprint("Background Task Canceled")) // Notify cancellation
			return
		case <-time.After(1 * time.Second): // Wait for a second before the next iteration
			e.addAndUpdateConsole(red.Sprint(count), yellow.Sprint(" Background task ", count, " working..."))
		}
	}
	e.cancel = nil
	e.addAndUpdateConsole(green.Sprint("Background Task Finished Successfully")) // Task completion message
}

// addAndUpdateConsole adds messages to the console and updates the UI
func (e *HiddifyExtension01) addAndUpdateConsole(message ...any) {
	e.console = fmt.Sprintln(message...) + e.console // Prepend new messages to the console output
	e.UpdateUI(e.GetUI())                            // Update the UI with the new console content
}

// SubmitData processes and validates form submission data
func (e *HiddifyExtension01) SubmitData(data map[string]string) error {
	// Validate and set the form data
	err := e.setFormData(data)
	if err != nil {
		e.ShowMessage("Invalid data", err.Error()) // Show error message if data is invalid
		return err                                 // Return the error
	}
	// Cancel any ongoing background task
	if e.cancel != nil {
		e.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background()) // Create a new context for the task
	e.cancel = cancel                                       // Store the cancel function

	go e.backgroundTask(ctx) // Run the background task concurrently

	return nil // Return nil if submission is successful
}

// Cancel stops the ongoing background task if it exists
func (e *HiddifyExtension01) Cancel() error {
	if e.cancel != nil {
		e.cancel()     // Cancel the task
		e.cancel = nil // Clear the cancel function
	}
	return nil // Return nil after cancellation
}

// Stop is called when the extension is closed
func (e *HiddifyExtension01) Stop() error {
	return e.Cancel() // Simply delegate to Cancel
}

// To Modify user's config before connecting, you can use this function
func (e *HiddifyExtension01) BeforeAppConnect(hiddifySettings *config.HiddifyOptions, singconfig *option.Options) error {
	return nil
}

// NewHiddifyExtension01 initializes a new instance of HiddifyExtension01 with default values
func NewHiddifyExtension01() ex.Extension {
	return &HiddifyExtension01{
		Base: ex.Base[HiddifyExtension01Data]{
			Data: HiddifyExtension01Data{ // Set default data
				Count: 4, // Default count value
			},
		},
		console: yellow.Sprint("Welcome to ") + green.Sprint("Hiddify-Extension-01\n"), // Default message
	}
}

// init registers the extension with the provided metadata
func init() {
	ex.RegisterExtension(
		ex.ExtensionFactory{
			Id:          "github.com/vAHiD55555/hiddify_extension_01/hiddify_extension", // Package identifier
			Title:       "Hiddify-Extension-01",                                                         // Display title of the extension
			Description: "Awesome Extension hiddify_extension_01 created by vAHiD55555",                                                     // Brief description of the extension
			Builder:     NewHiddifyExtension01,                                                       // Function to create a new instance
		},
	)
}
