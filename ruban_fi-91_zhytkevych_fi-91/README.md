# DDDB

### Variant

Variant №8 - **Collection of text documents with full-text search**

Documents database with inverted index.


## Team
1. Denis Ruban FI-91
2. Ivan Zhytkevych FI-91

## Schedule

| Era        | TBD | Deadline | Complexity |
| -------    | ------------ | -------- | ---------- |
| Stagnation | Writing Storage, Indexer and Domain. | 04.10.2021 | `***` |
| AllDone, Johny | Writing Parser and UI. Linking components. | 01.11.2021 | `*` |
| `O(n!)`   | Evaluating algorithms complexity. Optimizing stage. | 29.11.2021 | `***` |


## Architecture

![architecture diagram](./architecture.png)


# Config

You can create `dddb-conf.json` in the directory you execute your binary in order to customize paths where **Indexer** and **Collection Storage** save their files. The DB considers current working directory as directory for **Indexer** and **Storage** by default if the config is not present.

| Field | Explanation | Type |
| ---   | ---         | ---- |
| `storage_path` | Working path for storage | String path |
| `storage_type` | Type of storage to use (currently only filesystem `fs` type is supported) | String |
| `indexer_path` | Working path for indexer | String path |
| `indexer_type` | Type of indexer to use (currently only filesystem `fs` type is supported) | String |

Config `dddb-conf.json` example:
```json
{
	"storage_path": "/tmp/dddb/storage",
	"storage_type": "fs",
	"indexer_path": "/tmp/dddb/indexer",
	"indexer_type": "fs"
}
```
