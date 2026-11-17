package main

import "fmt"

const fallbackProvider = "anthropic"

// TODO (шаг 1): объяви здесь свои ошибки — именованные и кастомную.

type Provider struct {
	Name    string
	Model   string
	Enabled bool
	Status  int
}

// check проверяет, может ли провайдер ответить на вопрос.
//
// TODO (шаг 1): верни ошибку в ситуациях 2–6 из задания:
//   - вопрос пустой;
//   - провайдер выключен (Enabled == false);
//   - API ответил не 200 (Status != 200).
func (p Provider) check(question string) error {
	return nil
}

func (p Provider) Ask(question string) (string, error) {
	if err := p.check(question); err != nil {
		return "", err
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
	if err := a.check(question); err != nil {
		return "", err
	}

	return "[" + a.Name + "/" + a.Model + "] " + "hello from api: " + a.APIVersion + " " + question, nil
}

type Asker interface {
	Ask(question string) (string, error)
}

// findProvider ищет провайдера по имени.
//
// TODO (шаг 1): если провайдера с таким именем нет, верни ошибку (ситуация 1).
// Сейчас функция в этом случае возвращает nil без ошибки — из-за этого
// программа падает с паникой на запросе 3.
func findProvider(providers map[string]Asker, name string) (Asker, error) {
	return providers[name], nil
}

// dispatch находит провайдера и задаёт ему вопрос.
func dispatch(providers map[string]Asker, name string, question string) (string, error) {
	provider, err := findProvider(providers, name)
	if err != nil {
		// TODO (шаг 2): реши, что делать с этой ошибкой — обработать здесь
		// или вернуть выше. Если возвращаешь — добавь контекст.
		return "", err
	}

	answer, err := provider.Ask(question)
	if err != nil {
		// TODO (шаг 2): здесь могут оказаться разные ошибки.
		// Одну из них dispatch обрабатывает сам: если провайдер выключен,
		// нужно задать тот же вопрос провайдеру fallbackProvider (ситуация 3).
		// Если и он не ответил — это ситуация 7.
		// Остальные ошибки верни выше с контекстом.
		return "", err
	}

	return answer, nil
}

// react получает ошибку от dispatch и возвращает строку-реакцию для main.
//
// TODO (шаг 3): определи, что это за ошибка, и верни реакцию из таблицы
// в задании (ситуации 1, 2, 4, 5, 6). Для всех остальных ошибок оставь
// реакцию ниже — это ситуация 7.
func react(err error) string {
	return "не удалось обработать запрос: " + err.Error()
}

// В main ничего менять не нужно: провайдеры и запросы подобраны так,
// чтобы сработала каждая ситуация из задания.
func main() {
	providers := map[string]Asker{
		"openai":    OpenAI{Provider{Name: "openai", Model: "gpt-4o", Enabled: true, Status: 200}},
		"anthropic": Anthropic{Provider{Name: "anthropic", Model: "ai-1", Enabled: true, Status: 200}, "v1.01"},
		"deepseek":  OpenAI{Provider{Name: "deepseek", Model: "chat", Enabled: false, Status: 200}},
		"grok":      OpenAI{Provider{Name: "grok", Model: "grok-4", Enabled: true, Status: 429}},
		"gemini":    OpenAI{Provider{Name: "gemini", Model: "flash", Enabled: true, Status: 401}},
		"mistral":   OpenAI{Provider{Name: "mistral", Model: "large", Enabled: true, Status: 503}},
	}

	requests := []struct {
		provider string
		question string
	}{
		{"openai", "what is a goroutine?"}, // успешный запрос
		{"openai", ""},                     // ситуация 2: пустой вопрос
		{"yandex", "hello"},                // ситуация 1: провайдера нет
		{"deepseek", "what is a channel?"}, // ситуация 3: провайдер выключен
		{"grok", "question 5"},             // ситуация 4: код 429
		{"gemini", "question 6"},           // ситуация 5: код 401
		{"mistral", "question 7"},          // ситуация 6: другой код
	}

	for i, r := range requests {
		answer, err := dispatch(providers, r.provider, r.question)
		if err != nil {
			fmt.Printf("%d: %s\n", i+1, react(err))
			continue
		}
		fmt.Printf("%d: %s\n", i+1, answer)
	}

	// Ситуация 7: deepseek выключен, и запасной anthropic тоже выключен.
	noFallback := map[string]Asker{
		"deepseek":  OpenAI{Provider{Name: "deepseek", Model: "chat", Enabled: false, Status: 200}},
		"anthropic": Anthropic{Provider{Name: "anthropic", Model: "ai-1", Enabled: false, Status: 200}, "v1.01"},
	}

	answer, err := dispatch(noFallback, "deepseek", "question 8")
	if err != nil {
		fmt.Printf("8: %s\n", react(err))
		return
	}
	fmt.Printf("8: %s\n", answer)
}
