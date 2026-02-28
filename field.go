package redcap

// Represents a REDCap data dictionary field.
type Field struct {
	Branching_logic                            string
	Custom_alignment                           string
	Field_label                                string
	Field_name                                 string
	Field_note                                 string
	Field_type                                 string
	Form_name                                  string
	Identifier                                 string
	Matrix_group_name                          string
	Matrix_ranking                             string
	Question_number                            string
	Required_field                             bool
	Section_header                             string
	Choices                                    []FieldChoice
	Calculations                               string
	Text_validation_max                        string
	Text_validation_min                        string
	Text_validation_type_or_show_slider_number string
	Value                                      string
}

type FieldChoice struct {
	ID    int
	Label string
}
