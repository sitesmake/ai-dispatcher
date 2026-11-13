package main

import (
	"errors"
	"fmt"
)

type Request struct {
	Question string
}

type Provider struct {
	Name  string
	Model string
}

func (p Provider) Ask(question string) (string, error) {
	if question == "" {
		return "", errors.New("empty question")
	}

	return "[" + p.Name + "/" + p.Model + "] " + "hello from api: " + question, nil
}

type OpenAI struct {
	Provider
}

type Anthropic struct {
	Provider
	APIVersion string
}

func (a Anthropic) Ask(question string) (string, error) {
	if question == "" {
		return "", errors.New("empty question")
	}

	return "[" + a.Name + "/" + a.Model + "] " + "hello from api: " + a.APIVersion + " " + question, nil
}

type Asker interface {
	Ask(question string) (string, error)
}

func useService(a Asker, question string) (string, error) {
	out, err := a.Ask(question)
	if err != nil {
		return "", err
	}
	return out, nil
}

func main() {
	request := Request{Question: "hello world!"}
	provider := Provider{Name: "openai", Model: "gpt-4o"}

	fmt.Println(request, provider)

	if _, err := provider.Ask(""); err != nil {
		fmt.Println("Ошибка:", err)
	}

	answer2, err := provider.Ask("whats time oclock?")

	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(answer2)
	}

	openAI := OpenAI{
		Provider: Provider{Name: "openai", Model: "gpt-5o"},
	}
	answer3, err := openAI.Ask("whats time oclock?")

	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(answer3)
	}

	anthropic := Anthropic{
		Provider:   Provider{Name: "anthropic", Model: "ai-1"},
		APIVersion: "v1.01",
	}
	answer4, err := anthropic.Ask("whats time oclock?")

	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(answer4)
	}

	answer5, err := useService(openAI, "new question 1")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(answer5)
	}

	answer6, err := useService(anthropic, "new question 2")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(answer6)
	}
}
