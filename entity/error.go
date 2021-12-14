package entity

import "errors"

// Token Maker
var (
	ErrInvalidTokenMethod  error = errors.New("INVALID_TOKEN_SIGNING_METHOD")
	ErrExpiredToken        error = errors.New("EXPIRED_TOKEN")
	ErrGenericInvalidToken error = errors.New("INVALID_TOKEN")
)

// General Erros
var (
	ErrInvalidCountryCode    error = errors.New("country: INVALID_COUNTRY_CODE")
	ErrInvalidCurrency       error = errors.New("currency: INVALID_CURRENCY_CODE")
	ErrInvalidBrazilCurrency error = errors.New("currency: must be BRL")
	ErrInvalidUsaCurrency    error = errors.New("currency: must be USD")
)

// Asset
var (
	ErrInvalidAssetEntityBlank         error = errors.New("asset: BLANK_FIELDS")
	ErrInvalidAssetPreferenceUndefined error = errors.New("asset: UNDEFINED_PREFERENCE")
	ErrInvalidAssetEntityValues        error = errors.New("asset: " + "INVALID_VALUES_COUNTRY_OR_TYPE")
	ErrInvalidAssetSymbol              error = errors.New("asset: SYMBOL_NOT_EXIST")
	ErrInvalidAssetSymbolExist         error = errors.New("asset: SYMBOL_ALREADY_EXIST.")
	ErrInvalidDeleteAsset              error = errors.New("deleteAsset: ASSET_NOT_EXIST")
)

// AssetType
var ErrInvalidAssetTypeName = errors.New("assetTypeName: INVALID_NAME")

// User
var (
	ErrInvalidUserNameBlank      error = errors.New("user: BLANK_USERNAME")
	ErrInvalidUserEmailBlank     error = errors.New("user: BLANK_EMAIL")
	ErrInvalidUserUidBlank       error = errors.New("user: BLANK_USER_UID")
	ErrInvalidUserTypeBlank      error = errors.New("user: BLANK_USER_TYPE")
	ErrInvalidUserToken          error = errors.New("user: INVALID_USER_TOKEN")
	ErrInvalidUserSendEmail      error = errors.New("user: EMAIL_NOT_SENT")
	ErrInvalidUserAdminPrivilege error = errors.New("user: WITHOUT_ADMIN_PERMISSION")
	ErrInvalidUserSearch         error = errors.New("searchUser: INVALID_UID")
)

// Order
var (
	ErrInvalidOrder               error = errors.New("orders: NO_ORDER_EXIST")
	ErrInvalidOrderType           error = errors.New("orders: INVALID_TYPE_VALUE")
	ErrInvalidOrderQuantityBrazil error = errors.New("orders: QUANTITY_MUST_BE_INTEGER")
	ErrInvalidOrderBuyQuantity    error = errors.New("orders: QUANTITY_MUST_BE_POSITIVE")
	ErrInvalidOrderSellQuantity   error = errors.New("orders: QUANTITY_MUST_BE_NEGATIVE")
	ErrInvalidOrderPrice          error = errors.New("orders: PRICE_MUST_BE_POSITIVE")
	ErrInvalidOrderOrderBy        error = errors.New("orders: INVALID_ORDER_BY_VALUE")
	ErrInvalidOrderLimit          error = errors.New("orders: LIMIT_MUST_BE_INTEGER")
	ErrInvalidOrderOffset         error = errors.New("orders: OFFSET_MUST_BE_INTEGER")
)

// Earning
var (
	ErrInvalidEarningsAmount            error = errors.New("earnings: AMOUNT_MUST_BE_POSITIVE")
	ErrInvalidEarningType               error = errors.New("earnings: INVALID_TYPE_VALUE")
	ErrInvalidEarningsCreateBlankFields error = errors.New("earnings: MISSING_FIELDS")
	ErrInvalidEarningsOrderBy           error = errors.New("earnings: INVALID_ORDER_BY_VALUE")
	ErrInvalidEarningsLimit             error = errors.New("earnings: LIMIT_MUST_BE_INTEGER")
	ErrInvalidEarningsOffset            error = errors.New("earnings: OFFSET_MUST_BE_INTEGER")
)

// Brokerage
var (
	ErrInvalidBrokerageSearchType      error = errors.New("brokerage: INVALID_SEARCH_TYPE")
	ErrInvalidBrokerageNameSearch      error = errors.New("brokerage: INVALID_NAME")
	ErrInvalidBrokerageNameSearchBlank error = errors.New("brokerage: BLANK_NAME")
)

// Sector
var ErrInvalidSectorSearchName = errors.New("sector: NAME_NOT_EXIST")

// AssetUser
var ErrInvalidAssetUser = errors.New("assetUser: RELATION_NOT_EXIST")
var ErrinvalidAssetUserAlreadyExists = errors.New("assetUser: RELATION_EXIST")

// Database Errors
var ErrInvalidAssetType = errors.New("SearchAssetsPerAssetType: There is no asset for this type in this country")

// API Errors: Query
var (
	ErrInvalidApiQuerySymbolBlank       error = errors.New("query: BLANK_SYMBOL_VALUE")
	ErrInvalidApiQueryTypeBlank         error = errors.New("query: BLANK_TYPE_VALUE")
	ErrInvalidApiQueryCountryBlank      error = errors.New("query: BLANK_COUNTRY_VALUE")
	ErrInvalidApiQueryWithOrderResume   error = errors.New("query: INVALID_WITH_ORDER_RESUME_VALUE")
	ErrInvalidApiQueryWithPrice         error = errors.New("query: INVALID_WITH_PRICE_VALUE")
	ErrInvalidApiQueryWithOrders        error = errors.New("query: INVALID_WITH_ORDERS_VALUE")
	ErrInvalidApiQueryMyUser            error = errors.New("query: INVALID_MY_USER_VALUE")
	ErrInvalidApiQueryLoginType         error = errors.New("query: INVALID_LOGIN_TYPE_VALUE")
	ErrInvalidApiQueryOAuth2Code        error = errors.New("query: INVALID_OAUTH2_CODE_VALUE")
	ErrInvalidApiQueryOAuth2CodeBlank   error = errors.New("query: MISSING_CODE_VALUE")
	ErrInvalidApiQueryStateDoesNotMatch error = errors.New("query: STATE_DOES_NOT_MATCH")
	ErrInvalidApiQueryStateBlank        error = errors.New("query: STATE_MISSING_VALUE")
	ErrInvalidApiQueryState             error = errors.New("query: STATE_")
	// ErrInvalidApiQueryInvalidToken      error = errors.New("query: STATE_INVALID_TOKEN")
)

// API Errors: Parameters
var ErrInvalidApiParamsCompany = errors.New("params: INVALID_COMPANY_VALUE")

// API Errors: JSON
var (
	ErrInvalidApiBody        error = errors.New("httpBody: WRONG_JSON")
	ErrInvalidApiOrderUpdate error = errors.New("updateOrder: MISSING_JSON_KEYS")
)

// API Error: General Messages
var (
	ErrMessageApiInternalError    error = errors.New("Internal Server Error. Please contact us to correct this error")
	ErrMessageApiAuthentication   error = errors.New("Authentication information is missing or invalid")
	ErrMessageApiAuthorization    error = errors.New("The user is not authorized to execute this request")
	ErrMessageApiRequest          error = errors.New("Invalid request. Please see our API documentation.")
	ErrMessageApiEarningAssetUser error = errors.New("This user does not have any registered earning for the requested Asset")
	ErrMessageApiAssetSymbolUser  error = errors.New("This symbol/asset does not exist in our database or in your asset table")
	ErrMessageApiOrderId          error = errors.New("The authenticated user does not have this order with the requested ID")
	ErrMessageApiEarningId        error = errors.New("The authenticated user does not have this earning with the requested ID")
	ErrMessageApiSectorName       error = errors.New("The database does not have this sector")
	ErrMessageApiEmail            error = errors.New("The email for password reset was not found")
)
