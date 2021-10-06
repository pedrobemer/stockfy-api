package entity

import "errors"

// Application Logic Errors: General Erros
var ErrInvalidCountryCode = errors.New("Invalid Country Value. It is only accepted BR and US")

// Application Logic Errors: Asset
var ErrInvalidAssetEntity = errors.New("Invalid Asset Entity: Blank Field")
var ErrInvalidAssetSymbol = errors.New("Invalid Asset Symbol. We could not find the specified symbol")

// Application Logic Errors: AssetType
var ErrInvalidAssetTypeName = errors.New("Invalid Asset Type Name. It is only accepted STOCK, ETF, FII and REIT")

// Application Logic Errors: User
var ErrInvalidUserName = errors.New("Invalid username. Blank field")
var ErrInvalidUserEmail = errors.New("Invalid User email. Blank field")
var ErrInvalidUserUid = errors.New("Invalid User UID. Blank field")
var ErrInvalidUserType = errors.New("Invalid user Type. Blank field")
var ErrInvalidUserToken = errors.New("Invalid User information to get the valid ID token")
var ErrInvalidUserEmailVerification = errors.New("Problems to send the email for user verification")
var ErrInvalidUserEmailForgotPassword = errors.New("Problems to send the email to update the password")

// Application Logic Errors: Brokerage
var ErrInvalidBrokerageSearchType = errors.New("This searcy type for brokerage firms is not valid")

// Database Errors
var ErrInvalidSectorName = errors.New("CreateSector: Impossible to create a blank sector")
var ErrInvalidSearchAssetName = errors.New("SearchAsset: There is no Asset in our database with this symbol")
var ErrInvalidSearchUser = errors.New("SearchUser: There is no user in our database with this UID")
var ErrInvalidAssetType = errors.New("SearchAssetsPerAssetType: There is no asset for this type in this country")
var ErrInvalidAssetUser = errors.New("AssetUser: This asset is not registered for this user")
var ErrInvalidDeleteAsset = errors.New("DeleteAsset: This asset does not exist")

// API Errors
var ErrInvalidApiRequest = errors.New("Wrong REST API. Please see our documentation to properly execute requests for our API.")
var ErrInvalidApiAuthorization = errors.New("This user is not authorized to execute this request")
var ErrInvalidApiOrderType = errors.New("Wrong value for the order type field in the order body. Please see our documentation to properly execute requests for our API")
var ErrInvalidApiBrazilOrderQuantity = errors.New("Quantity value must have a integer value")
var ErrInvalidApiBrazilOrderCurrency = errors.New("Currency does not match for Brazil investment")
var ErrInvalidApiUsaOrderCurrency = errors.New("Currency does not match for USA investment")
var ErrInvalidApiOrderBuyQuantity = errors.New("Buy Order must have a positive quantity")
var ErrInvalidApiOrderSellQuantity = errors.New("Buy Order must have a negative quantity")
var ErrInvalidApiOrderPrice = errors.New("Order price field must have a positive or zero value")
