from os import urandom
from hashlib import sha256
from random import expovariate
import time
import sys

if len(sys.argv) > 1:
    rate = float(sys.argv[1])
else:
    rate = 1.0          # default rate: 1 Hz

if len(sys.argv) > 2:
    max_events = int(sys.argv[2])
else:
    max_events = None

event_count = 0
while max_events is None or event_count < max_events:
    event_count += 1
    print("%s %s" % (time.time(), sha256(urandom(20)).hexdigest()))
    time.sleep(expovariate(rate))
