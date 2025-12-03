package installer

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - professional cyan/blue theme
var (
	// Primary colors
	colorPrimary   = lipgloss.Color("86")  // Cyan
	colorSecondary = lipgloss.Color("39")  // Blue
	colorAccent    = lipgloss.Color("117") // Light cyan

	// Status colors
	colorSuccess = lipgloss.Color("42")  // Green
	colorWarning = lipgloss.Color("226") // Yellow
	colorMuted   = lipgloss.Color("241") // Gray

	// UI colors
	colorHighlight   = lipgloss.Color("117") // Light cyan for highlights
	colorBorder      = lipgloss.Color("86")  // Cyan for borders
	colorSelectedBg  = lipgloss.Color("24")  // Dark blue background
	colorSelectedFg  = lipgloss.Color("231") // White text
)

// Title style - bold, colored, padded
var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorPrimary).
	Padding(0, 1).
	MarginBottom(1)

// Border box style - rounded corners
var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder).
	Padding(1, 2).
	MarginTop(1).
	MarginBottom(1)

// Selected menu item style - highlighted background, contrasting text
var selectedItemStyle = lipgloss.NewStyle().
	Background(colorSelectedBg).
	Foreground(colorSelectedFg).
	Bold(true).
	Padding(0, 1)

// Normal menu item style
var normalItemStyle = lipgloss.NewStyle().
	Foreground(colorSecondary).
	Padding(0, 1)

// Cursor style - colored arrow
var cursorStyle = lipgloss.NewStyle().
	Foreground(colorHighlight).
	Bold(true)

// Help text style - muted/gray
var helpStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Italic(true).
	MarginTop(1)

// Prompt text style
var promptStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true)

// Banner style - for ASCII art header
var bannerStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	MarginBottom(1)

// ASCII art banner
const banner = `
   ╔═══════════════════════════════════════════════╗
   ║   🔧  C L A U D E   C O D E   F O U N D R Y   ║
   ╚═══════════════════════════════════════════════╝
`
