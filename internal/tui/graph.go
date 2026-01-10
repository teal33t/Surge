package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderMultiLineGraph creates a multi-line bar graph with grid lines.
// data: speed history data points
// width, height: dimensions of the graph
// maxVal: maximum value for scaling
// color: color for the data bars
func renderMultiLineGraph(data []float64, width, height int, maxVal float64, color lipgloss.Color) string {
	if width < 1 || height < 1 {
		return ""
	}

	// Styles
	gridStyle := lipgloss.NewStyle().Foreground(ColorGray) // Faint color for lines
	barStyle := lipgloss.NewStyle().Foreground(color)      // Bright color for data

	// 1. Prepare the canvas with a Grid
	rows := make([][]string, height)
	for i := range rows {
		rows[i] = make([]string, width)
		for j := range rows[i] {
			// Draw horizontal lines on every other row (or every row if you prefer)
			// Using "╌" gives a nice technical dashed look. "─" is a solid line.
			if i%2 == 0 {
				rows[i][j] = gridStyle.Render("╌")
			} else {
				rows[i][j] = " "
			}
		}
	}

	// 2. Slice data to fit width
	var visibleData []float64
	if len(data) > width {
		visibleData = data[len(data)-width:]
	} else {
		// If not enough data, we process what we have.
		// We do NOT pad with 0s here, because we want the grid
		// to show through on the left side.
		visibleData = data
	}

	// Block characters
	blocks := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	// 3. Draw Data Columns
	// We calculate the offset so data fills from the RIGHT
	offset := width - len(visibleData)

	for x, val := range visibleData {
		// Actual X position on the canvas
		canvasX := offset + x

		if val < 0 {
			val = 0
		}

		// Calculate height in "sub-blocks"
		pct := val / maxVal
		if pct > 1.0 {
			pct = 1.0
		}
		totalSubBlocks := pct * float64(height) * 8.0

		// Fill rows from bottom up
		for y := 0; y < height; y++ {
			rowIndex := height - 1 - y // 0 is top, height-1 is bottom

			// Calculate block value for this specific row
			rowValue := totalSubBlocks - float64(y*8)

			var char string
			if rowValue <= 0 {
				// No data for this height? Keep the grid background!
				continue
			} else if rowValue >= 8 {
				char = "█"
			} else {
				char = blocks[int(rowValue)]
			}

			// Overwrite the grid with the data bar
			rows[rowIndex][canvasX] = barStyle.Render(char)
		}
	}

	// 4. Join rows
	var s strings.Builder
	for i, row := range rows {
		s.WriteString(strings.Join(row, ""))
		if i < height-1 {
			s.WriteRune('\n')
		}
	}

	return s.String()
}
