package test

import (
	"fmt"
	"math"
	"testing"
)

func main() {
	v := Abs(3)
	fmt.Println(v)
}

// Abs возвращает абсолютное значение.
// Например: 3.1 => 3.1, -3.14 => 3.14, -0 => 0.
// Покрыть тестами нужно эту функцию.
func Abs(value float64) float64 {
	return math.Abs(value)
}

func TestAbs(t *testing.T) {
	tests := []struct { // добавился слайс тестов
		name   string
		values float64
		want   float64
	}{
		{
			name:   "simple test #1", // описывается каждый тест
			values: -3,               // значения, которые будет принимать функция
			want:   3,                // ожидаемое значение
		},
		{
			name:   "one",
			values: 1,
			want:   1,
		},
		{
			name:   "with negative values",
			values: 3,
			want:   3,
		},
		{
			name:   "with negative zero",
			values: -2.000001,
			want:   2.000001,
		},
		{
			name:   "a lot of values",
			values: -0.000000003,
			want:   0.000000003,
		},
	}
	for _, tt := range tests { // цикл по всем тестам
		t.Run(tt.name, func(t *testing.T) {
			if res := Abs(tt.values); res != tt.want {
				t.Errorf("Add() = %v, want %v", res, tt.want)
			}
		})
	}
}
