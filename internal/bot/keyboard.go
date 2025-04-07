package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Categories represents the available categories for interviews
var Categories = []string{
	"Algorithms",
	"SystemDesign",
	"DataScience",
	"Frontend",
	"ML",
	"DevOps",
	"DBA",
}

// ProgrammingLanguages represents the available programming languages
var ProgrammingLanguages = []string{
	"Go",
	"Cpp",
	"C",
	"C#",
	"Rust",
	"JS",
	"Java",
	"Ruby",
	"Python",
	"Kotlin",
	"Swift",
	"PHP",
}

// OtherCategories represents additional categories
var OtherCategories = []string{
	"Cybersecurity",
	"QA",
}

// ExperienceLevels represents the available experience levels
var ExperienceLevels = []string{
	"Intern",
	"Junior",
	"Middle",
	"Senior",
}

// CreateCategoriesKeyboard creates a keyboard with categories and programming languages
func CreateCategoriesKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	
	// Add main categories in pairs
	for i := 0; i < len(Categories); i += 2 {
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(Categories[i], "category:"+Categories[i]))
		
		if i+1 < len(Categories) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(Categories[i+1], "category:"+Categories[i+1]))
		}
		
		rows = append(rows, row)
	}
	
	// Add programming languages in groups of 3
	for i := 0; i < len(ProgrammingLanguages); i += 3 {
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(ProgrammingLanguages[i], "category:"+ProgrammingLanguages[i]))
		
		if i+1 < len(ProgrammingLanguages) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(ProgrammingLanguages[i+1], "category:"+ProgrammingLanguages[i+1]))
		}
		
		if i+2 < len(ProgrammingLanguages) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(ProgrammingLanguages[i+2], "category:"+ProgrammingLanguages[i+2]))
		}
		
		rows = append(rows, row)
	}
	
	// Add other categories
	var otherRow []tgbotapi.InlineKeyboardButton
	for _, category := range OtherCategories {
		otherRow = append(otherRow, tgbotapi.NewInlineKeyboardButtonData(category, "category:"+category))
	}
	rows = append(rows, otherRow)
	
	// Add "Not found" option
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Category not found", "category:notfound"),
	})
	
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// CreateLevelsKeyboard creates a keyboard with experience levels
func CreateLevelsKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	
	// Add experience levels in a single row each
	for _, level := range ExperienceLevels {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(level, "level:"+level),
		})
	}
	
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}