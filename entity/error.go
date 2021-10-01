package entity

import "errors"

// Application Logic Errors
var ErrInvalidAssetEntity = errors.New("Invalid Asset Entity: Blank Field")
var ErrInvalidCountryCode = errors.New("Invalid Country Value. It is only accepted BR and US")
var ErrInvalidAssetTypeName = errors.New("Invalid Asset Type Name. It is only accepted STOCK, ETF, FII and REIT")
var ErrInvalidAssetSymbol = errors.New("Invalid Asset Symbol. We could not find the specified symbol")

// Database Errors
var ErrInvalidSectorName = errors.New("CreateSector: Impossible to create a blank sector")
