package untrust_test

import (
	"testing"

	"github.com/flier/goutil/pkg/untrust"
)

func BenchmarkInput_Clone(b *testing.B) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	input := untrust.Input(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = input.Clone()
	}
}

func BenchmarkInput_AsSliceLessSafe(b *testing.B) {
	input := untrust.Input([]byte("hello world"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = input.AsSliceLessSafe()
	}
}

func BenchmarkReader_ReadByte(b *testing.B) {
	data := make([]byte, b.N)
	for i := range data {
		data[i] = byte(i % 256)
	}
	input := untrust.Input(data)
	r := untrust.NewReader(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.ReadByte()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader_ReadBytes(b *testing.B) {
	data := make([]byte, b.N)
	for i := range data {
		data[i] = byte(i % 256)
	}
	input := untrust.Input(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := untrust.NewReader(input)

		_, err := r.ReadBytes(len(data))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader_Clone(b *testing.B) {
	input := untrust.Input([]byte("hello world"))
	r := untrust.NewReader(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Clone()
	}
}

func BenchmarkReadAll(b *testing.B) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	input := untrust.Input(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (string, error) {
			// Read all bytes
			bytes, err := r.ReadBytesToEnd()
			if err != nil {
				return "", err
			}

			return string(bytes), nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadPartial(b *testing.B) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	input := untrust.Input(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset r for each iteration
		r := untrust.NewReader(input)

		_, _, err := untrust.ReadPartial(r, func(r *untrust.Reader) (string, error) {
			bytes, err := r.ReadBytes(100)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
