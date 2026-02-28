package redcap

type Event struct {
	Name 						string 	`json:"event_name"`		
	ArmNum 					int 		`json:"arm_num"`		
	DayOffset 			string 	`json:"day_offset"`		
	OffsetMin 			string 	`json:"offset_min"`		
	OffsetMax 			string 	`json:"offset_max"`		
	UniqueEventName string	`json:"unique_event_name"`
}
