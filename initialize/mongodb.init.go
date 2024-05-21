package initialize

import (
	"core/connection"
	"core/env"
	"core/interfaces"
	"core/model"
	"core/util"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	collection := database.Collection("user")
	exist := model.User{}

	err = collection.FindOne(ctx, bson.M{
		"username": "admin",
	}).Decode(&exist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Encryption := util.Encryption{}
			hashedPassword, err := Encryption.Encode("admin123")
			if err != nil {
				log.Fatalln("could not hash password")
				return
			}
			NowAt := primitive.NewDateTimeFromTime(time.Now())
			_, err = collection.InsertOne(ctx, model.User{
				Name:      "Administrator",
				Username:  "admin",
				Password:  string(hashedPassword),
				IsVerify:  true,
				CreatedAt: NowAt,
			})
			if err != nil {
				log.Fatalln("cannot inserted")
				return
			}
			log.Println("✅ user admin created!")
		}
	}

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

	MongoDB.CreateIndex(ctx, database, "user_login_history", []interfaces.IndexMongoDB{
		{
			Name:   "user_login_history_unique_per_item",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
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

	MongoDB.CreateIndex(ctx, database, "product_category", []interfaces.IndexMongoDB{
		{
			Name:   "product_category_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "product_etalase", []interfaces.IndexMongoDB{
		{
			Name:   "product_etalase_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_merchant_id", Value: 1},
				{Key: "name", Value: 1},
			},
		},
	})

	MongoDB.CreateIndex(ctx, database, "product_review", []interfaces.IndexMongoDB{
		{
			Name:   "product_review_user_id",
			Unique: false,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Name:   "product_review_invoice_number",
			Unique: false,
			Keys: bson.D{
				{Key: "invoice_number", Value: 1},
			},
		},
		{
			Name:   "product_review_product_id",
			Unique: false,
			Keys: bson.D{
				{Key: "product_id", Value: 1},
			},
		},
		{
			Name:   "product_review_unique_per_item",
			Unique: true,
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "invoice_number", Value: 1},
				{Key: "product_id", Value: 1},
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
			Name:   "transaction_status",
			Unique: false,
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
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
