package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// declare server model

type Server struct {
	gorm.Model
	IP       string    `json:"ip"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	PingAt   time.Time `json:"ping_at"`
	IsOnline bool      `json:"is_online"`
}

// method to make []Server a ListServers type
func WrapServers(servers []Server) ListServers {
	var l ListServers
	for _, server := range servers {
		l = append(l, server)
	}
	return l
}
func (s *Server) PrintOne() {
	fmt.Printf("| %15s | %15s | %18s | %s\n", s.IP, s.Status, s.Name, s.PingAt.Format("2006-01-02 15:04:05"))
}

func (s *Server) Print() {
	fmt.Printf("| %15s | %15s | %s \n", "     IP    ", "    Status   ", "     Name ")
	fmt.Printf("| %15s | %15s | %s \n", s.IP, s.Status, s.Name)
}

type ListServers []Server

func PrintListServers(ListServers []Server) {
	fmt.Println("")
	fmt.Printf("Found %d servers\n", len(ListServers))
	fmt.Println("")
	fmt.Printf("| %15s | %15s | %18s | %s |\n", "     IP    ", "    Status   ", "     Name    ", "	 Ping At    ")
	fmt.Printf("|%17s|%17s|%20s| %s \n", "-----------------", "-----------------", "--------------------", "----------------------------------")
	for _, server := range ListServers {
		server.PrintOne()
	}
}
