package handle

type Article struct {
}
type DeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
}
