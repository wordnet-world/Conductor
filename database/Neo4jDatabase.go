package database

import (
	"bufio"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/wordnet-world/Conductor/models"
)

// Neo4jDatabase is a struct that implements the Graph interface
// providing access to the underlying graph used in the game
type Neo4jDatabase struct {
	driver neo4j.Driver
}

// NewNeo4jDatabase creates a new neo4j database
func NewNeo4jDatabase() *Neo4jDatabase {
	return &Neo4jDatabase{}
}

// Connect Initializes the connection to neo4j, creating the driver
func (db *Neo4jDatabase) Connect(uri, username, password string) error {
	useConsoleLogger := func(level neo4j.LogLevel) func(config *neo4j.Config) {
		return func(config *neo4j.Config) {
			config.Log = neo4j.ConsoleLogger(level)
		}
	}

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), useConsoleLogger(neo4j.ERROR))
	if err != nil {
		return err
	}
	db.driver = driver

	return nil
}

// PopulateDummy populates the dummy data
func (db *Neo4jDatabase) PopulateDummy(uri, username, password string) error {
	err := db.Connect(uri, username, password)
	if err != nil {
		return err
	}
	stage := os.Getenv("STAGE")
	if stage != "" {
		err := initializeWithDummyData(db.driver)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close closes the connection with the graph database
func (db *Neo4jDatabase) Close() {
	db.driver.Close()
}

func initializeWithDummyData(driver neo4j.Driver) error {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}

	file, err := os.Open("config/dummy_neo4j.cypher")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		query := scanner.Text()
		log.Println(query)
		_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(query, nil)
			return result, err
		})
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// GetRoot returns the root node
func (db *Neo4jDatabase) GetRoot() (models.Node, error) {
	roots, err := db.getNodes("MATCH (n:Root) RETURN n.Text, ID(n)", map[string]interface{}{})
	if err != nil {
		return models.Node{}, err
	}
	return roots[0], nil
}

// GetNeighbors returns the neighbors of a node
func (db *Neo4jDatabase) GetNeighbors(node models.Node) ([]models.Node, error) {
	neighbors, err := db.getNodes("MATCH (n) - [] - (a) MATCH (n) WHERE id(n)=$id RETURN a.Text, ID(a)", map[string]interface{}{"id": node.ID})
	if err != nil {
		return nil, err
	}
	return neighbors, nil
}

// GetNeighborsNodeID returns the neighbors of a node
func (db *Neo4jDatabase) GetNeighborsNodeID(nodeID int64) ([]models.Node, error) {
	neighbors, err := db.getNodes("MATCH (n) - [] - (a) MATCH (n) WHERE id(n)=$id RETURN a.Text, ID(a)", map[string]interface{}{"id": nodeID})
	if err != nil {
		return nil, err
	}
	return neighbors, nil
}

func (db *Neo4jDatabase) getNodes(query string, params map[string]interface{}) ([]models.Node, error) {
	session, err := db.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	response, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(query, params)

		if err != nil {
			return nil, err
		}
		nodes := make([]models.Node, 0)
		for {
			if result.Next() {
				text := result.Record().GetByIndex(0).(string)
				id := result.Record().GetByIndex(1).(int64)
				nodes = append(nodes, models.Node{
					Text: text,
					ID:   id,
				})
			} else {
				break
			}
		}
		return nodes, nil
	})
	if err != nil {
		return nil, err
	}

	final := response.([]models.Node)
	return final, nil
}
