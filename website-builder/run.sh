
API_KEY="$1"
MODELS="$2"

if [ -z "$API_KEY" ]; then
    echo "Usage: $0 <api-key> <models>"
    exit 1
fi

if [ -z "$MODELS" ]; then
    echo "Usage: $0 <api-key> <models>"
    echo "Example: $0 your-api-key 'gemini-2.5-flash,gemini-2.0-flash'"
    exit 1
fi

export API_KEY
export MODELS
echo $MODELS
go run main.go web api -sse-write-timeout=10m a2a webui
 
