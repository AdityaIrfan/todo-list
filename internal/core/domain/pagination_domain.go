package domain

type Pagination struct {
	CurrentPage    int   `json:"Current_Page"`
	MaxDataPerPage int   `json:"Max_Data_Per_Page"`
	MaxPage        int   `json:"Max_Page"`
	TotalAllData   int64 `json:"Total_All_Data"`
}
