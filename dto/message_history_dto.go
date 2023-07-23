package dto

type MessageHistoryDTO struct {
	Data            []*MessageHistoryItemDTO `json:"data"`
	RecordsTotal    int64                    `json:"recordsTotal"`
	RecordsFiltered int                      `json:"recordsFiltered"`
}

type MessageHistoryItemDTO struct {
	Id      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}
