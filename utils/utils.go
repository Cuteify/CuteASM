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
	}
	return 0
}
