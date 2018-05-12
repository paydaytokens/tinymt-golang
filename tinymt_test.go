package tinymt

import "fmt"
import "testing"

func TestTinyMT64Init(t *testing.T) {
	Tinymt64_init(0xffff4332)
}

func TestTinyMT64Generate(t *testing.T) {
	for seeds := uint64(0); seeds < uint64(10); seeds++ {
		gen := Tinymt64_init(0xffff4332 * seeds)
		cycles := 1000000
		classes := uint64(10)
		nums := make(map[uint64]int, cycles)
		for i := 0; i < cycles; i++ {
			nums[gen.Generate()%classes] += 1
		}
		fmt.Printf("%v\n", nums)
	}
}

func BenchmarkTinyMT64Init(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Tinymt64_init(0xdeadbeefdead)
	}
	b.StopTimer()
}

func BenchmarkTinyMT64Generate(b *testing.B) {
	gen := Tinymt64_init(0xdeadbeefdead)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.Generate()
	}

}
