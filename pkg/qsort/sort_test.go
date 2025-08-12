package qsort

import (
	"reflect"
	"sort"
	"testing"

	comparator "github.com/leonid-voroshilov/mm-qsort/pkg/Comparator"
)

// Тестовые компараторы
type IntComparator struct{}

func (c IntComparator) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

type StringComparator struct{}

func (c StringComparator) Compare(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// Вспомогательные функции
func isSorted[T any](data []T, comp comparator.Comparator[T]) bool {
	for i := 1; i < len(data); i++ {
		if comp.Compare(data[i-1], data[i]) > 0 {
			return false
		}
	}
	return true
}

func generateSortedInts(size int) []int {
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}
	return data
}

func generateReversedInts(size int) []int {
	data := make([]int, size)
	for i := range data {
		data[i] = size - i
	}
	return data
}

func copySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// Тесты для ParallelQuickSort
func TestParallelQuickSort(t *testing.T) {
	comp := IntComparator{}

	tests := []struct {
		name string
		data []int
	}{
		{"Empty slice", []int{}},
		{"Single element", []int{42}},
		{"Two elements sorted", []int{1, 2}},
		{"Two elements reversed", []int{2, 1}},
		{"Small sorted", []int{1, 2, 3, 4, 5}},
		{"Small reversed", []int{5, 4, 3, 2, 1}},
		{"Small random", []int{3, 1, 4, 1, 5, 9, 2, 6, 5}},
		{"Duplicates", []int{1, 1, 1, 1, 1}},
		{"Mixed duplicates", []int{3, 1, 4, 1, 5, 3, 2, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := copySlice(tt.data)
			expected := copySlice(tt.data)
			sort.Ints(expected)

			ParallelQuickSort(data, comp)

			if !reflect.DeepEqual(data, expected) {
				t.Errorf("ParallelQuickSort() = %v, want %v", data, expected)
			}

			if !isSorted(data, comp) {
				t.Errorf("Result is not sorted: %v", data)
			}
		})
	}
}

func TestParallelQuickSortLargeData(t *testing.T) {
	comp := IntComparator{}
	sizes := []int{1000, 5000, 10000}

	for _, size := range sizes {
		t.Run(t.Name()+"_size_"+string(rune(size)), func(t *testing.T) {
			data := GenerateRandomInts(size)
			expected := copySlice(data)
			sort.Ints(expected)

			ParallelQuickSort(data, comp)

			if !isSorted(data, comp) {
				t.Errorf("Large data set (size %d) is not sorted", size)
			}

			if !reflect.DeepEqual(data, expected) {
				t.Errorf("Large data set (size %d) doesn't match expected result", size)
			}
		})
	}
}

// Тесты для SequentialQuickSort
func TestSequentialQuickSort(t *testing.T) {
	comp := IntComparator{}

	tests := []struct {
		name string
		data []int
	}{
		{"Empty slice", []int{}},
		{"Single element", []int{42}},
		{"Small random", []int{3, 1, 4, 1, 5, 9, 2, 6, 5}},
		{"Already sorted", []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := copySlice(tt.data)
			expected := copySlice(tt.data)
			sort.Ints(expected)

			SequentialQuickSort(data, comp)

			if !reflect.DeepEqual(data, expected) {
				t.Errorf("SequentialQuickSort() = %v, want %v", data, expected)
			}

			if !isSorted(data, comp) {
				t.Errorf("Result is not sorted: %v", data)
			}
		})
	}
}

// Тесты для partition
func TestPartition(t *testing.T) {
	comp := IntComparator{}

	tests := []struct {
		name string
		data []int
	}{
		{"Single element", []int{42}},
		{"Two elements", []int{2, 1}},
		{"Multiple elements", []int{3, 1, 4, 1, 5, 9, 2, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := copySlice(tt.data)

			if len(data) == 0 {
				return
			}

			pivotIndex := partition(data, comp)
			pivot := data[pivotIndex]

			// Проверяем, что элементы слева меньше пивота
			for i := 0; i < pivotIndex; i++ {
				if comp.Compare(data[i], pivot) >= 0 {
					t.Errorf("Element %v at index %d should be less than pivot %v", data[i], i, pivot)
				}
			}

			// Проверяем, что элементы справа больше или равны пивоту
			for i := pivotIndex + 1; i < len(data); i++ {
				if comp.Compare(data[i], pivot) < 0 {
					t.Errorf("Element %v at index %d should be >= pivot %v", data[i], i, pivot)
				}
			}
		})
	}
}

func TestPartitionEmptySlice(t *testing.T) {
	comp := IntComparator{}
	data := []int{}
	result := partition(data, comp)

	if result != 0 {
		t.Errorf("partition([]) should return 0, got %d", result)
	}
}

// Тесты для ParallelQuickSortWithThreshold
func TestParallelQuickSortWithThreshold(t *testing.T) {
	comp := IntComparator{}
	data := GenerateRandomInts(2000)
	expected := copySlice(data)
	sort.Ints(expected)

	tests := []struct {
		name      string
		threshold int
	}{
		{"Low threshold", 100},
		{"Medium threshold", 500},
		{"High threshold", 1500},
		{"Very high threshold", 5000}, // Должен использовать последовательную сортировку
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData := copySlice(data)

			ParallelQuickSortWithThreshold(testData, comp, tt.threshold)

			if !isSorted(testData, comp) {
				t.Errorf("Data is not sorted with threshold %d", tt.threshold)
			}

			if !reflect.DeepEqual(testData, expected) {
				t.Errorf("Result doesn't match expected with threshold %d", tt.threshold)
			}
		})
	}
}

// Тесты со строками
func TestParallelQuickSortStrings(t *testing.T) {
	comp := StringComparator{}
	data := []string{"zebra", "apple", "banana", "cherry", "date"}
	expected := []string{"apple", "banana", "cherry", "date", "zebra"}

	ParallelQuickSort(data, comp)

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("ParallelQuickSort(strings) = %v, want %v", data, expected)
	}

	if !isSorted(data, comp) {
		t.Errorf("String result is not sorted: %v", data)
	}
}

// Тесты на производительность
func BenchmarkParallelQuickSort(b *testing.B) {
	comp := IntComparator{}
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		data := GenerateRandomInts(size)
		b.Run("Random_"+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testData := copySlice(data)
				b.StartTimer()
				ParallelQuickSort(testData, comp)
				b.StopTimer()
			}
		})

		sortedData := generateSortedInts(size)
		b.Run("Sorted_"+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testData := copySlice(sortedData)
				b.StartTimer()
				ParallelQuickSort(testData, comp)
				b.StopTimer()
			}
		})

		reversedData := generateReversedInts(size)
		b.Run("Reversed_"+string(rune(size)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testData := copySlice(reversedData)
				b.StartTimer()
				ParallelQuickSort(testData, comp)
				b.StopTimer()
			}
		})
	}
}

func BenchmarkSequentialQuickSort(b *testing.B) {
	comp := IntComparator{}
	data := GenerateRandomInts(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testData := copySlice(data)
		b.StartTimer()
		SequentialQuickSort(testData, comp)
		b.StopTimer()
	}
}

// Сравнительные бенчмарки
func BenchmarkCompareParallelVsSequential(b *testing.B) {
	comp := IntComparator{}
	data := GenerateRandomInts(50000)

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testData := copySlice(data)
			b.StartTimer()
			ParallelQuickSort(testData, comp)
			b.StopTimer()
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testData := copySlice(data)
			b.StartTimer()
			SequentialQuickSort(testData, comp)
			b.StopTimer()
		}
	})
}

// Тесты на граничные случаи
func TestEdgeCases(t *testing.T) {
	comp := IntComparator{}

	t.Run("Very large slice", func(t *testing.T) {
		data := GenerateRandomInts(1000000)
		expected := copySlice(data)
		sort.Ints(expected)

		ParallelQuickSort(data, comp)

		if !isSorted(data, comp) {
			t.Error("Very large slice is not sorted")
		}
	})

	t.Run("All same elements", func(t *testing.T) {
		data := make([]int, 1000)
		for i := range data {
			data[i] = 42
		}

		ParallelQuickSort(data, comp)

		if !isSorted(data, comp) {
			t.Error("Slice with all same elements is not sorted")
		}

		for _, v := range data {
			if v != 42 {
				t.Error("Elements changed in all-same slice")
				break
			}
		}
	})
}

// Тест на корректность работы с горутинами
func TestConcurrencySafety(t *testing.T) {
	comp := IntComparator{}

	// Запускаем несколько сортировок параллельно
	const numGoroutines = 10
	const dataSize = 5000

	results := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			data := GenerateRandomInts(dataSize)
			ParallelQuickSort(data, comp)
			results <- isSorted(data, comp)
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		if !<-results {
			t.Error("Concurrent sorting failed")
		}
	}
}
