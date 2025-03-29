package http_model

type RequestWelcomeWordParam struct {
}
type ResponseWelcomeWordParam struct {
	WelcomeWords     []string `json"welcome_word"`
	ToGuessQuestions []string `json:"to_guess_questions"`
}
