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

## 📡 Real-time Telemetry Streaming

The **Device Service** now supports real-time data streaming using **SSE (Server-Sent Events)**. This allows the frontend to receive instant updates whenever an IoT device pushes new telemetry data, without needing to refresh the page or poll the API.

### 🏗️ The Broker Pattern
To handle multiple concurrent users watching the same or different devices, the service implements a **Broadcaster/Broker pattern**:

1.  **Publisher (IoT Device):** When a device sends data to `/api/iot/`, the `InsertTelemetry` service saves it to the database and simultaneously pushes the data into the **Broker's Notifier channel**.
2.  **Broker (Internal Service):** The Broker maintains a registry of all active client connections. It listens for new data and "broadcasts" a copy to every connected client's unique channel.
3.  **Subscriber (Frontend):** The SSE handler creates a persistent HTTP connection. It filters the broadcasted data by `device_id` and streams matching events to the client in real-time.

### 🗺️ Updated Telemetry Routes:

| Access | Method | Route | Description |
| :---: | :---: | :--- | :--- |
| **IoT Device** | `POST` | `/api/iot/` | Pushes new sensor data (Triggers Real-time Broadcast). |
| **Authorized** | `GET` | `/api/telemetry/:device_id` | Retrieves the latest single telemetry record from DB. |
| **Authorized** | `GET` | `/api/telemetry/:device_id/stream` | **Real-time SSE Stream:** Keeps a persistent connection open for live updates. |

### 🛠️ How to Consume the Stream (Frontend Example)
Since it uses the standard SSE protocol, you can easily consume the data using the native `EventSource` API in JavaScript:

```javascript
const eventSource = new EventSource('http://localhost:8081/api/telemetry/device-123/stream');

eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log("New Telemetry Received:", data);
    // Update your charts or dashboard UI here
};

eventSource.onerror = (err) => {
    console.error("SSE Connection Failed:", err);
};

```

### 1. Clone the Repository

```bash
git clone <YOUR_REPO_URL>
cd <YOUR_PROJECT_DIRECTORY_NAME>
