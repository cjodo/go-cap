package redcap

type ExportOptions struct {
    Records                []string // Record IDs
    Fields                []string // Field names
    Forms                 []string // Form names
    Events                []string // Event names (longitudinal)
    RawOrLabel            string   // "raw" | "label" | "both"
    RawOrLabelHeaders    string   // "raw" | "label"
    ExportCheckboxLabel   bool     // Export checkbox values as labels
    ExportSurveyFields    bool     // Include survey fields
    ExportDataAccessGroups bool    // Include DAG field
    Format                string   // "json" | "csv" | "odm" | "xml"
    DateRangeBegin        string   // Date filter begin (ISO 8601)
    DateRangeEnd          string   // Date filter end
    FilterLogic           string   // Filter logic expression
}

type ImportOptions struct {
    Format            string   // "json" | "csv" | "odm" | "xml" | "spss" | "r"
    Type              string   // "flat" | "eav"
    OverwriteBehavior string   // "normal" | "overwrite" | "Upsert"
    ForceAutoNumber   bool
    DateFormat        string   // "YMD" | "MDY" | "DMY"
    ReturnContent     string   // "count" | "ids" | "auto_ids"
}

type ImportResult struct {
    Count     int      `json:"count"`
    IDs       []string `json:"ids,omitempty"`
    Error     string   `json:"error,omitempty"`
}
