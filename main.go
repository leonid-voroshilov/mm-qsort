package main

import (
	"fmt"
	"time"

	comparator "github.com/leonid-voroshilov/mm-qsort/pkg/Comparator"
	"github.com/leonid-voroshilov/mm-qsort/pkg/qsort"
)

func main() {
	n := 100000
	// Тестирование с целыми числами
	fmt.Printf("Тестирование с целыми числами: \n N = %d\n", n)
	intData := qsort.GenerateRandomInts(n)

	// Копируем данные для сравнения
	intDataCopy := make([]int, len(intData))
	copy(intDataCopy, intData)

	intComp := comparator.IntC{}

	// Измеряем время параллельной сортировки
	start := time.Now()
	qsort.ParallelQuickSort(intData, intComp)
	parallelTime := time.Since(start)

	// Измеряем время последовательной сортировки
	start = time.Now()
	qsort.SequentialQuickSort(intDataCopy, intComp)
	sequentialTime := time.Since(start)

	fmt.Printf("	Параллельная сортировка: %v\n", parallelTime)
	fmt.Printf("	Последовательная сортировка: %v\n", sequentialTime)
	fmt.Printf("	Ускорение: %.2fx\n", float64(sequentialTime)/float64(parallelTime))

	// Проверяем корректность сортировки
	fmt.Printf("Массив отсортирован корректно: %v\n", isSorted(intData, intComp))

	// Тестирование со строками
	fmt.Println("Тестирование со строками:")
	stringData := []string{"zebra", "apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
	stringComp := comparator.StringC{}

	fmt.Println("	До сортировки:", stringData)
	qsort.ParallelQuickSort(stringData, stringComp)
	fmt.Println("	После сортировки:", stringData)
}

// Функция для проверки корректности сортировки
func isSorted[T any](data []T, comp comparator.Comparator[T]) bool {
	for i := 1; i < len(data); i++ {
		if comp.Compare(data[i-1], data[i]) > 0 {
			return false
		}
	}
	return true
}
