package entity

import "errors"

// General Erros
var ErrInvalidCountryCode = errors.New("country: INVALID_COUNTRY_CODE")
var ErrInvalidCurrency = errors.New("currency: INVALID_CURRENCY_CODE")
var ErrInvalidBrazilCurrency = errors.New("currency: must be BRL")
var ErrInvalidUsaCurrency = errors.New("currency: must be USD")

// Asset
var ErrInvalidAssetEntityBlank = errors.New("asset: BLANK_FIELDS")
var ErrInvalidAssetPreferenceUndefined = errors.New("asset: UNDEFINED_PREFERENCE")
var ErrInvalidAssetEntityValues = errors.New("asset: INVALID_VALUES_COUNTRY_OR_TYPE")
var ErrInvalidAssetSymbol = errors.New("asset: SYMBOL_NOT_EXIST")
var ErrInvalidAssetSymbolExist = errors.New("asset: SYMBOL_ALREADY_EXIST.")
var ErrInvalidDeleteAsset = errors.New("deleteAsset: ASSET_NOT_EXIST")

// AssetType
var ErrInvalidAssetTypeName = errors.New("assetTypeName: INVALID_NAME")

// User
var ErrInvalidUserNameBlank = errors.New("user: BLANK_USERNAME")
var ErrInvalidUserEmailBlank = errors.New("user: BLANK_EMAIL")
var ErrInvalidUserUidBlank = errors.New("user: BLANK_USER_UID")
var ErrInvalidUserTypeBlank = errors.New("user: BLANK_USER_TYPE")
var ErrInvalidUserToken = errors.New("user: INVALID_USER_TOKEN")
var ErrInvalidUserSendEmail = errors.New("user: EMAIL_NOT_SENT")
var ErrInvalidUserAdminPrivilege = errors.New("user: WITHOUT_ADMIN_PERMISSION")
var ErrInvalidUserSearch = errors.New("searchUser: INVALID_UID")

// Order
var ErrInvalidOrder = errors.New("orders: NO_ORDER_EXIST")
var ErrInvalidOrderType = errors.New("orders: INVALID_TYPE_VALUE")
var ErrInvalidOrderQuantityBrazil = errors.New("orders: QUANTITY_MUST_BE_INTEGER")
var ErrInvalidOrderBuyQuantity = errors.New("orders: QUANTITY_MUST_BE_POSITIVE")
var ErrInvalidOrderSellQuantity = errors.New("orders: QUANTITY_MUST_BE_NEGATIVE")
var ErrInvalidOrderPrice = errors.New("orders: PRICE_MUST_BE_POSITIVE")

// Earning
var ErrInvalidEarningsAmount = errors.New("earnings: AMOUNT_MUST_BE_POSITIVE")
var ErrInvalidEarningType = errors.New("earnings: INVALID_TYPE_VALUE")
var ErrInvalidEarningsCreateBlankFields = errors.New("createEarning: MISSING_FIELDS")

// Brokerage
var ErrInvalidBrokerageSearchType = errors.New("brokerage: INVALID_SEARCH_TYPE")
var ErrInvalidBrokerageNameSearch = errors.New("brokerage: INVALID_NAME")
var ErrInvalidBrokerageNameSearchBlank = errors.New("brokerage: BLANK_NAME")

// Sector
var ErrInvalidSectorSearchName = errors.New("sector: NAME_NOT_EXIST")

// AssetUser
var ErrInvalidAssetUser = errors.New("assetUser: RELATION_NOT_EXIST")

// Database Errors
var ErrInvalidAssetType = errors.New("SearchAssetsPerAssetType: There is no asset for this type in this country")

// API Errors: Query
var ErrInvalidApiQuerySymbolBlank = errors.New("query: BLANK_SYMBOL_VALUE")
var ErrInvalidApiQueryTypeBlank = errors.New("query: BLANK_TYPE_VALUE")
var ErrInvalidApiQueryCountryBlank = errors.New("query: BLANK_COUNTRY_VALUE")
var ErrInvalidApiQueryWithOrderResume = errors.New("query: INVALID_WITH_ORDER_RESUME_VALUE")
var ErrInvalidApiQueryWithOrders = errors.New("query: INVALID_WITH_ORDERS_VALUE")
var ErrInvalidApiQueryMyUser = errors.New("query: INVALID_MY_USER_VALUE")

// API Errors: JSON
var ErrInvalidApiBody = errors.New("httpBody: WRONG_JSON")
var ErrInvalidApiOrderUpdate = errors.New("updateOrder: MISSING_JSON_KEYS")

// API Error: General Messages
var ErrMessageApiInternalError = errors.New("Internal Server Error. Please contact us to correct this error")
var ErrMessageApiAuthentication = errors.New("Authentication information is missing or invalid")
var ErrMessageApiAuthorization = errors.New("The user is not authorized to execute this request")
var ErrMessageApiRequest = errors.New("Invalid request. Please see our API documentation.")
var ErrMessageApiEarningAssetUser = errors.New("This user does not have any registered earning for the requested Asset")
var ErrMessageApiAssetSymbolUser = errors.New("This symbol/asset does not exist in our database or in your asset table")
var ErrMessageApiOrderId = errors.New("The authenticated user does not have this order with the requested ID")
var ErrMessageApiEarningId = errors.New("The authenticated user does not have this earning with the requested ID")
var ErrMessageApiSectorName = errors.New("The database does not have this sector")
var ErrMessageApiEmail = errors.New("The email for password reset was not found")
