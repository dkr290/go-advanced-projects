import argparse
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


def safe_execute(command, *args, **kwargs):
    """Execute a Redis command with retry on connection loss."""
    global r
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            return getattr(r, command)(*args, **kwargs)
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
            r = connect_with_retry()
    raise RuntimeError("Redis command failed after maximum retries.")


def insert_records(client, num_records, sleep_ms):
    print(f"Inserting {num_records} records...")
    for i in range(num_records):
        key = f"record:{i}"
        value = f"value_{i}_{random.randint(1000, 9999)}"
        safe_execute("set", key, value)
        print(f"inserted {key} and {value}")
        time.sleep(sleep_ms / 1000)
        if (i + 1) % 500 == 0:
            print(f"  Inserted {i + 1}/{num_records}")
    print("Done inserting.")


def read_random_records(client, num_records, num_reads):
    if num_records <= 0:
        print("No records available for random reads (num_records <= 0).")
        return
    print(f"\nReading {num_reads} random records...")
    keys = [f"record:{random.randint(0, num_records - 1)}" for _ in range(num_reads)]
    for key in keys:
        value = safe_execute("get", key)
        print(f"  {key} -> {value}")
    print("\nDone.")


def parse_args():
    parser = argparse.ArgumentParser(description="Redis load test (insert/read).")
    parser.add_argument(
        "--insert",
        action="store_true",
        help="Run insert workload.",
    )
    parser.add_argument(
        "--read",
        action="store_true",
        help="Run read workload.",
    )
    parser.add_argument(
        "--num-records",
        type=int,
        default=NUM_RECORDS,
        help=f"Number of records for insert/random read range (default: {NUM_RECORDS}).",
    )
    parser.add_argument(
        "--num-reads",
        type=int,
        default=NUM_READS,
        help=f"Number of random reads (default: {NUM_READS}).",
    )
    parser.add_argument(
        "--sleep-ms",
        type=int,
        default=SLEEP_MS,
        help=f"Sleep between inserts in ms (default: {SLEEP_MS}).",
    )
    return parser.parse_args()


def main():
    args = parse_args()

    # If no mode is provided, keep previous behavior: run both.
    run_insert = args.insert or (not args.insert and not args.read)
    run_read = args.read or (not args.insert and not args.read)

    global r
    r = connect_with_retry()

    if run_insert:
        insert_records(r, args.num_records, args.sleep_ms)

    if run_read:
        read_random_records(r, args.num_records, args.num_reads)


if __name__ == "__main__":
    main()
