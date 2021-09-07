package main

import "golang.org/x/sync/semaphore"

type PortScanner struct {
	ip string
	lock *semaphore.Weighted
}
