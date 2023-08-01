# Proxy Checker

# Go 
-  go run ./cmd/server .

# Make
- make compile
- make install
- make run
  
# Docker
- docker build -t proxy-checker .
- docker run --restart always -p 8080:8080 -d proxy-checker
