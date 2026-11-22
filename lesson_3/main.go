package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

type Provider struct {
	Name  string
	Calls int
}

// ---------- Шаг 1. Указатели, стек и куча ----------

// TODO (шаг 1.1): увеличь p.Calls на 1 и напечатай адрес p через %p.
func countByValue(p Provider) {
	p.Calls++
	fmt.Printf("адрес p: %p\n", &p)
}

// TODO (шаг 1.1): увеличь p.Calls на 1 и напечатай адрес, который лежит в p, через %p.
func countByPointer(p *Provider) {
	p.Calls++
	fmt.Printf("адрес p: %p\n", &p)
}

// TODO (шаг 1.2): создай внутри функции Provider с именем name
// и верни его Calls + 1. Сам провайдер из функции не возвращай.
func localProvider(name string) int {
	p := Provider{Name: name}
	p.Calls++
	return p.Calls
}

// TODO (шаг 1.2): создай внутри функции Provider с именем name
// и верни указатель на него.
func newProvider(name string) *Provider {
	p := Provider{Name: name}
	return &p
}

func step1() {
	fmt.Println("--- Шаг 1 ---")

	// TODO (шаг 1.1): создай провайдера, напечатай его адрес через %p.
	// Вызови countByValue, напечатай Calls.
	// Вызови countByPointer, напечатай Calls.
	google := Provider{Name: "google"}
	countByValue(google)
	countByPointer(&google)

	// TODO (шаг 1.2): вызови localProvider и newProvider, напечатай результаты.
	fmt.Printf("%s\n", localProvider("google"))
	fmt.Printf("%p\n", newProvider("google"))
}

// ---------- Шаг 2. Массив и слайс ----------

// TODO (шаг 2.2): напечатай имя, len, cap и адрес массива, на который смотрит слайс.
// Адрес массива: unsafe.SliceData(s), печатается через %p.
func describe(name string, s []string) {
	fmt.Println(name)
	fmt.Println(len(s))
	fmt.Println(cap(s))
	fmt.Printf("%p\n", unsafe.SliceData(s))
}

func step2() {
	fmt.Println("--- Шаг 2 ---")

	// TODO (шаг 2.1): создай массив из 3 строк ([3]string) и скопируй его в другую переменную.
	// Поменяй в копии первый элемент. Напечатай обе переменные.
	// Потом сделай то же самое со слайсом из 3 строк ([]string).
	days := [3]string{"mon", "tue", "wed"}
	week := days
	week[0] = "sun"
	fmt.Println(days)
	fmt.Println(week)

	slice_days := []string{"mon", "tue", "wed"}
	slice_week := slice_days
	slice_week[0] = "sun"
	fmt.Println(slice_days)
	fmt.Println(slice_week)

	// TODO (шаг 2.1): напечатай unsafe.Sizeof для:
	//   - массива [3]string и массива [100]string;
	string3 := [3]string{}
	string100 := [100]string{}
	fmt.Println(unsafe.Sizeof(string3))
	fmt.Println(unsafe.Sizeof(string100))

	//   - слайса из 3 строк и слайса из 100 строк (make([]string, 100)).
	slice3 := make([]string, 3)
	slice100 := make([]string, 100)
	fmt.Println(unsafe.Sizeof(slice3))
	fmt.Println(unsafe.Sizeof(slice100))

	// TODO (шаг 2.2): создай слайс questions из 5 строк.
	// Вызови describe для:
	//   - questions;
	//   - same := questions;
	//   - part := questions[1:3];
	//   - grown := append(questions, "q6").
	questions := []string{"q1", "q2", "q3", "q4", "q5"}
	same := questions
	part := questions[1:3]
	grown := append(questions, "q6")

	describe("questions", questions)
	describe("same", same)
	describe("part", part)
	describe("grown", grown)
}

// ---------- Шаг 3. Предвыделение памяти ----------

// fillNoPrealloc собирает слайс из n чисел, не зная размера заранее.
func fillNoPrealloc(n int) []int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

// TODO (шаг 3.1): собери такой же слайс из n чисел,
// но выдели память под него один раз, до цикла.
func fillPrealloc(n int) []int {
	return nil
}

// TODO (шаг 3.2): добавляй в слайс n чисел через append и считай,
// сколько раз за это время изменилась ёмкость (cap). Верни это число.
func countGrowths(n int) int {
	return 0
}

// measure уже готова, менять её не нужно.
// Она вызывает fill(n) несколько раз (rounds) и печатает средние значения
// на один вызов: сколько было выделений памяти в куче (аллокаций),
// сколько байт выделено и сколько времени это заняло.
func measure(name string, n int, rounds int, fill func(int) []int) {
	runtime.GC()

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	start := time.Now()

	total := 0
	for i := 0; i < rounds; i++ {
		total += len(fill(n))
	}

	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)

	r := uint64(rounds)
	fmt.Printf("%-18s n=%-8d собрано=%-8d аллокаций=%-3d байт=%-9d время=%v\n",
		name,
		n,
		total/rounds,
		(after.Mallocs-before.Mallocs)/r,
		(after.TotalAlloc-before.TotalAlloc)/r,
		elapsed/time.Duration(rounds),
	)
}

func step3() {
	fmt.Println("--- Шаг 3 ---")

	measure("без предвыделения", 100, 100000, fillNoPrealloc)
	measure("с предвыделением", 100, 100000, fillPrealloc)

	measure("без предвыделения", 1000000, 20, fillNoPrealloc)
	measure("с предвыделением", 1000000, 20, fillPrealloc)

	fmt.Println("ёмкость менялась раз:", countGrowths(100), "для n=100")
	fmt.Println("ёмкость менялась раз:", countGrowths(1000000), "для n=1000000")
}

func main() {
	step1()
	step2()
	// step3()
}
