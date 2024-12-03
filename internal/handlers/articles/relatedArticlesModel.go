package articles

// Модель зв'язаних статей
type RelatedArticles struct {
	ParentArticleID  int `json:"parent_article_id" gorm:"primaryKey"`
	RelatedArticleID int `json:"related_article_id"`
}

func (RelatedArticles) TableName() string {
	return "related_articles"
}
