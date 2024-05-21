package module

import (
	"context"
	"core/config"
	"core/connection"
	"core/env"
	"core/model"
	"core/util"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

func (ref Auth) Route(api fiber.Router) {
	handler := AuthHandler{}
	route := api.Group("/auth")

	route.Post("/register", handler.Register)
	route.Post("/register-otp-verify", handler.RegisterOtpVerify)
	route.Post("/forgot-password", handler.ForgotPassword)
	route.Post("/forgot-password-resend-otp", handler.ForgotPasswordResendOTP)
	route.Post("/forgot-password-otp-valid", handler.ForgotPasswordOtpValid)
	route.Post("/forgot-password-submit", handler.ForgotPasswordSubmit)

	// JWT Handler with JTI
	route.Post("/login", handler.Login)
	route.Post("/login-by-google", handler.LoginByGoogle) // on callback
	route.Get("/token-validation", handler.TokenValidation)
	route.Get("/refresh-token", handler.RefreshToken)
}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type AuthHandler struct{}

func (handler AuthHandler) Register(c *fiber.Ctx) error {
	var body model.UserRegisterBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}

	err = collection.FindOne(ctx, bson.M{"username": body.Username, "is_verify": true}).Decode(&exist)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "username already exists",
		})
	} else if err != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error on query",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not hash password",
		})
	}

	NowAt := primitive.NewDateTimeFromTime(time.Now())

	err = collection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&exist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(ctx, model.User{
				Name: body.Name,

				Username: body.Username,
				Password: string(hashedPassword),
				IsVerify: false,

				CreatedAt: NowAt,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "cannot inserted",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error on query",
			})
		}
	} else {
		_, err = collection.UpdateOne(ctx, bson.M{
			"username": body.Username,
		}, bson.M{
			"$set": bson.M{
				"name": body.Name,

				"password": string(hashedPassword),

				"updated_at": NowAt,
			},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "cannot update",
			})
		}
	}

	otp_ref, otp_expired, err := createOTP(c, database, ctx, bson.M{
		"username": body.Username,
	})

	// send notification...

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "user pending",
		"otp_ref":     otp_ref,
		"otp_expired": otp_expired,
	})
}

func (handler AuthHandler) RegisterOtpVerify(c *fiber.Ctx) error {
	var body model.UserRegisterOtpBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}

	err = collection.FindOne(ctx, bson.M{
		"ref":  body.Ref,
		"code": body.Code,
	}).Decode(&exist)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error on query",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	NowAt := primitive.NewDateTimeFromTime(time.Now())

	_, err = collection.UpdateOne(ctx, bson.M{
		"ref":  body.Ref,
		"code": body.Code,
	}, bson.M{
		"$set": bson.M{
			"is_verify":  true,
			"updated_at": NowAt,
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot update",
		})
	}

	token, statusCode, err := generateToken(database, ctx, c.Get("user-agent"), exist.ID.Hex())
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = clearOTP(c, database, ctx, bson.M{
		"ref":  body.Ref,
		"code": body.Code,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
	})
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	var body model.UserLoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	// headers := c.Request().Header.RawHeaders()
	// fmt.Println("headers:", string(headers))

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}

	err = collection.FindOne(ctx, bson.M{
		"username": body.Username,
	}).Decode(&exist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "username or password is wrong 1",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	Encryption := util.Encryption{}
	decodePassword, err := Encryption.Decode(exist.Password)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error decryption",
		})
	}
	if body.Password != decodePassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "username or password is wrong 2",
		})
	}

	user_id := exist.ID.Hex()
	statusCode, err := checkJti(database, ctx, bson.M{
		"user_id": user_id,
	})
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, statusCode, err := generateToken(database, ctx, c.Get("user-agent"), user_id)
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (handler AuthHandler) LoginByGoogle(c *fiber.Ctx) error {
	url := config.GoogleOauthConfig.AuthCodeURL(config.GoogleOauthStateString)
	return c.Redirect(url)
}

func (handler AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var err error

	var body model.UserForgotPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}
	var otp_ref string
	var otp_code string
	if body.Email != nil {
		err = collection.FindOne(ctx, bson.M{
			"email": *body.Email,
		}).Decode(&exist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "account not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}
		otp_ref, otp_code, err = createOTP(c, database, ctx, bson.M{
			"email": *body.Email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	} else if body.Username != nil {
		err = collection.FindOne(ctx, bson.M{
			"username": *body.Username,
		}).Decode(&exist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "account not found",
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "internal server error",
				})
			}
		}
		otp_ref, otp_code, err = createOTP(c, database, ctx, bson.M{
			"username": *body.Username,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	} else if body.Email != nil && body.Username != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot use email and username same time",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// In a real application, you would send the OTP via email
	// For demonstration, we just return it in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "OTP sent",
		"otp_ref":  otp_ref,
		"otp_code": otp_code,
	})
}

func (handler AuthHandler) ForgotPasswordResendOTP(c *fiber.Ctx) error {
	var err error

	var body model.UserForgotPasswordResendBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request"})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}
	err = collection.FindOne(ctx, bson.M{
		"otp_ref": body.Ref,
	}).Decode(&exist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "account not found",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	otp_ref, otp_code, err := createOTP(c, database, ctx, bson.M{
		"otp_ref": body.Ref,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "OTP resent",
		"otp_ref":  otp_ref,
		"otp_code": otp_code,
	})
}

func (handler AuthHandler) ForgotPasswordOtpValid(c *fiber.Ctx) error {
	var err error

	var body model.UserForgotPasswordOtpValidBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}
	err = collection.FindOne(ctx, bson.M{
		"otp_ref":  body.Ref,
		"otp_code": body.Code,
	}).Decode(&exist)
	is_valid := true
	if err != nil {
		if err == mongo.ErrNoDocuments {
			is_valid = false
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}
	if is_valid {
		if exist.OtpExpired != nil && exist.OtpExpired.Time().Before(time.Now()) {
			is_valid = false
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_valid": is_valid,
	})
}

func (handler AuthHandler) ForgotPasswordSubmit(c *fiber.Ctx) error {
	var err error

	var body model.UserForgotPasswordOtpSubmitBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	collection := database.Collection("user")
	exist := model.User{}
	err = collection.FindOne(ctx, bson.M{
		"otp_ref":  body.Ref,
		"otp_code": body.Code,
	}).Decode(&exist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "account not found",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not hash password",
		})
	}

	_, err = collection.UpdateOne(ctx, bson.M{
		"otp_ref":  body.Ref,
		"otp_code": body.Code,
	}, bson.M{
		"$unset": bson.M{
			"password": string(hashedPassword),
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot clear otp",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "password reset successful",
	})
}

func (handler AuthHandler) TokenValidation(c *fiber.Ctx) error {
	claims, err := validateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}

	return c.Status(fiber.StatusOK).JSON(claims)
}

func (handler AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var err error

	claims, err := validateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}
	JWT := &util.JWT{}

	user_id := claims["user_id"].(string)
	newToken, jti, exp, err := JWT.Generate(user_id, time.Hour*1) // Short-lived token for access
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not generate token"})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "connect mongodb " + err.Error(),
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	statusCode, err := insertLogin(database, ctx, c.Get("user-agent"), user_id, jti, exp)
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": newToken,
	})
}

// ----------------------------------------------------------------

func insertLogin(database *mongo.Database, ctx context.Context, user_agent string, user_id string, jti string, exp time.Time) (int, error) {
	var err error

	user_revoke := database.Collection("user_revoke")
	user_login_history := database.Collection("user_login_history")
	ExpiredAt := primitive.NewDateTimeFromTime(exp)
	NowAt := primitive.NewDateTimeFromTime(time.Now())

	_, err = user_revoke.InsertOne(ctx, model.UserRevoke{
		UserID:    user_id,
		JwtID:     jti,
		ExpiredAt: ExpiredAt,
		LoginAt:   NowAt,
	})
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("cannot inserted user_revoke")
	}

	_, err = user_login_history.InsertOne(ctx, model.UserLoginHistory{
		UserID:    user_id,
		UserAgent: user_agent,
		LoginAt:   NowAt,
	})
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("cannot inserted user_login_history")
	}

	return 0, nil
}

func generateToken(database *mongo.Database, ctx context.Context, user_agent string, user_id string) (string, int, error) {
	var err error

	JWT := util.JWT{}
	token, jti, exp, err := JWT.Generate(user_id, time.Hour*1) // Short-lived token for access
	if err != nil {
		return "", fiber.StatusInternalServerError, fmt.Errorf("could not generate token")
	}

	statusCode, err := insertLogin(database, ctx, user_agent, user_id, jti, exp)
	if err != nil {
		return "", statusCode, err
	}
	return token, 0, nil
}

func validateToken(c *fiber.Ctx) (jwt.MapClaims, error) {
	JWT := &util.JWT{}

	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "missing authorization header"})
	}

	token := authorizationHeader[len("Bearer "):]
	if token == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "missing token"})
	}

	claims, err := JWT.Validate(token)
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}

	return claims, nil
}

func checkJti(database *mongo.Database, ctx context.Context, filter bson.M) (int, error) {
	var err error
	count, err := database.Collection("user_revoke").CountDocuments(ctx, filter)
	if err != nil {
		return fiber.StatusInternalServerError, fmt.Errorf("internal server error")
	}
	max_login_attempts := env.GetMaxLoginAttempts()
	if int(count) >= int(max_login_attempts) {
		return fiber.StatusForbidden, fmt.Errorf("maximum login attempts exceeded")
	}
	return 0, nil
}

func createOTP(c *fiber.Ctx, database *mongo.Database, ctx context.Context, filter bson.M) (string, string, error) {
	var err error

	Generate := util.Generate{}
	OtpRef := Generate.UUIDv4()
	OtpCode := Generate.OTP(6)
	otpExpired := primitive.NewDateTimeFromTime(time.Now().Add(5 * time.Minute))
	collection := database.Collection("user")
	_, err = collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"otp_ref":     OtpRef,
			"otp_code":    OtpCode,
			"otp_expired": otpExpired,
		},
	})
	if err != nil {
		return "", "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot create otp",
		})
	}
	return OtpRef, OtpCode, nil
}

func checkExpiredOTP(database *mongo.Database, ctx context.Context, filter bson.M) (bool, error) {
	var err error

	collection := database.Collection("user")
	var user model.User
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, err // No documents found
		}
		return false, err // Other errors
	}
	if user.OtpExpired != nil && user.OtpExpired.Time().Before(time.Now()) {
		return true, nil // OTP has expired
	}
	return false, nil // OTP is still valid
}

func clearOTP(c *fiber.Ctx, database *mongo.Database, ctx context.Context, filter bson.M) error {
	var err error

	collection := database.Collection("user")
	_, err = collection.UpdateOne(ctx, filter, bson.M{
		"$unset": bson.M{
			"otp_ref":     "",
			"otp_code":    "",
			"otp_expired": "",
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cannot clear otp",
		})
	}
	return nil
}
