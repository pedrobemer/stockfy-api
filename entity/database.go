package entity

import "time"

type OrderInfos struct {
	TotalQuantity        float64 `json:"totalQuantity,omitempty"`
	WeightedAdjPrice     float64 `json:"weightedAdjPrice,omitempty"`
	WeightedAveragePrice float64 `json:"weightedAveragePrice,omitempty"`
}

type Sector struct {
	Id        string    `db:"id" json:",omitempty"`
	Name      string    `db:"name" json:",omitempty"`
	CreatedAt time.Time `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:",omitempty"`
}

type Brokerage struct {
	Id        string    `db:"id" json:",omitempty"`
	Name      string    `db:"name" json:",omitempty"`
	Fullname  string    `db:"fullname" json:",omitempty"`
	Country   string    `db:"country" json:",omitempty"`
	CreatedAt time.Time `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:",omitempty"`
}

type AssetType struct {
	Id        string    `db:"id" json:",omitempty"`
	Type      string    `db:"type" json:",omitempty"`
	Name      string    `db:"name" json:",omitempty"`
	Country   string    `db:"country" json:",omitempty"`
	CreatedAt time.Time `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:",omitempty"`
	Assets    []Asset   `db:"assets" json:",omitempty"`
}

type Asset struct {
	Id         string      `db:"id"`
	Preference *string     `db:"preference"`
	Fullname   string      `db:"fullname"`
	Symbol     string      `db:"symbol"`
	Sector     *Sector     `db:"sector" json:",omitempty"`
	AssetType  *AssetType  `db:"asset_type" json:",omitempty"`
	CreatedAt  time.Time   `db:"created_at" json:",omitempty"`
	UpdatedAt  time.Time   `db:"updated_at" json:",omitempty"`
	OrderInfo  *OrderInfos `db:"orders_info" json:",omitempty"`
	OrdersList []Order     `db:"orders_list" json:",omitempty"`
}

type Order struct {
	Id        string     `db:"id" json:",omitempty"`
	Quantity  float64    `db:"quantity" json:",omitempty"`
	Price     float64    `db:"price" json:",omitempty"`
	Currency  string     `db:"currency" json:",omitempty"`
	OrderType string     `db:"order_type" json:",omitempty"`
	Date      time.Time  `db:"date" json:",omitempty"`
	Brokerage *Brokerage `db:"brokerage" json:",omitempty"`
	Asset     *Asset     `db:"asset" json:",omitempty"`
	UserUid   string     `db:"user_uid" json:",omitempty"`
	CreatedAt time.Time  `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time  `db:"updated_at" json:",omitempty"`
}

type Earnings struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Earning   float64   `json:"earning"`
	Currency  string    `json:"currency"`
	Date      time.Time `json:"date"`
	Asset     *Asset    `db:"asset" json:",omitempty"`
	UserUid   string    `db:"user_uid" json:",omitempty"`
	CreatedAt time.Time `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:",omitempty"`
}

type AssetUsers struct {
	AssetId string `db:"asset_id"`
	UserUid string `db:"user_uid"`
}

type Users struct {
	Id        string    `db:"id" json:"id,omitempty"`
	Uid       string    `db:"uid" json:"uid"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:",omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:",omitempty"`
}

type SymbolLookup struct {
	Fullname string `json:",omitempty"`
	Symbol   string `json:",omitempty"`
	Type     string `json:",omitempty"`
}

type SymbolPrice struct {
	Symbol         string  `json:",omitempty"`
	CurrentPrice   float64 `json:",omitempty"`
	HighPrice      float64 `json:",omitempty"`
	LowPrice       float64 `json:",omitempty"`
	OpenPrice      float64 `json:",omitempty"`
	PrevClosePrice float64 `json:",omitempty"`
	MarketCap      float64 `json:",omitempty"`
}

type UserInfo struct {
	DisplayName string
	Email       string
	UID         string
}

type ReqIdToken struct {
	Token              string `json:"token,omitempty"`
	RequestSecureToken bool   `json:"requestSecureToken,omitempty"`
	Kind               string `json:"kind,omitempty"`
	IdToken            string `json:"idToken,omitempty"`
	IsNewUser          bool   `json:"isNewUser,omitempty"`
}

type EmailVerificationResponse struct {
	UserIdToken string                 `json:",omitempty"`
	Email       string                 `json:"email,omitempty"`
	Error       map[string]interface{} `json:"error,omitempty"`
}

type EmailForgotPasswordResponse struct {
	Email string                 `json:"email"`
	Error map[string]interface{} `json:"error,omitempty"`
}
