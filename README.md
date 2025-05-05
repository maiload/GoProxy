# GoProxy

---

### Generating TLS Keys and Certificates with OpenSSL

#### **1. Install OpenSSL**

- **macOS**: Install via Homebrew:
  ```bash
  brew install openssl
  ```
- **Window**: Download and install OpenSSL from [here](https://slproweb.com/products/Win32OpenSSL.html).

#### **2. Generate Private Key (key.pem)**

  ```bash
  openssl genpkey -algorithm RSA -out key.pem -aes256
  ```

#### **3. Generate Certificate Signing Request (CSR)**

  ```bash
  openssl req -new -key key.pem -out request.csr
  ```

#### **4. Generate Self-Signed Certificate (cert.pem)**

  ```bash
  openssl x509 -req -in request.csr -signkey key.pem -out cert.pem
  ```
