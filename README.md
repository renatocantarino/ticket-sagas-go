# 🎟️ Ticket Purchase System with Sagas in Go

[![Go](https://img.shields.io/badge/Go-1.22+-007d9c?logo=go&logoColor=white)](https://golang.org)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org)
[![Clean Architecture](https://img.shields.io/badge/Pattern-Clean_Architecture-9c27b0)](#)
[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com)

A practical implementation of the **Saga Pattern** in Go using **Clean Architecture** to manage distributed transactions with **compensating actions**.

This project simulates a ticket purchase flow across multiple services, ensuring consistency even when failures occur.

---

## 🧩 Features

- Reserve a ticket
- Process payment
- Confirm purchase
- Send confirmation email
 Handle failures with **compensations**
- Persist domain events in SQLite    [WIP]
- Context-based logging with `sagaID` for tracing   [WIP]