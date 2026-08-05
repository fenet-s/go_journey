package main

import (
	"fmt"
	"regexp"
	"strings"
)

// IsPalindrome checks if a string is a palindrome.
func IsPalindrome(text string) bool {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove punctuation and spaces
	re := regexp.MustCompile(`[^a-z0-9]`)
	text = re.ReplaceAllString(text, "")

	// Check from both ends
	left := 0
	right := len(text) - 1

	for left < right {
		if text[left] != text[right] {
			return false
		}
		left++
		right--
	}

	return true
}

func main() {
	fmt.Println(IsPalindrome("man"))
	fmt.Println(IsPalindrome("apple"))
	fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println(IsPalindrome("madam"))
}
