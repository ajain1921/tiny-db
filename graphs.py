# csv format:
# timestamp, sent_timestamp, bytes_received

import pandas as pd
import matplotlib.pyplot as plt

DATA_FILE_NAME = 'data copy.csv'

df = pd.read_csv(DATA_FILE_NAME)

# First turn timestamp and sent_timestamp to float for computation
df['timestamp'] = df['timestamp'].apply(
    lambda timestamp: float(timestamp))
df['sent_timestamp'] = df['sent_timestamp'].apply(
    lambda timestamp: float(timestamp))
# Rename sent column to delay and subtract
df = df.rename(columns={'sent_timestamp': 'delay'})


def mapper(row):
    return row[0] - row[1]


df['delay'] = df.apply(mapper, axis=1)

# Turn timestamp to int for bucketing
df['timestamp'] = df['timestamp'].apply(
    lambda timestamp: int(timestamp))

df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s').dt.time

# Get delay frame by filtering for timestamp and delay columns
delay_frame = df.filter(items=['timestamp', 'delay'])
min_frame = delay_frame.groupby(['timestamp']).min()
max_frame = delay_frame.groupby(['timestamp']).max()
median_frame = delay_frame.groupby(['timestamp']).median()
ninetieth_frame = delay_frame.groupby(['timestamp']).quantile(0.9)

ax = min_frame.plot()
max_frame.plot(ax=ax)
median_frame.plot(ax=ax)
ninetieth_frame.plot(ax=ax)

plt.savefig('delay_graph.png', format='png')
plt.show()

# Get bandwidth frame by filtering for timestamp and bytes_received columns
bandwidth_frame = df.filter(items=['timestamp', 'bytes_received'])
avg_frame = bandwidth_frame.groupby(['timestamp']).mean()

avg_frame.plot()

plt.savefig('bandwidth_graph.png', format='png')
plt.show()
