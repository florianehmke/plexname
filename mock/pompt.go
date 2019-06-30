package mock

import "github.com/florianehmke/plexname/prompt"

type AskNumberFn func(question string) (int, error)
type ConfirmFn func(question string) (bool, error)

func NewMockPrompter(askNumberFn AskNumberFn, confirmFn ConfirmFn) prompt.Prompter {
	return &mockPrompter{
		askNumberFn: askNumberFn,
		confirmFn:   confirmFn,
	}
}

type mockPrompter struct {
	askNumberFn AskNumberFn
	confirmFn   ConfirmFn
}

func (p mockPrompter) AskNumber(question string) (int, error) {
	return p.askNumberFn(question)
}

func (p mockPrompter) Confirm(question string) (bool, error) {
	return p.confirmFn(question)
}
