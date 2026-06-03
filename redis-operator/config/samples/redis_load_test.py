import random
import time

import redis

REDIS_HOST = "localhost"
REDIS_PORT = 6379
REDIS_PASSWORD = "yourpassword"
NUM_RECORDS = 10000
NUM_READS = 200
SLEEP_MS = 40
RETRY_SLEEP_SEC = 5
MAX_RETRIES = 20


def connect_with_retry():
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            client = redis.Redis(
                host=REDIS_HOST,
                port=REDIS_PORT,
                password=REDIS_PASSWORD,
                decode_responses=True,
                socket_connect_timeout=5,
            )
            client.ping()
            print(f"Connected to Redis on attempt {attempt}.")
            return client
        except (
            redis.ConnectionError,
            redis.AuthenticationError,
            redis.TimeoutError,
        ) as e:
            print(
                f"Attempt {attempt}/{MAX_RETRIES} failed: {e}. Retrying in {RETRY_SLEEP_SEC}s..."
            )
            time.sleep(RETRY_SLEEP_SEC)
    raise RuntimeError("Could not connect to Redis after maximum retries.")


def safe_execute(fn, *args, **kwargs):
    """Execute a Redis command with retry on connection loss."""
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            return fn(*args, **kwargs)
        except redis.ReadOnlyError as e:
            print(
                f"Redis is read-only (attempt {attempt}/{MAX_RETRIES}): {e}. "
                f"Waiting for promotion to master in {RETRY_SLEEP_SEC}s..."
            )
            time.sleep(RETRY_SLEEP_SEC)
        except (redis.ConnectionError, redis.TimeoutError) as e:
            print(
                f"Redis command failed (attempt {attempt}/{MAX_RETRIES}): {e}. Retrying in {RETRY_SLEEP_SEC}s..."
            )
            time.sleep(RETRY_SLEEP_SEC)
            global r
            r = connect_with_retry()
    raise RuntimeError("Redis command failed after maximum retries.")



r = connect_with_retry()
# Insert 10000 records
print(f"Inserting {NUM_RECORDS} records...")
for i in range(NUM_RECORDS):
    key = f"record:{i}"
    value = f"value_{i}_{random.randint(1000, 9999)}"
    safe_execute(r.set, key, value)
    print(f"inserted {key} and {value}")
    time.sleep(SLEEP_MS / 1000)
    if (i + 1) % 500 == 0:
        print(f"  Inserted {i + 1}/{NUM_RECORDS}")

print("Done inserting.")

# Read 200 random records
print(f"\nReading {NUM_READS} random records...")
keys = [f"record:{random.randint(0, NUM_RECORDS - 1)}" for _ in range(NUM_READS)]
for key in keys:
    value = safe_execute(r.get, key)
    print(f"  {key} -> {value}")

print("\nDone.")
