//go:build !386 && !arm && linux

package raft

const initialMmapSize = 100 * 1024 * 1024 * 1024 // 100GB
