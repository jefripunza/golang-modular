package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Status int

// Status constants
const (
	pending_payment Status = iota
	failed_payment
	awaiting_shipment
	cancelled
	processing
	shipped
	in_transit
	out_for_delivery
	delivered
	returned
	refunded
	pending_review
	hold
	success
)

// StatusEnum struct
var StatusEnum = struct {
	PendingPayment   Status
	FailedPayment    Status
	AwaitingShipment Status
	Cancelled        Status
	Processing       Status
	Shipped          Status
	InTransit        Status
	OutForDelivery   Status
	Delivered        Status
	Returned         Status
	Refunded         Status
	PendingReview    Status
	Hold             Status
	Success          Status
}{
	PendingPayment:   pending_payment,
	FailedPayment:    failed_payment,
	AwaitingShipment: awaiting_shipment,
	Cancelled:        cancelled,
	Processing:       processing,
	Shipped:          shipped,
	InTransit:        in_transit,
	OutForDelivery:   out_for_delivery,
	Delivered:        delivered,
	Returned:         returned,
	Refunded:         refunded,
	PendingReview:    pending_review,
	Hold:             hold,
	Success:          success,
}

// String method for Status type
func (s Status) String() string {
	return [...]string{
		"PENDING_PAYMENT",   //+1 Pembayaran dari pelanggan belum diterima atau dikonfirmasi.
		"FAILED_PAYMENT",    //-1 Pembayaran dari pelanggan gagal, misalnya karena dana tidak mencukupi atau kesalahan teknis.
		"AWAITING_SHIPMENT", //+2 Pembayaran berhasil diterima, dan pesanan sedang menunggu untuk diproses dan dikirim.
		"CANCELLED",         //-2 Pesanan dibatalkan oleh pelanggan atau penjual sebelum pengiriman.
		"PROCESSING",        //+3 Pesanan sedang diproses untuk pengepakan dan pengiriman.
		"SHIPPED",           //+4 Pesanan telah diproses dan dikirim ke alamat pelanggan.
		"IN_TRANSIT",        //+5 Pesanan sedang dalam perjalanan menuju alamat pelanggan.
		"OUT_FOR_DELIVERY",  //+6 Pesanan sedang dalam proses pengiriman terakhir menuju alamat pelanggan.
		"DELIVERED",         //+7 Pesanan telah berhasil sampai dan diterima oleh pelanggan.
		"RETURNED",          //-7 Pesanan telah dikembalikan oleh pelanggan, mungkin karena cacat atau tidak sesuai dengan harapan.
		"REFUNDED",          //-7 Dana telah dikembalikan ke pelanggan setelah proses pengembalian barang selesai.
		"PENDING_REVIEW",    //+8.1 Pesanan sedang menunggu untuk ditinjau, mungkin karena ada masalah atau permintaan khusus dari pelanggan.
		"HOLD",              //+8.2 Pesanan ditahan, bisa jadi karena ada masalah dengan stok atau informasi pembayaran yang belum jelas.
		"SUCCESS",           //0 Pesanan telah selesai dengan sukses, semua proses dari pembayaran, pengiriman, hingga penerimaan telah berjalan lancar tanpa masalah.
	}[s]
}

//-> main collection
type Transaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// unique
	UserID        string `bson:"user_id"`
	InvoiceNumber string `bson:"name"`

	Products []TransactionProduct `bson:"products"`
	Status   string               `bson:"status"`
	Tracks   []TransactionTrack   `bson:"tracks"`

	CreatedAt  primitive.DateTime  `bson:"created_at"`
	SuccessAt  *primitive.DateTime `bson:"success_at,omitempty"`
	RejectedAt *primitive.DateTime `bson:"rejected_at,omitempty"`
}

type TransactionProduct struct {
	ProductID string `bson:"product_id"`

	// replace
	Name        string         `bson:"name"`
	Description string         `bson:"description"`
	Price       int            `bson:"price"`
	WeightGram  int            `bson:"weight_gram"`
	CategoryID  string         `bson:"category_id"`
	Images      []ProductImage `bson:"images"`

	UseQty int `bson:"use_qty"`
}

type TransactionTrack struct {
	Name   string `bson:"name"`
	Status string `bson:"status"`

	CreatedAt primitive.DateTime `bson:"created_at"`
}
