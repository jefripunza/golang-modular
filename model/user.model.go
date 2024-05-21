package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -> main collection
type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// info
	Name     string  `bson:"name"`
	ImageURL *string `bson:"image_url,omitempty"`

	// auth, unique:email
	Email *string `bson:"email,omitempty"` // login by email

	// auth, username-password
	Username string `bson:"username"`
	Password string `bson:"password"`
	IsVerify bool   `bson:"is_verify"` // jika login by email / baru, maka auto true

	// OTP
	OtpRef     *string             `bson:"otp_ref,omitempty"`
	OtpCode    *string             `bson:"otp_code,omitempty"`
	OtpExpired *primitive.DateTime `bson:"otp_expired,omitempty"`

	CreatedAt primitive.DateTime  `bson:"created_at"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty"`
}

type UserLoginHistory struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID    string `bson:"user_id"`
	UserAgent string `bson:"user_agent"`

	LoginAt primitive.DateTime `bson:"login_at"`
}

type UserRevoke struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID    string             `bson:"user_id"`
	JwtID     string             `bson:"jwt_id"`
	ExpiredAt primitive.DateTime `bson:"expired_at"`

	LoginAt primitive.DateTime `bson:"login_at"`
}

type UserMerchant struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID string `bson:"user_id"`

	// unique
	Username string `bson:"username"`

	// info
	Name  string `bson:"name"`
	Image string `bson:"image"`

	CreatedAt primitive.DateTime `bson:"created_at"`
}

// ------------------------------------------------------------------------

type UserRegisterBody struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterOtpBody struct {
	Ref  string `json:"ref"`
	Code string `json:"code"`
}

type UserLoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginByEmailBody struct {
	Email string `json:"email"`
}

type UserForgotPasswordBody struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

type UserForgotPasswordResendBody struct {
	Ref string `json:"ref"`
}

type UserForgotPasswordOtpValidBody struct {
	Ref  string `json:"ref"`
	Code string `json:"code"`
}

type UserForgotPasswordOtpSubmitBody struct {
	Ref      string `json:"ref"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
