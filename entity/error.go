package entity

import "errors"

// Application Logic Errors: General Erros
var ErrInvalidCountryCode = errors.New("country: INVALID_COUNTRY_CODE")
var ErrInvalidCurrency = errors.New("Invalid Currency Value. It is only accepted USD and BRL")
var ErrInvalidBrazilCurrency = errors.New("BRL currency does not match the Country code")
var ErrInvalidUnitedStatesCurrency = errors.New("USD currency does not match the Country code")

// Application Logic Errors: Asset
var ErrInvalidAssetEntity = errors.New("Invalid Asset Entity: Blank Field")
var ErrInvalidAssetSymbol = errors.New("asset: SYMBOL_NOT_EXIST")

// Application Logic Errors: AssetType
var ErrInvalidAssetTypeName = errors.New("assetTypeName: INVALID_NAME")

// Application Logic Errors: User
var ErrInvalidUserName = errors.New("Invalid username. Blank field")
var ErrInvalidUserEmail = errors.New("Invalid User email. Blank field")
var ErrInvalidUserUid = errors.New("Invalid User UID. Blank field")
var ErrInvalidUserType = errors.New("Invalid user Type. Blank field")
var ErrInvalidUserToken = errors.New("Invalid User information to get the valid ID token")
var ErrInvalidUserEmailVerification = errors.New("Problems to send the email for user verification")
var ErrInvalidUserEmailForgotPassword = errors.New("Problems to send the email to update the password")

// Application Logic Errors: Order
var ErrInvalidOrdersFromAssetUser = errors.New("orders: NO_ORDER_EXIST")
var ErrInvalidOrderId = errors.New("There is no order with this ID for your user")

// Application Logic Errors: Brokerage
var ErrInvalidBrokerageSearchType = errors.New("brokerage: INVALID_SEARCH_TYPE")
var ErrInvalidBrokerageNameSearch = errors.New("brokerage: INVALID_NAME")
var ErrInvalidBrokerageNameSearchBlank = errors.New("brokerage: BLANK_NAME")

// Application Logic Errors: Earning
var ErrInvalidEarningId = errors.New("There is no earning with this ID for your user")

// Application Logic Errors: Sector
var ErrInvalidSectorSearchName = errors.New("sector: NAME_NOT_EXIST")

// Database Errors
var ErrInvalidSectorName = errors.New("CreateSector: Impossible to create a blank sector")
var ErrInvalidSearchAssetName = errors.New("SearchAsset: There is no Asset in our database with this symbol")
var ErrInvalidSearchUser = errors.New("SearchUser: There is no user in our database with this UID")
var ErrInvalidAssetType = errors.New("SearchAssetsPerAssetType: There is no asset for this type in this country")
var ErrInvalidAssetUser = errors.New("assetUser: RELATION_NOT_EXIST")
var ErrInvalidDeleteAsset = errors.New("deleteAsset: ASSET_NOT_EXIST")

// API Errors
var ErrInvalidApiInternalError = errors.New("Internal Server Error. Please contact us to correct this error")
var ErrInvalidApiAuthentication = errors.New("Authentication information is missing or invalid")
var ErrInvalidApiAuthorization = errors.New("The user is not authorized to execute this request")
var ErrInvalidApiRequest = errors.New("Invalid request. Please see our API documentation.")
var ErrInvalidApiEmail = errors.New("The email for password reset was not found")
var ErrInvalidApiAssetSymbolUser = errors.New("This symbol/asset does not exist in our database or in your asset table")
var ErrInvalidApiSectorName = errors.New("The database does not have this sector")
var ErrInvalidApiOrderId = errors.New("The authenticated user does not have this order with the requested ID")
var ErrInvalidApiEarningId = errors.New("The authenticated user does not have this earning with the requested ID")
var ErrInvalidApiUserAdminPrivilege = errors.New("user: WITHOUT_ADMIN_PERMISSION")
var ErrInvalidApiAssetSymbolExist = errors.New("asset: SYMBOL_EXIST.")
var ErrInvalidApiAssetSymbol = errors.New("asset: INVALID_SYMBOL")
var ErrInvalidApiQuerySymbolBlank = errors.New("symbol: BLANK_VALUE")
var ErrInvalidApiQueryTypeBlank = errors.New("type: BLANK_VALUE")
var ErrInvalidApiQueryCountryBlank = errors.New("country: BLANK_VALUE")
var ErrInvalidApiQueryWithOrderResume = errors.New("withOrderResume: INVALID_VALUE")
var ErrInvalidApiQueryWithOrders = errors.New("withOrders: INVALID_VALUE")
var ErrInvalidApiQueryMyUser = errors.New("myUser: INVALID_VALUE")
var ErrInvalidApiBody = errors.New("Wrong body request")

// var ErrInvalidApiRequest = errors.New("Wrong REST API. Please see our documentation to properly execute requests for our API.")
var ErrInvalidApiEarningsCreate = errors.New("createEarning: MISSING_JSON_KEYS")
var ErrInvalidApiOrderUpdate = errors.New("updateOrder: MISSING_JSON_KEYS")
var ErrInvalidApiOrderType = errors.New("orderType: INVALID_VALUE")
var ErrInvalidApiBrazilOrderQuantity = errors.New("quantity: must be integer")
var ErrInvalidApiBrazilOrderCurrency = errors.New("currency: must be BRL")
var ErrInvalidApiUsaOrderCurrency = errors.New("currency: must be USD")
var ErrInvalidApiOrderBuyQuantity = errors.New("quantity: must be positive")
var ErrInvalidApiOrderSellQuantity = errors.New("quantity: must be negative")
var ErrInvalidApiOrderPrice = errors.New("price: must be positive or zero")

// var ErrInvalidApiAssetSymbol = errors.New("Wrong value for the Asset symbol. Please see our documentation to properly execute requests for our API")
var ErrInvalidApiEarningsAmount = errors.New("amount: must have a positive value")
var ErrInvalidApiEarningType = errors.New("earningType: INVALID_TYPE")
var ErrInvalidApiEarningSymbol = errors.New("This user does not have this asset to register a earning")
var ErrInvalidApiEarningAssetUser = errors.New("This user does not have any registered earning for the requested Asset")
