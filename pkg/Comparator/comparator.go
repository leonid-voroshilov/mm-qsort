package comparator

// Компаратор - интерфейс для сравнения элементов
type Comparator[T any] interface {
	Compare(a, b T) int // -1 если a < b, 0 если a == b, 1 если a > b
}

// Компаратор для целых чисел
type IntC struct{}

func (c IntC) Compare(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// Компаратор  для строк
type StringC struct{}

func (c StringC) Compare(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
