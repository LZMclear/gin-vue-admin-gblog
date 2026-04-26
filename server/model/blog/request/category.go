package request

type CategoryUpsert struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"categoryName"`
}
