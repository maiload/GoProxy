# GoProxy

---

## Generating TLS Keys and Certificates with OpenSSL

### **1. Install OpenSSL**

- **macOS**: Install via Homebrew:
  ```bash
  brew install openssl
  ```
- **Window**: Download and install OpenSSL from [here](https://slproweb.com/products/Win32OpenSSL.html).

### **2. Generate Private Key (key.pem)**

  ```bash
  openssl genpkey -algorithm RSA -out key.pem
  ```

### **3. Generate Certificate Signing Request (CSR)**

  ```bash
  openssl req -new -key key.pem -out request.csr
  ```

### **4. Generate Self-Signed Certificate (cert.pem)**

  ```bash
  openssl x509 -req -in request.csr -signkey key.pem -out cert.pem
  ```

---

## How to Write a Config.yml

The HTTP server is executed based on the presence of the `server.port` value, while the HTTPS server is executed based on the presence of the three values in `server.ssl`

---

## How to Run with Docker Compose

### Requirements

- Docker
- Docker Compose
- cert.pem & key.pem

### 1. Move SSL Certificates

  ```bash
  mkdir ssl
  mv cert.pem key.pem ssl/
  ```

### 2. Move to the Docker Directory

  ```bash
  cd docker
  ```

### 3. Start Containers with Docker Compose

  ```bash
  docker-compose up -d
  ```

### 4. Check the Status of Running Containers

  ```bash
  docker-compose ps
  ```

### 5. Stop and Remove Containers

  ```bash
  docker-compose down
  ```