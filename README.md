# 🚀 Device Service V2 & Auth Service V2

This project implements a modular and isolated **Microservice** architecture, developed using the **Go** language and the **Gin** web framework. The system consists of two main independent services: the **Auth Service** (Identity Management) and the **Device Service** (Device Management).

## 🛠️ System Requirements

* Go **1.21+**
* Docker & Docker Compose (for running databases and service orchestration)

## 🏗️ Project Structure (Layered Architecture)

Each Go service in this project follows a strict layered structure to maintain clear separation of concerns. 

| Directory | Description |
| :--- | :--- |
| `cmd/` | Contains the main entry point (`main.go`) for each application. |
| `internal/domain` | **Domain Models** (`models.User`, `models.Device`, etc.). Pure data structures for business logic. |
| `internal/service` | **Service Layer** (Business Logic). Implements use cases and interacts with the Repository. |
| `internal/repository` | **Repository Layer**. Manages data interaction (GORM/Database). Contains Repository Models (DB Schema). |
| `internal/http/handler` | **Handler/Controller**. Receives HTTP requests, performs binding, and calls the Service Layer. |
| `internal/http/middleware` | Common middleware functions, such as JWT authentication and authorization. |

---

## 🔑 Core Service Components

### 1. Auth Service (`auth-service-v2`)

This service is solely responsible for user identity management, authentication, and credential updates.

#### 🗺️ Routes (Endpoints):

| Access | Method | Route | Description |
| :---: | :---: | :--- | :--- |
| **Public** | `POST` | `/register` | Identity registration. |
| **Public** | `POST` | `/login` | Authentication and Token Issuance. |
| **Authenticated**| `GET` | `/api/me` | Retrieves the current user's identity data. |
| **Authenticated**| `PUT` | `/api/me/password` | Modifies user credentials (password). |
| **Authenticated**| `PUT` | `/api/me/name` | Modifies user identity data (name). |
| **General** | `GET` | `/api/health`| Service health status check. |

### 2. 📱 Device Service (`device-service-v2`)

This service manages all device and telemetry data, including CRUD operations requiring owner authorization.

#### 🗺️ Routes (Endpoints):

| Access | Method | Route | Description |
| :---: | :---: | :--- | :--- |
| **Authorized**| `POST` | `/api/devices/` | Creates a new Device resource. |
| **Authorized**| `GET` | `/api/devices/` | Retrieves a list of Device resources (by owner). |
| **Authorized**| `GET` | `/api/devices/:id` | Retrieves specific details for a single Device. |
| **Authorized**| `PUT` | `/api/devices/:id` | Updates an existing Device resource. |
| **Authorized**| `DELETE` | `/api/devices/:id` | Deletes a Device resource. |

---

## 🔗 Authentication and Authorization Mechanism

Both services are integrated using the **JWT (JSON Web Token)** standard to secure access to `/api/` routes.

1.  The **Auth Service** issues a JWT upon successful `/login`.
2.  This token must be included in the `Authorization: Bearer <token>` header to access all protected routes on both services.
3.  Dedicated middleware (`AuthRequired` / `DeviceRequired`) validates the token's integrity and the user's authorization before allowing access to the handlers.

---

## ⚙️ Running the Services (Development Guide)

### 1. Clone the Repository

```bash
git clone <YOUR_REPO_URL>
cd <YOUR_PROJECT_DIRECTORY_NAME>
