package redcap

// Instrument represents a REDCap instrument/eCRF.
type Instrument struct {
	Name  string `json:"instrument_name"`
	Label string `json:"instrument_label"`
}

type Form struct {
	Name       string
	Fields     map[string]*Field
	FieldOrder []*Field
	Key        Field
}
