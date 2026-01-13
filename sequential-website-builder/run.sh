
API_KEY="$1"
MODELS="$1"
if [ -z "$API_KEY" ]; then
    echo "Usage: $0 <api-key> $1 <models>" 
    exit 1
fi

export API_KEY
export MODELS
go run main.go web api -sse-write-timeout=10m a2a webui
 
