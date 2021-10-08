package entity

import (
	"strings"
)

var ListValidBrETF = [5]string{"BOVA11", "SMAL11", "IVVB11", "HASH11", "ECOO11"}

func ConvertAssetLookup(symbol string, fullname string,
	symbolType string) SymbolLookup {

	fullnameTitle := strings.Title(strings.ToLower(fullname))

	for _, s := range strings.Fields(fullnameTitle) {
		if s == "Sa" || s == "Edp" || s == "Etf" || s == "Ftse" ||
			s == "Msci" || s == "Usa" {
			fullnameTitle = strings.ReplaceAll(fullnameTitle, s,
				strings.ToUpper(s))
		}
	}

	symbolLookup := SymbolLookup{
		Symbol:   strings.ReplaceAll(symbol, ".SA", ""),
		Fullname: fullnameTitle,
		Type:     symbolType,
	}

	return symbolLookup
}

func ConvertAssetPrice(symbol string, openPrice string, highPrice string,
	lowPrice string, currentPrice string, prevClosePrice string) SymbolPrice {

	symbolPrice := SymbolPrice{
		Symbol:         strings.ReplaceAll(symbol, ".SA", ""),
		OpenPrice:      StringToFloat64(openPrice),
		HighPrice:      StringToFloat64(highPrice),
		LowPrice:       StringToFloat64(lowPrice),
		CurrentPrice:   StringToFloat64(currentPrice),
		PrevClosePrice: StringToFloat64(prevClosePrice),
	}

	return symbolPrice
}

func ConvertUserInfo(email string, displayName string, userUid string) UserInfo {
	return UserInfo{
		Email:       email,
		DisplayName: displayName,
		UID:         userUid,
	}
}

func ConvertUserTokenInfo(idToken string, email string, emailVerified bool) UserTokenInfo {
	return UserTokenInfo{
		UserID:        idToken,
		Email:         email,
		EmailVerified: emailVerified,
	}
}
