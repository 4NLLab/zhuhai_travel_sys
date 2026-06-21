#!/usr/bin/env python3
import http.server
import os
import socketserver
import urllib.error
import urllib.request


ROOT = "/www/wwwroot/backend"
API_TARGET = "http://127.0.0.1:8080"


class Handler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        if not self.path.startswith("/api/v1/") and self.path.split("?", 1)[0].endswith((".html", "/")):
            self.send_header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
            self.send_header("Pragma", "no-cache")
            self.send_header("Expires", "0")
        super().end_headers()

    def do_GET(self):
        if self.path.startswith("/api/v1/"):
            self.proxy()
            return
        if self.path == "/":
            self.path = "/index.html"
        super().do_GET()

    def do_POST(self):
        if self.path.startswith("/api/v1/"):
            self.proxy()
            return
        self.send_error(404)

    def do_PUT(self):
        if self.path.startswith("/api/v1/"):
            self.proxy()
            return
        self.send_error(404)

    def do_DELETE(self):
        if self.path.startswith("/api/v1/"):
            self.proxy()
            return
        self.send_error(404)

    def proxy(self):
        length = int(self.headers.get("Content-Length", "0") or "0")
        body = self.rfile.read(length) if length else None
        headers = {
            key: value
            for key, value in self.headers.items()
            if key.lower() not in ("host", "content-length", "connection")
        }
        req = urllib.request.Request(
            API_TARGET + self.path,
            data=body,
            headers=headers,
            method=self.command,
        )
        try:
            with urllib.request.urlopen(req, timeout=30) as resp:
                self.send_response(resp.status)
                for key, value in resp.headers.items():
                    if key.lower() not in ("transfer-encoding", "connection"):
                        self.send_header(key, value)
                self.end_headers()
                self.wfile.write(resp.read())
        except urllib.error.HTTPError as err:
            self.send_response(err.code)
            for key, value in err.headers.items():
                if key.lower() not in ("transfer-encoding", "connection"):
                    self.send_header(key, value)
            self.end_headers()
            self.wfile.write(err.read())
        except Exception as err:
            self.send_error(502, str(err))


class ThreadingHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    daemon_threads = True
    allow_reuse_address = True


if __name__ == "__main__":
    os.chdir(ROOT)
    server = ThreadingHTTPServer(("0.0.0.0", 80), Handler)
    server.serve_forever()
