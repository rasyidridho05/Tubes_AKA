package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const MIN_MERGE = 32 

func insertionSort(arr []int, left, right int) {
	for i := left + 1; i <= right; i++ {
		key := arr[i]
		j := i - 1

		for j >= left && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func merge(arr []int, left, mid, right int) {

	leftArr := make([]int, mid-left+1)
	rightArr := make([]int, right-mid)

	copy(leftArr, arr[left:mid+1])
	copy(rightArr, arr[mid+1:right+1])

	i, j, k := 0, 0, left

	for i < len(leftArr) && j < len(rightArr) {
		if leftArr[i] <= rightArr[j] {
			arr[k] = leftArr[i]
			i++
		} else {
			arr[k] = rightArr[j]
			j++
		}
		k++
	}

	for i < len(leftArr) {
		arr[k] = leftArr[i]
		i++
		k++
	}

	for j < len(rightArr) {
		arr[k] = rightArr[j]
		j++
		k++
	}
}

func timSortRecursiveHelper(arr []int, left, right, minMerge int) {
	size := right - left + 1

	if size <= minMerge {
		insertionSort(arr, left, right)
		return
	}

	mid := left + (right-left)/2

	timSortRecursiveHelper(arr, left, mid, minMerge)
	timSortRecursiveHelper(arr, mid+1, right, minMerge)

	merge(arr, left, mid, right)
}

func timSortRecursive(arr []int) {
	if len(arr) <= 1 {
		return
	}
	timSortRecursiveHelper(arr, 0, len(arr)-1, MIN_MERGE)
}

func timSortIterative(arr []int) {
	n := len(arr)
	for start := 0; start < n; start += MIN_MERGE {
		end := min(start+MIN_MERGE-1, n-1)
		insertionSort(arr, start, end)
	}
	size := MIN_MERGE

	for size < n {
		for start := 0; start < n; start += size * 2 {
			mid := start + size - 1
			end := min(start+size*2-1, n-1)

			if mid < end {
				merge(arr, start, mid, end)
			}
		}
		size *= 2
	}
}


func generateData(n int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = rand.Intn(n)
	}
	return data
}

func measureTime(fn func()) float64 {
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	return float64(elapsed.Nanoseconds()) / 1_000_000.0 
}

func copySlice(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	return result
}

type Result struct {
	n    int
	rec  float64
	iter float64
}

func printTable(results []Result) {
	fmt.Println("\n╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║          TIM SORT - PERBANDINGAN WAKTU EKSEKUSI       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝\n")
	fmt.Println("-------------------------------------------------")
	fmt.Println("|   N Data   | Recursive (ms) | Iterative (ms) |")
	fmt.Println("-------------------------------------------------")

	var sumRec, sumIter float64
	for _, row := range results {
		fmt.Printf("| %-10d | %-14.3f | %-14.3f |\n", row.n, row.rec, row.iter)
		sumRec += row.rec
		sumIter += row.iter
	}
	fmt.Println("-------------------------------------------------")

	avgRec := sumRec / float64(len(results))
	avgIter := sumIter / float64(len(results))
	fmt.Printf("| %-10s | %-14.3f | %-14.3f |\n", "AVERAGE", avgRec, avgIter)
	fmt.Println("-------------------------------------------------")
}

func printChart(results []Result) {
	fmt.Println("\n╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║              GRAFIK PERBANDINGAN (█ = 1 ms)           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝\n")

	for _, r := range results {
		fmt.Printf("n = %d\n", r.n)
		fmt.Printf("Recursive : %s\n", strings.Repeat("█", int(math.Round(r.rec))))
		fmt.Printf("Iterative : %s\n", strings.Repeat("█", int(math.Round(r.iter))))
		fmt.Println()
	}
}

func askInput() []int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan ukuran data (pisahkan dengan spasi, contoh: 1000 5000 10000): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Fields(input)
	sizes := make([]int, 0)

	for _, part := range parts {
		if num, err := strconv.Atoi(part); err == nil && num > 0 {
			sizes = append(sizes, num)
		}
	}

	return sizes
}

func main() {
	fmt.Printf("MIN_MERGE size: %d elements\n\n", MIN_MERGE)

	sizes := askInput()

	if len(sizes) == 0 {
		fmt.Println("\n❌ Input tidak valid")
		return
	}

	results := make([]Result, 0)

	fmt.Println("\n⏳ Sedang melakukan benchmark...\n")

	for _, n := range sizes {
		data := generateData(n)

		arrRec := copySlice(data)
		recTime := measureTime(func() {
			timSortRecursive(arrRec)
		})

		arrIter := copySlice(data)
		iterTime := measureTime(func() {
			timSortIterative(arrIter)
		})

		results = append(results, Result{
			n:    n,
			rec:  recTime,
			iter: iterTime,
		})

		fmt.Printf("✓ Completed: n = %d\n", n)
	}

	printTable(results)
	printChart(results)
}