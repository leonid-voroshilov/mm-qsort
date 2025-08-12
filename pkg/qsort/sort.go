package qsort

import (
	"math/rand"
	"runtime"
	"sync"

	comparator "github.com/leonid-voroshilov/mm-qsort/pkg/Comparator"
)

func GenerateRandomInts(size int) []int {
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(10000)
	}
	return data
}

// ParallelQuickSort - параллельная быстрая сортировка
func ParallelQuickSort[T any](data []T, comp comparator.Comparator[T]) {
	if len(data) <= 1 {
		return
	}

	maxGoroutines := runtime.NumCPU() // для оптимальности, максимум горутин берём как количество ядер, чтобы не было простоя
	parallelQuickSort(data, comp, maxGoroutines)
}

func parallelQuickSort[T any](data []T, comp comparator.Comparator[T], maxGoroutines int) {
	if len(data) <= 1 {
		return
	}

	// Для небольших массивов используем последовательную сортировку
	if len(data) < 1000 || maxGoroutines <= 1 {
		SequentialQuickSort(data, comp)
		return
	}

	pivotIndex := partition(data, comp)

	var wg sync.WaitGroup

	// Распределяем горутины между левой и правой частями
	leftGoroutines := maxGoroutines / 2
	rightGoroutines := maxGoroutines - leftGoroutines

	// Сортируем левую часть в отдельной горутине
	if pivotIndex > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelQuickSort(data[:pivotIndex], comp, leftGoroutines)
		}()
	}

	// Сортируем правую часть в отдельной горутине
	if pivotIndex < len(data)-1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelQuickSort(data[pivotIndex+1:], comp, rightGoroutines)
		}()
	}

	wg.Wait()
}

// SequentialQuickSort — последовательная быстрая сортировка для небольших массивов
func SequentialQuickSort[T any](data []T, comp comparator.Comparator[T]) {
	if len(data) <= 1 {
		return
	}

	pivotIndex := partition(data, comp)

	SequentialQuickSort(data[:pivotIndex], comp)
	SequentialQuickSort(data[pivotIndex+1:], comp)
}

// partition разбивает массив относительно опорного элемента
// Возвращает индекс опорного элемента после разбиения
func partition[T any](data []T, comp comparator.Comparator[T]) int {
	if len(data) <= 1 {
		return 0
	}

	// Используем медиану из трех для выбора опорного элемента
	pivotIndex := medianOfThree(data, comp)

	// Помещаем опорный элемент в конец массива
	lastIndex := len(data) - 1
	data[pivotIndex], data[lastIndex] = data[lastIndex], data[pivotIndex]
	pivot := data[lastIndex]

	// Индекс для размещения следующего элемента, меньшего опорного
	storeIndex := 0

	// Проходим по всем элементам кроме последнего (опорного)
	for i := 0; i < lastIndex; i++ {
		if comp.Compare(data[i], pivot) <= 0 {
			data[i], data[storeIndex] = data[storeIndex], data[i]
			storeIndex++
		}
	}

	// Помещаем опорный элемент на правильную позицию
	data[storeIndex], data[lastIndex] = data[lastIndex], data[storeIndex]

	return storeIndex
}

// medianOfThree выбирает медиану из первого, среднего и последнего элементов
// для лучшего выбора опорного элемента
func medianOfThree[T any](data []T, comp comparator.Comparator[T]) int {
	length := len(data)
	if length < 3 {
		return 0
	}

	first, middle, last := 0, length/2, length-1

	// Сортируем индексы по значениям элементов
	if comp.Compare(data[first], data[middle]) > 0 {
		first, middle = middle, first
	}
	if comp.Compare(data[middle], data[last]) > 0 {
		middle, last = last, middle
		if comp.Compare(data[first], data[middle]) > 0 {
			first, middle = middle, first
		}
	}

	return middle
}

// Альтернативная версия с настраиваемым порогом параллелизма
func ParallelQuickSortWithThreshold[T any](data []T, comp comparator.Comparator[T], threshold int) {
	if len(data) <= 1 {
		return
	}

	maxGoroutines := runtime.NumCPU()
	parallelQuickSortWithThreshold(data, comp, maxGoroutines, threshold)
}

func parallelQuickSortWithThreshold[T any](data []T, comp comparator.Comparator[T], maxGoroutines, threshold int) {
	if len(data) <= 1 {
		return
	}

	if len(data) < threshold || maxGoroutines <= 1 {
		SequentialQuickSort(data, comp)
		return
	}

	pivotIndex := partition(data, comp)

	var wg sync.WaitGroup
	leftGoroutines := maxGoroutines / 2
	rightGoroutines := maxGoroutines - leftGoroutines

	if pivotIndex > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelQuickSortWithThreshold(data[:pivotIndex], comp, leftGoroutines, threshold)
		}()
	}

	if pivotIndex < len(data)-1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parallelQuickSortWithThreshold(data[pivotIndex+1:], comp, rightGoroutines, threshold)
		}()
	}

	wg.Wait()
}
