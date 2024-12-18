# archive-api
Doodocs backend challenge

This project implements a REST API for managing archives and sending files via email. It follows clean code and architecture principles, ensures code extensibility. The API is built with [Gin](https://github.com/gin-gonic/gin) and includes testing for reliability.

---

## Features

1. **Archive Information**  
   Extracts detailed information from an uploaded archive file.
    - Supported format: `.zip`

2. **Archive Creation**  
   Combines multiple valid files into a `.zip` archive.
    - Supported file types:
        - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
        - `application/xml`
        - `image/jpeg`
        - `image/png`

3. **File Emailing**  
   Sends an uploaded file to multiple email addresses.
    - Supported file types:
        - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
        - `application/pdf`

---

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Install dependencies:

    ```bash
    go mod tidy    
    ```

3. Configure environment variables:
Create a .env file in the root directory with the following:
    ```
    SMTP_HOST=smtp.gmail.com
    SMTP_PORT=587
    SMTP_USER=<your-email@example.com>
    SMTP_PASSWORD=<your-email-password>
    ```

4. Run the server:
    ```bash
    go run main.go
    ```

## Endpoints and cURL Examples

### 1. Get Archive Information
- **URL**: `/api/archive/information`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`

**cURL Example**:
```
curl -X POST http://localhost:8080/api/archive/information \
    -H "Content-Type: multipart/form-data" \
    -F "file=@my_archive.zip" \
    -w "\nHTTP Status Code: %{http_code}\n"
```
Response:

```
{
    "filename": "my_archive.zip",
    "archive_size": 4102029.312,
    "total_size": 6836715.52,
    "total_files": 2,
    "files": [
        {
            "file_path": "photo.jpg",
            "size": 2516582.4,
            "mimetype": "image/jpeg"
        },
        {
            "file_path": "directory/document.docx",
            "size": 4320133.12,
            "mimetype": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        }
    ]
}
```
2. Create Archive
- **URL**: `/api/archive/files`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`

**cURL Example**:

```
curl -X POST http://localhost:8080/api/archive/files \
    -H "Content-Type: multipart/form-data" \
    -F "files[]=@document.docx" \
    -F "files[]=@avatar.png" \
    -w "\nHTTP Status Code: %{http_code}\n"
```
Response:

Returns the .zip archive as binary data.

To save the output as a .zip file:

```
curl -X POST http://localhost:8080/api/archive/files \
    -H "Content-Type: multipart/form-data" \
    -F "files[]=@document.docx" \
    -F "files[]=@avatar.png" \
    -w "\nHTTP Status Code: %{http_code}\n" \
    --output your_archive.zip
```
3. Send File via Email
- **URL**: `/api/mail/file`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
cURL Example:
```
curl -X POST http://localhost:8080/api/mail/file \
    -H "Content-Type: multipart/form-data" \
    -F "file=@document.docx" \
    -F "emails=elonmusk@x.com,jeffbezos@amazon.com,zuckerberg@meta.com" \
    -w "\nHTTP Status Code: %{http_code}\n"
```