# Data Model

The system uses a **single-table DynamoDB design**.

---

## Table: Resources

### Partition Key (PK)

USER#<telegram_user_id>

### Sort Key (SK)

RES#<timestamp>#<uuid>

---

## Resource attributes

- title
- url
- type
- status
- tags (String Set)
- notes
- created_at

---

This design enables:
- Fast per-user queries
- Time-ordered results
- Easy future extensions
