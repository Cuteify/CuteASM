package utils

func GetLength(name string) int {
	switch name {
	case "BB":
		return 1
	case "WW":
		return 2
	case "DW":
		return 4
	case "QW":
		return 8
	case "TW":
		return 10
	case "OW":
		return 16
	case "YW":
		return 32
	case "ZW":
		return 64
	}
	return 0
}
