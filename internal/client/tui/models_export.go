package tui

// Re-export types from models package for backward compatibility with screens

import "bloop-tunnel/internal/client/tui/models"

type (
	InputFieldModel            = models.InputFieldModel
	InputFieldOpts             = models.InputFieldOpts
	InputFieldFocusedMsg       = models.InputFieldFocusedMsg
	InputFieldValueChangeMsg   = models.InputFieldValueChangeMsg

	SelectFieldModel           = models.SelectFieldModel
	SelectFieldOpts            = models.SelectFieldOpts

	ListViewModel              = models.ListViewModel
	ListItem                   = models.ListItem

	StatusModel                = models.StatusModel
	StatusLoadingMsg           = models.StatusLoadingMsg
	StatusSuccessMsg           = models.StatusSuccessMsg
	StatusErrorMsg             = models.StatusErrorMsg
)

// Re-export constructors
func NewInputField(opts InputFieldOpts) InputFieldModel {
	return models.NewInputField(models.InputFieldOpts{
		Label:       opts.Label,
		Placeholder: opts.Placeholder,
		Value:        opts.Value,
		Validation:   opts.Validation,
		IsPassword:   opts.IsPassword,
	})
}

func NewSelectField(opts SelectFieldOpts) SelectFieldModel {
	return models.NewSelectField(models.SelectFieldOpts{
		Label:    opts.Label,
		Options:  opts.Options,
		Selected: opts.Selected,
	})
}

func NewListViewModel(items []ListItem) *ListViewModel {
	itemsConv := make([]models.ListItem, len(items))
	for i, item := range items {
		itemsConv[i] = models.ListItem{
			ID:      item.ID,
			Label:   item.Label,
			Details: item.Details,
		}
	}
	m := models.NewListViewModel(itemsConv)
	return &m
}

func NewStatusModel() *StatusModel {
	return models.NewStatusModel()
}
