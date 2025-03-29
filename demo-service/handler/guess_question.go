package handler

type GuessQuestioner interface {
	GetGuessQuestions(lang string) []string
}

type GuessQuestionFromDB struct {
}

// GetGuessQuestions
func (gDB *GuessQuestionFromDB) GetGuessQuestions(lang string) []string {
	return nil
}

type GuessQuestionFromConf struct {
}

// GetGuessQuestions
func (gDB *GuessQuestionFromConf) GetGuessQuestions(lang string) []string {
	return nil
}
