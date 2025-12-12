package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#bd93f9") // Dracula Purple
	ColorSecondary = lipgloss.Color("#ff79c6") // Dracula Pink
	ColorSuccess   = lipgloss.Color("#50fa7b") // Dracula Green
	ColorError     = lipgloss.Color("#ff5555") // Dracula Red
	ColorWarning   = lipgloss.Color("#ffb86c") // Dracula Orange
	ColorText      = lipgloss.Color("#f8f8f2") // Dracula Foreground
	ColorSubtext   = lipgloss.Color("#6272a4") // Dracula Comment
	ColorBorder    = lipgloss.Color("#44475a") // Dracula Selection

	// Styles
	AppStyle = lipgloss.NewStyle().
			Padding(DefaultPaddingX, 2).
			Foreground(ColorText)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(DefaultPaddingY, DefaultPaddingX).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(DefaultPaddingY, DefaultPaddingX)

	FocusedPanelStyle = PanelStyle.
				BorderForeground(ColorSecondary)

	// List Styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Bold(true)

	ItemStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	// Status Bar Styles
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("#282a36")). // Dracula Background
			Padding(DefaultPaddingY, DefaultPaddingX)

	// Progress Styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	// New Styles for Refactor
	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true).
			Padding(DefaultPaddingY, DefaultPaddingX).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorPrimary).
			BorderBottom(true)

	// Stats Style in Header
	StatsStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Padding(DefaultPaddingY, DefaultPaddingX)

	// Base Card Style
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(DefaultPaddingY, DefaultPaddingX).
			Margin(DefaultPaddingY, DefaultPaddingX)

	// Selected Card Style (highlighted border)
	SelectedCardStyle = CardStyle.
				BorderForeground(ColorSecondary)

	// Text inside the card
	CardTitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	CardStatsStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Italic(true)
)
