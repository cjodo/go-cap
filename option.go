package redcap

import "strconv"

type ExportOptions struct {
	Records                []string // Record IDs
	Fields                 []string // Field names
	Forms                  []string // Form names
	Events                 []string // Event names (longitudinal)
	RawOrLabel             string   // "raw" | "label" | "both"
	RawOrLabelHeaders      string   // "raw" | "label"
	ExportCheckboxLabel    bool     // Export checkbox values as labels
	ExportSurveyFields     bool     // Include survey fields
	ExportDataAccessGroups bool     // Include DAG field
	Format                 string   // "json" | "csv" | "odm" | "xml"
	DateRangeBegin         string   // Date filter begin (ISO 8601)
	DateRangeEnd           string   // Date filter end
	FilterLogic            string   // Filter logic expression
}

type ImportOptions struct {
	Format            string // "json" | "csv" | "odm" | "xml" | "spss" | "r"
	Type              string // "flat" | "eav"
	OverwriteBehavior string // "normal" | "overwrite" | "Upsert"
	ForceAutoNumber   bool
	DateFormat        string // "YMD" | "MDY" | "DMY"
	ReturnContent     string // "count" | "ids" | "auto_ids"
}

type ImportResult struct {
	Count int      `json:"count"`
	IDs   []string `json:"ids,omitempty"`
	Error string   `json:"error,omitempty"`
}

// ExportOption is a functional option for ExportRecords.
type ExportOption func(map[string]string)

// ExportRecords filters records by the given record IDs.
func ExportRecordsFilter(records []string) ExportOption {
	return func(p map[string]string) {
		p["records"] = commaJoin(records)
	}
}

// ExportFields limits exported fields.
func ExportFields(fields []string) ExportOption {
	return func(p map[string]string) {
		p["fields"] = commaJoin(fields)
	}
}

// ExportForms limits exported forms.
func ExportForms(forms []string) ExportOption {
	return func(p map[string]string) {
		p["forms"] = commaJoin(forms)
	}
}

// ExportEvents filters by events (longitudinal).
func ExportEvents(events []string) ExportOption {
	return func(p map[string]string) {
		p["events"] = commaJoin(events)
	}
}

// ExportRawOrLabel sets raw/label option.
func ExportRawOrLabel(v string) ExportOption {
	return func(p map[string]string) {
		p["rawOrLabel"] = v
	}
}

// ExportFormat sets the export format.
func ExportFormat(format string) ExportOption {
	return func(p map[string]string) {
		p["format"] = format
	}
}

// ExportCheckboxLabel exports checkbox values as labels.
func ExportCheckboxLabel(b bool) ExportOption {
	return func(p map[string]string) {
		p["exportCheckboxLabel"] = strconv.FormatBool(b)
	}
}

// ExportSurveyFields includes survey fields.
func ExportSurveyFields(b bool) ExportOption {
	return func(p map[string]string) {
		p["exportSurveyFields"] = strconv.FormatBool(b)
	}
}

// ExportDataAccessGroups includes DAG field.
func ExportDataAccessGroups(b bool) ExportOption {
	return func(p map[string]string) {
		p["exportDataAccessGroups"] = strconv.FormatBool(b)
	}
}

// ExportFilterLogic applies filter logic.
func ExportFilterLogic(logic string) ExportOption {
	return func(p map[string]string) {
		p["filterLogic"] = logic
	}
}

// ImportOption is a functional option for ImportRecords.
type ImportOption func(map[string]string)

// ImportFormat sets the import format.
func ImportFormat(format string) ImportOption {
	return func(p map[string]string) {
		p["format"] = format
	}
}

// ImportOverwriteBehavior sets the overwrite behavior.
func ImportOverwriteBehavior(behavior string) ImportOption {
	return func(p map[string]string) {
		p["overwriteBehavior"] = behavior
	}
}

// ImportForceAutoNumber forces auto-numbering.
func ImportForceAutoNumber(b bool) ImportOption {
	return func(p map[string]string) {
		p["forceAutoNumber"] = strconv.FormatBool(b)
	}
}

// ImportReturnContent sets what to return.
func ImportReturnContent(content string) ImportOption {
	return func(p map[string]string) {
		p["returnContent"] = content
	}
}

func commaJoin(s []string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += "," + s[i]
	}
	return result
}
