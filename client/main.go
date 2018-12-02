package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/adigunhammedolalekan/araknet"
	"log"
	"net"
	"os"
	"strings"
)

func printInstruction() {

	fmt.Println("Commands ****")
	fmt.Println("- Create Network -> create_network [host] [network_name] [username] [secret]")
	fmt.Println("- Join Network -> join [host] [network_name] [username] [secret]")
	fmt.Println("- Quit -> quit")

	fmt.Println()
	fmt.Print(">>> ")
}
func createPrompt() {

	printInstruction()

	for {

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		input = strings.Replace(input, "\n", "", -1)
		fmt.Println(input)
		if input == "quit" {
			break
		}

		processCommand(input)
		fmt.Print(">>> ")

	}

}

func main() {
	createPrompt()
}

func processCommand(input string)  {

	commands := strings.Split(input, " ")
	count := len(commands)
	if input != "quit" && (count < 5 || count > 5) {

		printInstruction()
		return
	}

	switch commands[0] {

	case araknet.TYPE_CREATE_NETWORK:
		createNetwork(commands[1], commands[2], commands[3], commands[4])
	case araknet.TYPE_JOIN_NETWORK:
		joinNetwork(commands[1], commands[2], commands[3], commands[4])
	}
}

func createNetwork(host string, username string, network string, secret string) {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal("Failed to connect to " + host, err)
	}

	defer conn.Close()

	message := &araknet.Message{}
	message.Name = username
	message.Network = network
	message.Type = araknet.TYPE_CREATE_NETWORK
	message.Secret = secret

	bytes, _ := json.Marshal(message)
	conn.Write(bytes)

	decoder := json.NewDecoder(conn)

	for  {

		fmt.Println(">>> ")
		m := &araknet.Message{}
		err := decoder.Decode(message)
		if err != nil {
			fmt.Println("Failed to decode message => ", err)
			return
		}

		if message.Name != m.Name {

			data := m.Format()
			fmt.Println(data)
		}
	}
}

func joinNetwork(host string, username string, network string, secret string) {

	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal("Failed to connect to " + host, err)
	}

	defer conn.Close()

	message := &araknet.Message{}
	message.Name = username
	message.Network = network
	message.Type = araknet.TYPE_JOIN_NETWORK
	message.Secret = secret

	bytes, _ := json.Marshal(message)
	conn.Write(bytes)

	decoder := json.NewDecoder(conn)

	for  {

		m := &araknet.Message{}
		err := decoder.Decode(message)
		if err != nil {
			fmt.Println("Failed to decode message => ", err)
			return
		}

		if message.Name != m.Name {

			data := m.Format()
			fmt.Println(data)
		}
	}
}
