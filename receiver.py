import http.server
import socketserver
import os

PORT = 8080

class CSVRequestHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers["Content-Length"])
        file_data = self.rfile.read(content_length)

        # Save the received data to a file (e.g., received.csv)
        with open("received.csv", "wb") as f:
            f.write(file_data)

        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"CSV file received and saved successfully!")

if __name__ == "__main__":
    os.chdir(os.path.dirname(__file__))  # Set current directory to script location
    with socketserver.TCPServer(("", PORT), CSVRequestHandler) as httpd:
        print("Server listening on port", PORT)
        httpd.serve_forever()
