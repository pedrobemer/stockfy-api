package entity

import "errors"

// Application Logic Errors
var ErrInvalidAssetEntity = errors.New("Invalid Asset Entity: Blank Field")
var ErrInvalidCountryCode = errors.New("Invalid Country Value. It is only accepted BR and US")
var ErrInvalidAssetTypeName = errors.New("Invalid Asset Type Name. It is only accepted STOCK, ETF, FII and REIT")
var ErrInvalidAssetSymbol = errors.New("Invalid Asset Symbol. We could not find the specified symbol")
var ErrInvalidUserName = errors.New("Invalid username. Blank field")
var ErrInvalidUserEmail = errors.New("Invalid User email. Blank field")
var ErrInvalidUserUid = errors.New("Invalid User UID. Blank field")
var ErrInvalidUserType = errors.New("Invalid user Type. Blank field")
var ErrInvalidUserToken = errors.New("Invalid User information to get the valid ID token")
var ErrInvalidUserEmailVerification = errors.New("Problems to send the email for user verification")

// Database Errors
var ErrInvalidSectorName = errors.New("CreateSector: Impossible to create a blank sector")
