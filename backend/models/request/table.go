package request

type CreateTableRequest struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}
