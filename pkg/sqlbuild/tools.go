package sqlbuild

import "strings"

func handleStringsSplice(strings []string, separator string, builder *strings.Builder) {
	for i, value := range strings {
		if i > 0 {
			builder.WriteString(separator)
		}
		builder.WriteString(value)
	}
}

// SliceToAnySlice 将普通类型的切片转换为any类型切片
func SliceToAnySlice[T any](list []T) []any {
	if list == nil {
		return nil
	}
	result := make([]any, len(list))
	for i := 0; i < len(list); i++ {
		result[i] = list[i]
	}
	return result
}
