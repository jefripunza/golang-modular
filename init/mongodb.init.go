package initialize

import (
	"log"
	"project/connection"
	"project/env"
	"project/interfaces"

	"go.mongodb.org/mongo-driver/bson"
)

func MongoDB() {
	var err error

	MongoDB := connection.MongoDB{}

	client, ctx, err := MongoDB.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	// -----------------------------------------------------------------
	// -----------------------------------------------------------------

	MongoDB.CreateIndex(ctx, database, "user", []interfaces.IndexMongoDB{
		{
			Name:   "user_otp",
			Unique: false,
			Keys: bson.D{
				{Key: "otp_ref", Value: 1},
				{Key: "otp_code", Value: 1},
			},
		},
		{
			Name:   "user_unique_otp_ref",
			Unique: true,
			Keys: bson.D{
				{Key: "otp_ref", Value: 1},
			},
		},
		{
			Name:   "user_unique_otp_code",
			Unique: true,
			Keys: bson.D{
				{Key: "otp_code", Value: 1},
			},
		},
		{
			Name:   "user_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "username", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "user_revoke", []interfaces.IndexMongoDB{
		{
			Name:   "user_revoke_unique_per_item",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "user_merchant", []interfaces.IndexMongoDB{
		{
			Name:   "user_merchant_user_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Name:   "user_merchant_username",
			Unique: false,
			Keys: bson.D{
				{Key: "username", Value: 1},
			},
		},
		{
			Name:   "user_merchant_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "username", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "product", []interfaces.IndexMongoDB{
		{
			Name:   "product_user_merchant_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_merchant_id", Value: 1},
			},
		},
		{
			Name:   "product_seo_url",
			Unique: false,
			Keys: bson.D{
				{Key: "seo_url", Value: 1},
			},
		},
		{
			Name:   "product_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_merchant_id", Value: 1},
				{Key: "seo_url", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "wishlist", []interfaces.IndexMongoDB{
		{
			Name:   "wishlist_user_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Name:   "wishlist_product_id",
			Unique: false,
			Keys: bson.D{
				{Key: "product_id", Value: 1},
			},
		},
		{
			Name:   "wishlist_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "product_id", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "cart", []interfaces.IndexMongoDB{
		{
			Name:   "cart_user_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Name:   "cart_product_id",
			Unique: false,
			Keys: bson.D{
				{Key: "product_id", Value: 1},
			},
		},
		{
			Name:   "cart_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "product_id", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "payment_method", []interfaces.IndexMongoDB{
		{
			Name:   "payment_method_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "transaction", []interfaces.IndexMongoDB{
		{
			Name:   "transaction_user_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Name:   "transaction_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "invoice_number", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "report", []interfaces.IndexMongoDB{
		{
			Name:   "report_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "code", Value: 1},
			},
		},
	})

	// -----------------------------------------------------------------
	// -----------------------------------------------------------------
}
