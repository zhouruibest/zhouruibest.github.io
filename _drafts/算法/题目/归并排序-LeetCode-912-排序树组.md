```go
func sortArray(nums []int) []int {
    temp := make([]int, len(nums))
    mergeSort(nums, temp, 0, len(nums)-1)
    return nums
}

// mergeSort归并排序
func mergeSort(nums, temp []int, left, right int) {
    // nums 要排序的数组
    // temp 临时数组，避免多次递归时创建新的
    // left, right 在nums上排序的范围

    // 结束的条件
    if left >= right {
        return
    }

    // 递
    mid := (left + right) / 2
    mergeSort(nums, temp, left, mid)
    mergeSort(nums, temp, mid+1, right)

    // 归
    merge(nums, temp, left, mid, right)
}

func merge(nums, temp []int, left, mid, right int) {
    a := left
    b := mid+1
    c := left
    for a <= mid && b <= right {
        if nums[a] <= nums[b] {
            temp[c] = nums[a]
            c++
            a++
        } else {
            temp[c] = nums[b]
            c++
            b++
        }
    }

    for a <= mid {
        temp[c] = nums[a]
        c++
        a++
    }

    for b <= right {
        temp[c] = nums[b]
        c++
        b++
    }

    for d:= left; d<= right; d++ {
        nums[d] = temp[d]
    }

}
```