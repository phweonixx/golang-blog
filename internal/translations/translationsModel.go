package translations

// Модель для перекладів
type Translations struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	ObjectID int    `json:"object_id"`
	Field    string `json:"field"`
	Language string `json:"language"`
	Content  string `json:"content"`
}
