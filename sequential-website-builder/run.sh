
API_KEY="$1"
MODELS="$2"
DEBUG_FLAG=$3
if [ -z "$API_KEY" ]; then
    echo "Usage: $0 <api-key>  <models> [debugflag]" 
    exit 1
fi

export API_KEY
export MODELS
if [ -n "$DEBUG_FLAG" ]; then 
    export DEBUG="$DEBUG_FLAG"
fi
go run main.go web api -sse-write-timeout=10m a2a webui
 
