# Integration Project 1: Mini Graph Database

## Overview
Build a minimal but functional graph database combining storage layer, graph structures, query execution, and transactions.

**Duration:** 40-50 hours
**Difficulty:** Very Hard

## Components Integration

### Phase 1: Storage Layer
- Buffer pool for page caching
- WAL for crash recovery
- Persistent storage

### Phase 2: Graph Structure
- CSR for efficient traversal
- Property storage for nodes/edges
- Index structures

### Phase 3: Query Engine
- Simple query language parser
- Query optimizer
- Execution engine

### Phase 4: Transactions
- Lock manager for concurrency
- MVCC for snapshot reads
- Transaction coordinator

## Query Language
```
CREATE NODE person {id: 1, name: "Alice", age: 30}
CREATE EDGE knows {from: 1, to: 2, since: 2020}

MATCH (p:person)-[:knows]->(friend)
WHERE p.name = "Alice"
RETURN friend.name, friend.age

BEGIN
CREATE NODE person {id: 3, name: "Bob"}
CREATE EDGE knows {from: 1, to: 3}
COMMIT
```

## Architecture
```
┌─────────────────────┐
│   Query Parser      │
├─────────────────────┤
│   Query Optimizer   │
├─────────────────────┤
│   Execution Engine  │
├─────────────────────┤
│   Transaction Mgr   │
├─────────────────────┤
│   Storage Layer     │
│ (Buffer + WAL)      │
└─────────────────────┘
```

## Features
- Node and edge CRUD operations
- Pattern matching queries
- Property filtering
- Transaction support (ACID)
- Crash recovery
- Concurrent access

## Test Scenarios
- End-to-end query execution
- Concurrent transactions
- Crash and recovery
- Performance benchmarks
- Social network dataset

## Performance Goals
- Simple query: <10ms
- Pattern match: <100ms (1M edges)
- Transaction commit: <50ms
- Crash recovery: <5s

## Time Estimate
Integration: 25-30 hours, Testing: 8-10 hours, Optimization: 5-8 hours
