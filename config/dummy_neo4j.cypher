MATCH (n) DETACH DELETE (n)
MERGE (n:Root {Text:"Apple"})
MERGE (a:Node {Text: "Pineapple"})
MERGE (a:Node {Text: "Granny Smith"})
MERGE (a:Node {Text: "Red Delicious"})
MERGE (a:Node {Text: "Fruit"})
MERGE (a:Node {Text: "Tree"})
MERGE (a:Node {Text: "Oak"})
MERGE (a:Node {Text: "Spruce"})
MERGE (a:Node {Text: "Redwood"})
MERGE (a:Node {Text: "Willow"})
MERGE (a:Node {Text: "Cherry"})
MERGE (a:Node {Text: "Pinecone"})
MERGE (a:Node {Text: "SpongeBob"})
MERGE (a:Node {Text: "Squidward"})
MERGE (a:Node {Text: "Patrick"})
MERGE (a:Node {Text: "Sandy"})
MERGE (a:Node {Text: "Mr. Krabs"})
MERGE (a:Node {Text: "Gary"})
MERGE (a:Node {Text: "AVL"})
MERGE (a:Node {Text: "Computer Science"})
MERGE (a:Node {Text: "Professor"})
MERGE (a:Node {Text: "Aaron Wilkin"})
MATCH (a:Root {Text:"AVL"}) MATCH (b:Node {Text: "Computer Science"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Computer Science"}) MATCH (b:Node {Text: "Professor"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Professor"}) MATCH (b:Node {Text: "Aaron Wilkin"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Apple"}) MATCH (b:Node {Text: "Pineapple"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Apple"}) MATCH (b:Node {Text: "Granny Smith"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Apple"}) MATCH (b:Node {Text: "Red Delicious"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Apple"}) MATCH (b:Node {Text: "Fruit"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Root {Text:"Apple"}) MATCH (b:Node {Text: "Tree"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "Oak"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "Spruce"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "Redwood"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "Willow"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "Cherry"}) MERGE (a) - [:Associated] - (b)
MATCH (a:Node {Text:"Tree"}) MATCH (b:Node {Text: "AVL"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"Pineapple"}) MATCH (b:Node {Text:"Pinecone"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"Pineapple"}) MATCH (b:Node {Text:"SpongeBob"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"SpongeBob"}) MATCH (b:Node {Text:"Squidward"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"SpongeBob"}) MATCH (b:Node {Text:"Patrick"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"SpongeBob"}) MATCH (b:Node {Text:"Sandy"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"SpongeBob"}) MATCH (b:Node {Text:"Mr. Krabs"}) MERGE (a) - [:Associated] - (b)
MATCH (a: Node {Text:"SpongeBob"}) MATCH (b:Node {Text:"Gary"}) MERGE (a) - [:Associated] - (b)