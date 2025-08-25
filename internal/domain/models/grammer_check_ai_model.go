package models

type Correction struct {
	OriginalPhrase  string `json:"original_phrase"`
	CorrectedPhrase string `json:"corrected_phrase"`
	Explanation     string `json:"explanation"`
}

type GrammarResponse struct {
	CorrectedText string       `json:"corrected_text"`
	Corrections   []Correction `json:"corrections"`
}
