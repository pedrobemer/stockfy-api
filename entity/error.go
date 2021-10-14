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
var ErrInvalidOrdersFromAssetUser = errors.New("There is no orders for this asset considering your user")
var ErrInvalidOrderId = errors.New("There is no order with this ID for your user")

// Application Logic Errors: Brokerage
var ErrInvalidBrokerageSearchType = errors.New("This search type for brokerage firms is not valid")
var ErrInvalidBrokerageNameSearch = errors.New("Invalid brokerage name to search")

// Application Logic Errors: Earning
var ErrInvalidEarningId = errors.New("There is no earning with this ID for your user")

// Application Logic Errors: Sector
var ErrInvalidSectorSearchName = errors.New("Impossible to find this sector")

// Database Errors
var ErrInvalidSectorName = errors.New("CreateSector: Impossible to create a blank sector")
var ErrInvalidSearchAssetName = errors.New("SearchAsset: There is no Asset in our database with this symbol")
var ErrInvalidSearchUser = errors.New("SearchUser: There is no user in our database with this UID")
var ErrInvalidAssetType = errors.New("SearchAssetsPerAssetType: There is no asset for this type in this country")
var ErrInvalidAssetUser = errors.New("assetUser: RELATION_NOT_EXIST")
var ErrInvalidDeleteAsset = errors.New("deleteAsset: ASSET_NOT_EXIST")

// API Errors
var ErrInvalidApiAuthentication = errors.New("Authentication information is missing or invalid")
var ErrInvalidApiAuthorization = errors.New("The user is not authorized to execute this request")
var ErrInvalidApiUserAdminPrivilege = errors.New("user: WITHOUT_ADMIN_PERMISSION")
var ErrInvalidApiRequest = errors.New("Invalid request. Please see our API documentation.")
var ErrInvalidApiEmail = errors.New("The email for password reset was not found")
var ErrInvalidApiAssetSymbolUser = errors.New("This symbol/asset does not exist in our database or in your asset table.")
var ErrInvalidApiAssetSymbolExist = errors.New("asset: SYMBOL_EXIST.")
var ErrInvalidApiAssetSymbol = errors.New("asset: INVALID_SYMBOL")
var ErrInvalidApiQueryTypeBlank = errors.New("type: BLANK_VALUE")
var ErrInvalidApiQueryCountryBlank = errors.New("country: BLANK_VALUE")
var ErrInvalidApiQueryWithOrderResume = errors.New("withOrderResume: INVALID_VALUE")
var ErrInvalidApiQueryWithOrders = errors.New("withOrders: INVALID_VALUE")
var ErrInvalidApiQueryMyUser = errors.New("myUser: INVALID_VALUE")
var ErrInvalidApiBody = errors.New("Wrong body request")

// var ErrInvalidApiRequest = errors.New("Wrong REST API. Please see our documentation to properly execute requests for our API.")
var ErrInvalidApiMissedKeysBody = errors.New("Missed JSON keys in the body request. Please read our REST API documentation.")
var ErrInvalidApiOrderUpdate = errors.New("Missed JSON keys in the body request. Please read our REST API documentation.")
var ErrInvalidApiOrderType = errors.New("Wrong value for the order type field in the order body. Please see our documentation to properly execute requests for our API")
var ErrInvalidApiBrazilOrderQuantity = errors.New("Quantity value must have a integer value")
var ErrInvalidApiBrazilOrderCurrency = errors.New("Currency does not match for Brazil investment")
var ErrInvalidApiUsaOrderCurrency = errors.New("Currency does not match for USA investment")
var ErrInvalidApiOrderBuyQuantity = errors.New("Buy Order must have a positive quantity")
var ErrInvalidApiOrderSellQuantity = errors.New("Buy Order must have a negative quantity")
var ErrInvalidApiOrderPrice = errors.New("Order price field must have a positive or zero value")

// var ErrInvalidApiAssetSymbol = errors.New("Wrong value for the Asset symbol. Please see our documentation to properly execute requests for our API")
var ErrInvalidApiEarningsAmount = errors.New("Wrong value for the field amount in the body request. Please read our REST API documentation")
var ErrInvalidApiEarningType = errors.New("Wrong value for the earning type in the body request. Please read our REST API documentation")
var ErrInvalidApiEarningSymbol = errors.New("This user does not have this asset to register a earning")
var ErrInvalidApiEarningAssetUser = errors.New("This user does not have any registered earning for the requested Asset")
