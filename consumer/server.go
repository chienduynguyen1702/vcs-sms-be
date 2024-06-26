package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// declare server model

type Server struct {
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
	fmt.Printf("| %15s | %15s | %s \n", s.IP, s.Status, s.Name)
}

func (s *Server) PrintResult() {
	fmt.Printf("|%20v| %15s | %15s | %s \n", s.PingAt.Format(DDMMYYYYhhmmss), s.IP, s.Status, s.Name)
}

func (s *Server) Print() {
	fmt.Printf("| %15s | %15s | %s \n", "     IP    ", "    Status   ", "     Name ")
	fmt.Printf("| %15s | %15s | %s \n", s.IP, s.Status, s.Name)
}

func (s *Server) UnmarshalJSON(b []byte) error {
	type serverAlias Server
	var server serverAlias
	err := json.Unmarshal(b, &server)
	if err != nil {
		return err
	}
	*s = Server(server)
	return nil
}

type ListServers []Server

func Print(ListServers []Server) {
	fmt.Println("")
	fmt.Printf("Found %d servers\n", len(ListServers))
	fmt.Println("")
	fmt.Printf("| %15s | %15s | %s \n", "     IP    ", "    Status   ", "     Name ")
	fmt.Printf("|%17s|%17s|%s\n", "-----------------", "-----------------", "----------------------------------")
	for _, server := range ListServers {
		server.PrintOne()
	}
}
