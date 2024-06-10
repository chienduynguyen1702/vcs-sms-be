package main

import "fmt"

// declare server model

type Server struct {
	IP     string
	Name   string
	Status string
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

func (s *Server) Print() {
	fmt.Printf("| %15s | %15s | %s \n", "     IP    ", "    Status   ", "     Name ")
	fmt.Printf("| %15s | %15s | %s \n", s.IP, s.Status, s.Name)
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
