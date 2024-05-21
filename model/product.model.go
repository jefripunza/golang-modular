package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//-> main collection
type Product struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// unique
	SeoURL string `bson:"seo_url"`

	// info
	Name        string `bson:"name"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	WeightGram  int    `bson:"weight_gram"`
	Qty         int    `bson:"qty"`
	Sold        int    `bson:"sold"`

	CategoryID string   `bson:"category_id"`
	EtalaseIDs []string `bson:"etalase_ids"`

	Images []ProductImage `bson:"images"`

	IsActive bool `bson:"is_active"`

	CreatedAt primitive.DateTime  `bson:"created_at"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty"`
}

type ProductImage struct {
	ProductID string `bson:"product_id"`

	URL       string `bson:"url"`
	IsPrimary bool   `bson:"is_primary"`

	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty"`
}

type ProductCategory struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name     string `bson:"name"`
	IconURL  string `bson:"icon_url"`
	IsActive bool   `bson:"is_active"`
}

type ProductEtalase struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// unique
	UserMerchantID string `bson:"user_merchant_id"`
	Name           string `bson:"name"`

	IsActive bool `bson:"is_active"`
}

type ProductReview struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// unique
	UserID        string `bson:"user_id"`
	InvoiceNumber string `bson:"invoice_number"`
	ProductID     string `bson:"product_id"`

	Comment string `bson:"comment"`
	Star    int    `bson:"star"`
}
