package redcap

type Form struct {
	Name 				string
	Fields  		map[string]*Field
	FieldOrder 	[]*Field
	Key					Field
}
