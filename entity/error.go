package entity

import "errors"

var ErrInvalidAssetEntity = errors.New("Invalid Asset Entity: Blank Field")
var ErrInvalidCountryCode = errors.New("Invalid Country Value. It is only accepted BR and US")
var ErrInvalidAssetTypeName = errors.New("Invalid Asset Type Name. It is only accepted STOCK, ETF, FII and REIT")
