package redcap

type Record struct {
	ID string
	// FieldName->Value
	Fields     map[string]any
	EventName  string
	Repetition FormRepetition
}

type FormRepetition struct {
	FormName          string
	CustomRecordLabel string
}
