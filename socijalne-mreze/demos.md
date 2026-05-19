# Content of the Demos

## Demo 1

This demo shows how you can create an undirected graph with the `gonum` library. 

## Demo 2

This demo shows how you can load some well-known social networks. The graphs can be found in the project's root directory as `.txt` files. Each line contains one edge in the format `ID_Node1` `spacebar` `ID_Node2`. Using the implemented `LoadGraph` function found in `lab/graph.go`, you can load any graph you want.

Every demo after this uses the functions declared for loading the graphs in the root directory.

## Demo 3

This demo presents the degree distribution of every node in the graph.

The degree of a node is the number of edges that node has.

## Demo 4

This demo shows the clustering coefficient for each node in the graph.

The clustering coefficient indicates whether a node is in a cluster. A value closer to 1 means that almost all of its neighbors have edges between them; a value closer to 0 means that the node is not part of a cluster but more of a bridge between clusters.

The formula for the clustering coefficient for a node is the following:

```
CC = num_of_edges_between_its_neighbors / num_of_possible_edges_between_its_neighbors
```

## Demo 5

This demo shows how many connected components a graph has. First it showcases this on one of the graphs in the root directory, then on a generated Erdos-Renyi graph.

A connected component of a graph is a subgraph in which there is a path between every pair of nodes.
Erdos-Renyi graph generation is a technique to generate a random graph. You are given two parameters: the number of nodes and the probability for every edge to exist. Using these two parameters you can generate the graph. You can see the implementation in `demo5/helper.go`.

## Demo 6

This demo shows the shortest distance between every node in the graph. It also shows the average shortest distance between nodes and the diameter of the graph.

It uses the function `path.DijkstraAllPaths` that returns the object `path.AllShortest`. You can call the method `Between` on that object to get information about the shortest path between nodes.

## Demo 7

This demo prints out the top 10 nodes for betweenness centrality, closeness centrality and eigenvector centrality.

Betweenness centrality measures how many times a node appears on shortest paths.
Closeness centrality represents how close a node is to the rest of the network.
Eigenvector centrality represents how important a node is in a network.

## Demo 8

This demo prints out the shell index for every node in the graph. It also prints the sizes of the k-core subgraphs.

A k-core graph is a graph where all of the nodes have degree of at least k.
The shell index of a node is k, where k is the highest k-core subgraph the node belongs to.

## Demo 9

This demo shows algorithms for community detection. You can see label propagation, greedy modularity optimization and the Louvain method.

Label propagation is an iterative method. At the start, each node has its label, and after every iteration it takes the label that is most popular among its neighbors. This process continues until the labels stop changing.
Greedy modularity optimization is also an iterative method. Modularity is a score for how good a particular choice of communities is. Every node starts in its own community, and in each iteration two communities are joined if doing so yields the best improvement in modularity. This process continues until no improvements can be made.
The Louvain algorithm has two phases. The first phase is similar to greedy modularity optimization, but instead of merging communities it moves nodes between communities. The second phase makes a new graph where each community is a node and the edges between them represent edges between communities. On that graph the first phase is then applied.

## Demo 10

This demo shows algorithms for link prediction. You can see the Adamic-Adar index and the Jaccard coefficient.

The Jaccard coefficient uses the number of common neighbors to predict links.
The Adamic-Adar index also uses the number of common neighbors to predict links, but each neighbor's weight is inversely proportional to the neighbor's degree.
