# Graph Report - pkg/module/payroll  (2026-07-18)

## Corpus Check
- 2 files · ~769 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 30 nodes · 43 edges · 6 communities (3 shown, 3 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bf90c990`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]

## God Nodes (most connected - your core abstractions)
1. `Repository` - 8 edges
2. `Service` - 7 edges
3. `Context` - 5 edges
4. `Context` - 4 edges
5. `NewRepository()` - 3 edges
6. `EmployeeSalaryDto` - 3 edges
7. `NewService()` - 3 edges
8. `DB` - 2 edges
9. `Repository` - 2 edges
10. `EmployeeSalaryDto` - 2 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Import Cycles
- None detected.

## Communities (6 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.29
Nodes (5): SavePayrollRequest, Context, EmployeeSalaryDto, EmployeeSalaryFindAllRequest, ResultPagination

### Community 1 - "Community 1"
Cohesion: 0.38
Nodes (5): Service, Repository, SimulatePayrollResultDto, NewService(), SimulatePayrollRequest

### Community 3 - "Community 3"
Cohesion: 0.83
Nodes (3): DB, Repository, NewRepository()

## Knowledge Gaps
- **9 isolated node(s):** `Time`, `SimulatePayrollResultDto`, `EmployeeSalaryFindAllRequest`, `ResultPagination`, `SimulatePayrollRequest` (+4 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Repository` connect `Community 3` to `Community 2`, `Community 4`, `Community 5`?**
  _High betweenness centrality (0.125) - this node is a cross-community bridge._
- **Why does `Service` connect `Community 1` to `Community 0`?**
  _High betweenness centrality (0.125) - this node is a cross-community bridge._
- **Why does `Context` connect `Community 2` to `Community 4`, `Community 5`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **What connects `Time`, `SimulatePayrollResultDto`, `EmployeeSalaryFindAllRequest` to the rest of the system?**
  _9 weakly-connected nodes found - possible documentation gaps or missing edges._