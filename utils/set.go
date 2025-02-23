package utils

type Set[T comparable] map[T]struct{}

// NewSet 创建新的 Set
func NewSet[T comparable]() Set[T] {
	m := make(map[T]struct{})
	return m
}

// Add 向 Set 中添加新元素
func (s Set[T]) Add(ele T) {
	s[ele] = struct{}{}
}

// Delete 从 Set 中移除元素
func (s Set[T]) Delete(ele T) {
	delete(s, ele)
}

// Has 如果值存在则返回 true
func (s Set[T]) Has(ele T) bool {
	_, exists := s[ele]
	return exists
}

// Clear 从 Set 中移除所有元素
func (s Set[T]) Clear() {
	s.ForEach(func(ele T) {
		delete(s, ele)
	})
}

// ForEach 为每个元素调用回调函数
func (s Set[T]) ForEach(f func(ele T)) {
	for key := range s {
		f(key)
	}
}

// Values 返回包含 Set 中所有值的迭代器
func (s Set[T]) Values() []T {
	arr := make([]T, 0, s.Size())
	s.ForEach(func(ele T) {
		arr = append(arr, ele)
	})
	return arr
}

// Keys 与 values() 相同
func (s Set[T]) Keys() []T {
	return s.Values()
}

// Entries 返回切片，其中包含 Set 中的 [value,value] 值值对
func (s Set[T]) Entries() map[T]T {
	arr := make(map[T]T, s.Size())
	for key := range s {
		arr[key] = key
	}
	return arr
}

// Size 获取set长度
func (s Set[T]) Size() int {
	return len(s)
}
