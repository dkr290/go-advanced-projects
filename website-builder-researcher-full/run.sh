
API_KEY="$1"
if [ -z "$API_KEY" ]; then
    echo "Usage: $0 <api-key>"
    exit 1
fi

export API_KEY
go run main.go web api -sse-write-timeout=10m webui
 
