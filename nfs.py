import sqlite3
import psutil
import time

# Database setup
db_connection = sqlite3.connect('nfs_metrics.db')
cursor = db_connection.cursor()

# Create table for NFS metrics
cursor.execute('''
    CREATE TABLE IF NOT EXISTS nfs_metrics (
        timestamp DATETIME,
        read_count INTEGER,
        write_count INTEGER,
        read_bytes INTEGER,
        write_bytes INTEGER
    )
''')
db_connection.commit()

def collect_nfs_metrics():
    # Collect NFS metrics. Here we're using psutil, but you might need to adjust
    # this depending on how your NFS is set up and what metrics you need.
    io_counters = psutil.net_io_counters(pernic=True)
    nfs_metrics = io_counters.get('nfs', None)

    if nfs_metrics:
        return {
            'read_count': nfs_metrics.packets_sent,
            'write_count': nfs_metrics.packets_recv,
            'read_bytes': nfs_metrics.bytes_sent,
            'write_bytes': nfs_metrics.bytes_recv
        }
    else:
        return None

def store_metrics(metrics):
    cursor.execute('''
        INSERT INTO nfs_metrics (timestamp, read_count, write_count, read_bytes, write_bytes)
        VALUES (CURRENT_TIMESTAMP, ?, ?, ?, ?)
    ''', (metrics['read_count'], metrics['write_count'], metrics['read_bytes'], metrics['write_bytes']))
    db_connection.commit()

while True:
    metrics = collect_nfs_metrics()
    if metrics:
        store_metrics(metrics)
    time.sleep(60)  # Collect and store metrics every 60 seconds

# Don't forget to close the database connection when done
db_connection.close()
