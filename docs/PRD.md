# Luna — AI Agent for Chartered Accountants

## What is Luna?

Luna is a CLI-based AI agent that automates the CA (Chartered Accountancy) workflow — from raw client data to filed returns. It does for accounting what Claude Code does for software engineering: understand the full context, automate the grunt work, let the professional focus on judgment calls.

Built on Anchor's graph architecture. Same tech (tree-sitter parsing, petgraph, MCP server), different domain.

## The Problem

CAs spend 60-70% of their time on mechanical tasks:
- Manual data entry from invoices/bank statements into Tally/journals
- Cross-checking figures across ledgers, trial balances, returns
- Filing repetitive forms on government portals (GST, TDS, ITR)
- Reconciliation — matching bank statements to books, finding mismatches

These are pattern-matching and data-flow problems. Exactly what a graph engine solves.

## Core Architecture

```
┌──────────────────────────────────────────────────┐
│                    Luna CLI                       │
│  luna ingest · luna validate · luna file · luna q │
├──────────────────────────────────────────────────┤
│                  Engine Layer                     │
│  ┌──────────┐ ┌───────────┐ ┌──────────────────┐│
│  │  Parser   │ │   Graph   │ │    Compute       ││
│  │ PDF/Excel │ │ Accounts  │ │ Tax/GST/TDS      ││
│  │ Tally XML │ │ Ledger    │ │ Validation       ││
│  │ Bank CSV  │ │ Relations │ │ Reconciliation   ││
│  └──────────┘ └───────────┘ └──────────────────┘│
├──────────────────────────────────────────────────┤
│               MCP Server (stdio)                  │
│  query · validate · reconcile · compute · file   │
├──────────────────────────────────────────────────┤
│              Portal Integration                   │
│  incometax.gov.in · GST portal · TDS portal      │
└──────────────────────────────────────────────────┘
```

### 1. Parser (Data Ingestion)

Extracts structured data from raw client documents:

| Source | Format | What's Extracted |
|--------|--------|-----------------|
| Invoices | PDF, image | Party, amount, GST, date, HSN/SAC |
| Bank statements | CSV, PDF, Excel | Date, narration, debit, credit, balance |
| Tally exports | XML, JSON | Ledger entries, vouchers, groups |
| Salary slips | PDF, Excel | Basic, HRA, deductions, TDS |
| Form 26AS / AIS | PDF | TDS credits, high-value transactions |

Tech: `pdf-extract` for PDFs, `calamine` for Excel, `serde` for CSV/XML/JSON. OCR via external service for scanned documents.

### 2. Graph (Regulatory Intelligence)

Petgraph-based graph where nodes are accounts/entries and edges are financial relationships:

**Node types:**
- `Account` — ledger account (Cash, Sales, GST Input, etc.)
- `Entry` — single journal entry (debit/credit pair)
- `Invoice` — source document
- `Party` — client, vendor, employee
- `Return` — GST return, ITR, TDS return
- `Rule` — tax rule/rate/threshold

**Edge types:**
- `Debits` / `Credits` — entry → account
- `BelongsTo` — account → group (Assets, Liabilities, etc.)
- `SourcedFrom` — entry → invoice/document
- `FiledIn` — entry → return
- `AppliesTo` — rule → account/transaction type
- `PaidBy` / `PaidTo` — party → entry

**What the graph enables:**
- Trace any amount back to its source document
- Find all entries that feed into a specific return line
- Detect circular references or double-counting
- Reconcile bank ↔ books by matching amounts + dates

### 3. Engine (Computation & Validation)

**Tax Computation:**
- GST: CGST/SGST/IGST calculation based on place of supply
- TDS: Section-wise deduction (194A, 194C, 194J, etc.) with threshold checks
- Income Tax: Slab computation, rebates, exemptions (Old vs New regime)
- Advance Tax: Quarterly liability estimation

**Validation Rules:**
- Double-entry balance check (total debits = total credits)
- GST input credit eligibility (blocked credits, reverse charge)
- TDS rate validation against section + party type
- PAN/GSTIN format validation
- Threshold checks (turnover limits, audit applicability)
- Mismatch detection: GSTR-2A vs books, 26AS vs books

**Reconciliation:**
- Bank reconciliation: match bank statement entries to book entries
- GST reconciliation: GSTR-1 vs GSTR-3B vs books
- TDS reconciliation: 26AS credits vs TDS deposited

### 4. Portal Integration

Automate filing on Indian government portals:

| Portal | What | Method |
|--------|------|--------|
| incometax.gov.in | ITR filing, 26AS download, advance tax challan | Headless browser / API |
| gst.gov.in | GSTR-1, GSTR-3B, GSTR-9 filing | GSP API (Goods and Services Tax Suvidha Provider) |
| TRACES | TDS return filing, Form 16/16A download | Headless browser |
| MCA | ROC filings, annual returns | Headless browser |

Phase 1 (MVP): JSON export in portal-compatible format. Manual upload.
Phase 2: Direct API integration via GSP for GST.
Phase 3: Full portal automation.

## CLI Commands

```bash
# Data ingestion
luna ingest <file_or_dir>          # Parse invoices, bank statements, Tally exports
luna ingest --watch <dir>          # Watch directory for new documents

# Querying (graph-powered)
luna query "GST liability for Jan 2026"
luna query "TDS on rent payments"
luna query "all entries for Party X"
luna trace <entry_id>              # Trace entry back to source document

# Validation
luna validate                      # Run all validation rules
luna validate --gst                # GST-specific checks
luna validate --tds                # TDS-specific checks
luna reconcile bank <statement>    # Bank reconciliation
luna reconcile gst                 # GSTR-2A vs books

# Computation
luna compute gst --period 2026-01  # GST liability for period
luna compute tds --period 2026-01  # TDS liability for period
luna compute tax --fy 2025-26      # Income tax computation
luna compute advance-tax --quarter Q4  # Advance tax estimate

# Filing
luna file gstr1 --period 2026-01   # Generate GSTR-1 JSON
luna file gstr3b --period 2026-01  # Generate GSTR-3B JSON
luna file itr --fy 2025-26         # Generate ITR JSON
luna file tds --quarter Q4-2026    # Generate TDS return

# System
luna build                         # Build/rebuild the financial graph
luna stats                         # Graph statistics
luna mcp                           # Start MCP server
luna init                          # Configure AI agent integration
```

## MCP Tools (for AI agents)

```
query     — Natural language query over the financial graph
validate  — Run validation rules, return violations
reconcile — Match two data sources, return mismatches
compute   — Calculate tax/GST/TDS for a period
trace     — Follow an amount from return → entry → document
```

## Data Model

```
Client/
├── FY2025-26/
│   ├── invoices/          # Raw PDFs, images
│   ├── bank/              # Bank statements
│   ├── tally/             # Tally exports
│   ├── salary/            # Salary data
│   ├── .luna/
│   │   ├── graph.bin      # Financial graph (petgraph + bincode)
│   │   ├── config.toml    # Client config (PAN, GSTIN, regime, etc.)
│   │   └── cache/         # Parsed document cache
│   ├── returns/
│   │   ├── gstr1/         # Generated GSTR-1 JSONs
│   │   ├── gstr3b/        # Generated GSTR-3B JSONs
│   │   ├── itr/           # Generated ITR JSONs
│   │   └── tds/           # Generated TDS returns
│   └── reports/
│       ├── pnl.json       # Profit & Loss
│       ├── bs.json        # Balance Sheet
│       └── reconciliation/
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go |
| Graph engine | Custom (port of Anchor's petgraph approach, or use Go graph lib) |
| PDF parsing | `pdfcpu` / `unidoc` |
| Excel parsing | `excelize` |
| CSV | stdlib `encoding/csv` |
| XML (Tally) | stdlib `encoding/xml` |
| CLI framework | `cobra` |
| MCP server | Go MCP SDK |
| Config | TOML (`pelletier/go-toml`) |
| Serialization | `encoding/gob` or Protocol Buffers |
| Portal automation | `chromedp` (headless Chrome) or direct API |

## MVP Scope (v0.1.0)

**In scope:**
1. Parse bank statements (CSV) and invoices (structured PDF)
2. Build financial graph (accounts, entries, parties)
3. GST computation (CGST/SGST/IGST)
4. Basic validation (double-entry balance, GST rate check)
5. Bank reconciliation (match by amount + date)
6. GSTR-1 and GSTR-3B JSON export
7. CLI commands: `ingest`, `build`, `query`, `validate`, `compute gst`, `file gstr1`
8. MCP server with `query` and `validate` tools

**Out of scope for MVP:**
- OCR / scanned document parsing
- Income tax computation (ITR)
- TDS computation
- Direct portal filing (API or headless browser)
- Multi-client management
- Tally integration
- Audit trail / revision history

## Competitive Landscape

| Tool | What it does | Gap |
|------|-------------|-----|
| Tally Prime | Manual bookkeeping + GST | No AI, no automation, desktop-only |
| Zoho Books | Cloud accounting | No AI agent, no CLI, SaaS lock-in |
| ClearTax | GST filing + ITR | Filing only, no bookkeeping, web-only |
| Suvit | AI data entry from Excel/PDF | Entry only, no computation/validation/filing |
| Legaltax AI | AI-assisted tax filing | Early stage, no graph, no reconciliation |

**Luna's edge:** End-to-end pipeline (ingest → graph → compute → validate → file) with graph intelligence. Not a SaaS — runs locally, CA owns their data. MCP server means any AI agent can use it.

## Success Metrics

1. **Time saved**: 60%+ reduction in manual data entry and cross-checking
2. **Accuracy**: Zero computation errors (validated against manual calculations)
3. **Coverage**: Handle 80% of a typical CA's monthly workflow (GST + bookkeeping)
4. **Adoption**: CA can go from raw documents to filed GSTR-1 in one session
