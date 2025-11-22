# Integration Project 2: Social Network Analyzer

## Overview
Build a social network analysis tool using graph algorithms, property storage, and interactive queries.

**Duration:** 30-40 hours
**Difficulty:** Hard

## Features
- Load graph from CSV (users, friendships)
- Friend recommendations (2-hop queries)
- Community detection (connected components)
- Influence analysis (PageRank)
- Interactive CLI
- Query statistics and visualization

## Algorithms
- BFS for shortest paths
- PageRank for influence
- Triangle counting for clustering
- Connected components for communities

## Dataset Format
```csv
users.csv: id,name,age,city
edges.csv: from,to,timestamp
```

## CLI Commands
```
> load users.csv edges.csv
> recommend friends 1
> pagerank --iterations 10
> community detect
> stats
> export results.json
```

## Performance Goals
- Load 1M users in <30s
- Friend recommendations: <100ms
- PageRank on 1M nodes: <10s
- Interactive response: <1s

## Time Estimate
Core: 20-25 hours, Testing: 5-7 hours, Polish: 5-8 hours
