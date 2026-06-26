package tools

import (
	"golang.org/x/exp/slices"
	"strings"
)

// StrInContains 用于判断字符串是否 包含某个集合当中的字符串
func StrInContains(u string, contain []string) bool {
	return slices.ContainsFunc(contain, func(s string) bool {
		return strings.Contains(u, s)
	})
}

// LevenshteinDistance 计算两个字符串之间的Levenshtein距离
func LevenshteinDistance(a, b string) int {
	// 初始化矩阵
	lenA := len(a)
	lenB := len(b)
	dp := make([][]int, lenA+1)
	for i := range dp {
		dp[i] = make([]int, lenB+1)
	}

	// 初始化第一行和第一列
	for i := range dp {
		dp[i][0] = i
	}
	for j := range dp[0] {
		dp[0][j] = j
	}

	// 计算Levenshtein距离
	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] // 字符相同
			} else {
				dp[i][j] = min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]) + 1
			}
		}
	}

	return dp[lenA][lenB]
}
