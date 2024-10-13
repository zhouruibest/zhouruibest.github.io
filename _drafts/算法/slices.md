# Max

numbers := []int{0, 10, -1, 8}
fmt.Println(slices.Max(numbers)) // 10

# MaxFunc

firstOldest := slices.MaxFunc(people, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age) // cmp包提供了常用的比较函数
	})

# Replace

names := []string{"Alice", "Bob", "Vera", "Zac"}
names = slices.Replace(names, 1, 3, "Bill", "Billie", "Cat")

# Reverse

names := []string{"alice", "Bob", "VERA"}
slices.Reverse(names)

# Sort 升序排列

s1 := []int8{0, 42, -10, 8}
slices.Sort(s1)

# SortFunc

比较函数的原型是func(a,b E) int, 是一个Less函数，E的类型和数组的是一样的

slices.SortFunc(names, func(a, b string) int {
	return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
})


# Index

index := slices.Index(nums, target) // index or -1

# IndexFunc
// IndexFunc returns the first index i satisfying f(s[i]), or -1 if none do.
i := slices.IndexFunc(numbers, func(n int) bool {
    return n < 0
})