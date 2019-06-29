package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Prompter interface {
	AskNumber(question string) (int, error)
	Confirm(question string) (bool, error)
}

type prompter struct {
	reader io.Reader
}

func NewPrompter() Prompter {
	return &prompter{reader: os.Stdin}
}

func (p *prompter) AskNumber(question string) (int, error) {
	r := bufio.NewReader(p.reader)
	fmt.Println(question)
	res, err := r.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("ask number prompt failed: %v", err)
	}
	n, err := strconv.Atoi(strings.ToLower(strings.TrimSpace(res)))
	if err != nil {
		fmt.Printf("invalid number: %v\n", err)
		return p.AskNumber(question)
	}
	return n, nil
}

func (p *prompter) Confirm(question string) (bool, error) {
	r := bufio.NewReader(p.reader)
	fmt.Println(question)
	res, err := r.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("confirm prompt failed: %v", err)
	}
	res = strings.ToLower(strings.TrimSpace(res))
	if res == "y" || res == "yes" {
		return true, nil
	} else if res == "n" || res == "no" {
		return false, nil
	}
	return p.Confirm(question)
}
