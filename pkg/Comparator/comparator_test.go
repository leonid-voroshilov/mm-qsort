package comparator

import (
	"testing"
)

func TestIntComparator(t *testing.T) {
	comp := IntC{}

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "a less than b",
			a:        5,
			b:        10,
			expected: -1,
		},
		{
			name:     "a greater than b",
			a:        15,
			b:        7,
			expected: 1,
		},
		{
			name:     "a equals b",
			a:        42,
			b:        42,
			expected: 0,
		},
		{
			name:     "negative numbers - a less than b",
			a:        -10,
			b:        -5,
			expected: -1,
		},
		{
			name:     "negative numbers - a greater than b",
			a:        -3,
			b:        -8,
			expected: 1,
		},
		{
			name:     "negative numbers equal",
			a:        -7,
			b:        -7,
			expected: 0,
		},
		{
			name:     "positive and negative - positive greater",
			a:        5,
			b:        -3,
			expected: 1,
		},
		{
			name:     "positive and negative - negative less",
			a:        -2,
			b:        8,
			expected: -1,
		},
		{
			name:     "zero comparisons",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "zero vs positive",
			a:        0,
			b:        1,
			expected: -1,
		},
		{
			name:     "zero vs negative",
			a:        0,
			b:        -1,
			expected: 1,
		},
		{
			name:     "maximum int values",
			a:        2147483647, // max int32
			b:        2147483647,
			expected: 0,
		},
		{
			name:     "minimum int values",
			a:        -2147483648, // min int32
			b:        -2147483648,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := comp.Compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("IntC.Compare(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestStringComparator(t *testing.T) {
	comp := StringC{}

	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{
			name:     "a less than b lexicographically",
			a:        "apple",
			b:        "banana",
			expected: -1,
		},
		{
			name:     "a greater than b lexicographically",
			a:        "zebra",
			b:        "apple",
			expected: 1,
		},
		{
			name:     "equal strings",
			a:        "hello",
			b:        "hello",
			expected: 0,
		},
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: 0,
		},
		{
			name:     "empty vs non-empty - empty less",
			a:        "",
			b:        "a",
			expected: -1,
		},
		{
			name:     "non-empty vs empty - non-empty greater",
			a:        "a",
			b:        "",
			expected: 1,
		},
		{
			name:     "case sensitivity - uppercase less than lowercase",
			a:        "A",
			b:        "a",
			expected: -1,
		},
		{
			name:     "case sensitivity - lowercase greater than uppercase",
			a:        "z",
			b:        "Z",
			expected: 1,
		},
		{
			name:     "numbers as strings",
			a:        "1",
			b:        "2",
			expected: -1,
		},
		{
			name:     "numbers vs letters",
			a:        "9",
			b:        "A",
			expected: -1,
		},
		{
			name:     "different lengths - shorter is prefix",
			a:        "cat",
			b:        "catch",
			expected: -1,
		},
		{
			name:     "different lengths - longer contains shorter",
			a:        "catch",
			b:        "cat",
			expected: 1,
		},
		{
			name:     "special characters",
			a:        "hello!",
			b:        "hello?",
			expected: -1,
		},
		{
			name:     "unicode characters",
			a:        "café",
			b:        "cafe",
			expected: 1,
		},
		{
			name:     "whitespace",
			a:        "hello world",
			b:        "hello  world",
			expected: 1,
		},
		{
			name:     "single character strings",
			a:        "a",
			b:        "b",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := comp.Compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("StringC.Compare(%q, %q) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Тест для проверки, что компараторы реализуют интерфейс Comparator
func TestComparatorInterface(t *testing.T) {
	t.Run("IntC implements Comparator", func(t *testing.T) {
		var _ Comparator[int] = IntC{}
	})

	t.Run("StringC implements Comparator", func(t *testing.T) {
		var _ Comparator[string] = StringC{}
	})
}

// Benchmark тесты
func BenchmarkIntComparator(b *testing.B) {
	comp := IntC{}
	for i := 0; i < b.N; i++ {
		comp.Compare(i, i+1)
	}
}

func BenchmarkStringComparator(b *testing.B) {
	comp := StringC{}
	s1 := "benchmark"
	s2 := "testing"

	for i := 0; i < b.N; i++ {
		comp.Compare(s1, s2)
	}
}

// Тест для проверки транзитивности (если a < b и b < c, то a < c)
func TestIntComparatorTransitivity(t *testing.T) {
	comp := IntC{}

	testCases := []struct {
		a, b, c int
	}{
		{1, 2, 3},
		{-3, -1, 5},
		{0, 10, 20},
	}

	for _, tc := range testCases {
		ab := comp.Compare(tc.a, tc.b)
		bc := comp.Compare(tc.b, tc.c)
		ac := comp.Compare(tc.a, tc.c)

		if ab == -1 && bc == -1 && ac != -1 {
			t.Errorf("Transitivity failed: %d < %d < %d, but Compare(%d, %d) = %d",
				tc.a, tc.b, tc.c, tc.a, tc.c, ac)
		}
	}
}

// Тест для проверки симметрии (если Compare(a, b) = x, то Compare(b, a) = -x)
func TestIntComparatorSymmetry(t *testing.T) {
	comp := IntC{}

	testCases := [][2]int{
		{5, 10},
		{-3, 7},
		{0, 0},
		{15, -2},
	}

	for _, tc := range testCases {
		ab := comp.Compare(tc[0], tc[1])
		ba := comp.Compare(tc[1], tc[0])

		if ab != -ba {
			t.Errorf("Symmetry failed: Compare(%d, %d) = %d, Compare(%d, %d) = %d",
				tc[0], tc[1], ab, tc[1], tc[0], ba)
		}
	}
}

func TestStringComparatorSymmetry(t *testing.T) {
	comp := StringC{}

	testCases := [][2]string{
		{"apple", "banana"},
		{"", "test"},
		{"same", "same"},
		{"Z", "a"},
	}

	for _, tc := range testCases {
		ab := comp.Compare(tc[0], tc[1])
		ba := comp.Compare(tc[1], tc[0])

		if ab != -ba {
			t.Errorf("Symmetry failed: Compare(%q, %q) = %d, Compare(%q, %q) = %d",
				tc[0], tc[1], ab, tc[1], tc[0], ba)
		}
	}
}
