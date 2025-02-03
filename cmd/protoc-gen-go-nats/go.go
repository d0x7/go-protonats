package main

import "strings"

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }
