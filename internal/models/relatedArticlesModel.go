package models

// Модель зв'язаних статей
type RelatedArticles struct {
	ParentArticleID  int `json:"parent_article_id" gorm:"column:parent_article_id;primaryKey"`
	RelatedArticleID int `json:"related_article_id" gorm:"column:related_article_id;type:int"`
}

func (RelatedArticles) TableName() string {
	return "related_articles"
}
