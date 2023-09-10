package data

type JsonResponse struct {
	Status   bool   `json:"status"`
	Info     string `json:"info"`
	Data     []User `json:"data"`
}
