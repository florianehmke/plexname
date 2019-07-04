package mock

import "github.com/florianehmke/plexname/prompt"

type AskNumberFn func(question string) (int, error)
type AskStringFn func(question string) (string, error)
type ConfirmFn func(question string) (bool, error)

func NewMockPrompter(askNumberFn AskNumberFn, askStringFn AskStringFn, confirmFn ConfirmFn) prompt.Prompter {
	return &mockPrompter{
		askNumberFn: askNumberFn,
		askStringFn: askStringFn,
		confirmFn:   confirmFn,
	}
}

type mockPrompter struct {
	askNumberFn AskNumberFn
	askStringFn AskStringFn
	confirmFn   ConfirmFn
}

func (p mockPrompter) AskNumber(question string) (int, error) {
	return p.askNumberFn(question)
}

func (p mockPrompter) AskString(question string) (string, error) {
	return p.askStringFn(question)
}

func (p mockPrompter) Confirm(question string) (bool, error) {
	return p.confirmFn(question)
}
