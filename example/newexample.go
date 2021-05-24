package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	croc "github.com/krickler/crocgodyl"
)

var (
	// Config is the global config variable.
	Config config
)

type config struct {
	PanelURL    string `json:"panel_url"`
	APIToken    string `json:"api_token"`
	ClientToken string `json:"client_token"`
}

func init() {
	//log.SetOutput(os.Stdout)
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("Error loading config.")
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Config)
}

func main() {
	application()
}

func application() {
	panel, err := croc.NewApp(Config.PanelURL, Config.APIToken)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("panel connection successful")
	}

	// validate the server is up and available
	if users, err := panel.GetUsers(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(users.Users[0].Attributes.Username)
	}

	printAllNodes(panel)

	newid := createNode(panel)
	time.Sleep(time.Second * 2)
	readNode(panel, newid)
	updateNode(panel, newid)
	time.Sleep(time.Second * 2)
	deleteNode(panel, newid)
}

func printAllNodes(panel *croc.AppConfig) {
	// All AppNodes
	fmt.Println("Listing all nodes on the panel.")
	nodesData, err := panel.GetNodes()

	if err != nil {
		log.Println("There was an error getting the locations.")
	}

	fmt.Println("All nodes on the panel")
	for _, node := range nodesData.Nodes {
		fmt.Print(node)
		//fmt.Printf("ID: %d Name: %s\n", node.Attributes.ID, node.Attributes.Name)
	}
}

func createNode(panel *croc.AppConfig) int {
	newNode := croc.NodeAttributes{
		Name:            "apitest",
		LocationID:      2,
		FQDN:            "test.ruecker.cloud",
		Scheme:          "https",
		Memory:          1500,
		MemoryOverAlloc: 0,
		Disk:            19000,
		DiskOverAlloc:   0,
		UploadSize:      100,
		DaemonListen:    8080,
		DaemonSftp:      2022,
	}

	newNodeInfo, err := panel.CreateNode(newNode)
	if err != nil {
		log.Println("Failed to create node.")
		log.Println(err)
	} else {
		log.Println("node created successfully.")
		fmt.Println("node info")
		fmt.Printf("ID: %d Name: %s\n", newNodeInfo.Attributes.ID, newNodeInfo.Attributes.Name)
	}

	newNodeAllocations := croc.AllocationAttributes{
		IP:    "2.2.2.2",
		Alias: "two.two.two.two",
		Ports: []string{"4000", "4001", "4002-4500"},
	}

	err = panel.CreateNodeAllocations(newNodeAllocations, newNodeInfo.Attributes.ID)
	if err != nil {
		log.Println("Failed to add node allocations.")
		log.Println(err)
	} else {
		log.Println("node allocations added successfully.")
	}
	return newNodeInfo.Attributes.ID
}

func readNode(panel *croc.AppConfig, nodeID int) {
	node, err := panel.GetNode(nodeID)
	if err == nil {
		fmt.Println(node.Attributes)
	} else {
		log.Println("Error reading node")
		log.Println(err)
	}

}

func updateNode(panel *croc.AppConfig, nodeID int) {
	editNode := croc.NodeAttributes{
		Name:            "apitest",
		LocationID:      2,
		FQDN:            "test.ruecker.cloud",
		Scheme:          "https",
		Memory:          1500,
		MemoryOverAlloc: 0,
		Disk:            19000,
		DiskOverAlloc:   0,
		UploadSize:      500,
		DaemonListen:    8080,
		DaemonSftp:      2022,
	}

	editNodeInfo, err := panel.EditNode(editNode, nodeID)
	if err != nil {
		log.Println("Failed to edit node.")
		log.Println(err)
	} else {
		log.Println("node edited successfully.")
		fmt.Println("node info")
		fmt.Printf("ID: %d Name: %s\n", editNodeInfo.Attributes.ID, editNodeInfo.Attributes.Name)
	}
}
func deleteNode(panel *croc.AppConfig, nodeID int) {
	err := panel.DeleteNode(nodeID)
	if err != nil {
		log.Println("failed to delete node")
		log.Println(err)
	} else {
		log.Println("Deleted node")
	}
}
