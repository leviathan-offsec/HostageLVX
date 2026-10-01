//go:build !windows

package main

func init() { enableVT() }

func enableVT() {}
