package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mrcook/time_warrior/manager"
	"github.com/mrcook/time_warrior/reports"
)

var agendaCmd = &cobra.Command{
	Use:   "agenda [flags] PROJECT",
	Short: "Generate an agenda view of your time tracking data",
	Long: `Generate a Mermaid diagram showing your time tracking data in an agenda view.
This view helps visualize how your time was spent across different tasks and projects.

Time Periods:
  - t=today, w=this week, m=this month, y=this year.
  - 1d=yesterday, 1w=last week, 1m=last month, 1y=last year.

Project name omitted: agenda is generated showing all projects.
Project name specified: agenda is generated for that specific project.

Examples:
  $ tw agenda -p m
  => Agenda view showing all projects for the current month.

  $ tw agenda -p 1d MyProject
  => Agenda view for MyProject for yesterday.`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}
		generateAgenda(projectName, period)
	},
}

func init() {
	agendaCmd.Flags().StringVarP(&period, "period", "p", "", `report for the time period: t, 1d, w, m, y.`)
	rootCmd.AddCommand(agendaCmd)
}

func generateAgenda(projectName, period string) {
	m := manager.NewFromConfig(initializeConfig())

	report := reports.New(period)

	if projectName == "" {
		for _, filename := range m.AllProjectFilenames() {
			report.ProcessProjectFile(filename)
		}
	} else {
		filename, ok := m.ProjectFilename(projectName)
		if !ok {
			fmt.Println("project file not found")
			return
		}
		report.ProcessProjectFile(filename)
	}

	// Calculate total time spent for each day
	type DayTotal struct {
		date  string
		total float64
		fee   float64
	}
	dayTotals := make(map[string]DayTotal)
	var periodTotalFee float64

	// First pass: calculate total time for each day
	for _, p := range report.Projects() {
		for _, t := range p.SortedTasks() {
			startTime := time.Unix(int64(t.Started()), 0)
			day := startTime.Format("2006-01-02")
			dayTotal := dayTotals[day]
			dayTotal.total += t.TimeSpent()
			dayTotal.fee += t.Fee()
			dayTotals[day] = dayTotal
			periodTotalFee += t.Fee()
		}
	}

	// Generate Mermaid diagram
	fmt.Println("```mermaid")
	fmt.Println("gantt")
	fmt.Println("    title Time Tracking Agenda")
	fmt.Println("    dateFormat  YYYY-MM-DD HH:mm")
	fmt.Println("    axisFormat %H:%M")

	// Add work day bars and tasks
	currentDay := ""
	var firstDay, lastDay time.Time
	firstDaySet := false

	for _, p := range report.Projects() {
		for _, t := range p.SortedTasks() {
			startTime := time.Unix(int64(t.Started()), 0)
			endTime := time.Unix(int64(t.Finished()), 0)
			if t.Finished() == 0 {
				endTime = time.Now()
			}

			// Track first and last day for total fee bar
			if !firstDaySet {
				firstDay = startTime
				firstDaySet = true
			}
			lastDay = endTime

			day := startTime.Format("2006-01-02")

			// Add work day bar if it's a new day
			if day != currentDay {
				workStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 9, 0, 0, 0, time.Local)
				workEnd := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 17, 0, 0, 0, time.Local)
				fmt.Printf("    Work Day %.1fh %.0f$ :workday, %s, %s\n",
					dayTotals[day].total,
					dayTotals[day].fee,
					workStart.Format("2006-01-02 15:04"),
					workEnd.Format("2006-01-02 15:04"))
				currentDay = day
			}

			// Use actual time spent from task data
			timeSpent := t.TimeSpent()
			percentage := (timeSpent / dayTotals[day].total) * 100
			fee := t.Fee()

			// Format the task name to be Mermaid-compatible
			taskName := t.Name()
			if taskName == "." {
				taskName = p.Name()
			}

			fmt.Printf("    %s %.1fh %.1f%% %.0f$ :%s, %s, %s\n",
				taskName,
				timeSpent,
				percentage,
				fee,
				taskName,
				startTime.Format("2006-01-02 15:04"),
				endTime.Format("2006-01-02 15:04"))
		}
	}

	// Add total fee bar
	fmt.Printf("    Total Fee %.0f$ :totalfee, %s, %s\n",
		periodTotalFee,
		firstDay.Format("2006-01-02 15:04"),
		lastDay.Format("2006-01-02 15:04"))
	fmt.Println("```")
}
