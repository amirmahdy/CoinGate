package utils

const (
	USDT = "USDT"
	BTC  = "BTC"
)

func IsSupportedCoin(coin string) bool {
	switch coin {
	case USDT, BTC:
		return true
	}
	return false
}
