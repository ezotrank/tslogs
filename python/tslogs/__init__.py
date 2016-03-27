__version__ = '0.1'

import time
import threading
import socket
import json
from httplib2 import Http
import logging
from copy import deepcopy

log = logging.getLogger(__name__)

SETTINGS = {
    'tick': 30,
    'buff_size': 32,
    'handler': 'thread',
    'hostname': socket.gethostname(),
}

def init(host, port, **kwargs):
    SETTINGS['tick'] = kwargs.get('tick', SETTINGS['tick'])
    SETTINGS['buff_size'] = kwargs.get('buff_size', SETTINGS['buff_size'])
    SETTINGS['handler'] = kwargs.get('buff_size', SETTINGS['handler'])
    SETTINGS['hostname'] = kwargs.get('hostname', SETTINGS['hostname'])
    SETTINGS['host'] = host
    SETTINGS['port'] = port
    SETTINGS['next_send_time'] = now() + SETTINGS['tick']
    SETTINGS['buff'] = []

def send(key, val, **tags):
    try:
        split_time = str(now()).split('.')
        timestamp = int(split_time[0] + split_time[1][:3])
        if not 'host' in tags:
            tags['host'] = SETTINGS['hostname']
        SETTINGS['buff'].append({"metric": key, "value": val, "tags": tags, "timestamp": timestamp})
        if need_send():
            thread = threading.Thread(target=_send, args=(deepcopy(SETTINGS['buff'])))
            thread.start()
            SETTINGS['buff'] = []
            SETTINGS['buff_size'] = now() + SETTINGS['tick']
    except Exception as e:
        log.exception("can't send metrics, error: %s", str(e))

def _send(*buff):
    response, content = Http().request('http://%s:%d/api/put' % (SETTINGS['host'], SETTINGS['port']), 'POST', json.dumps(buff).encode('utf-8'))
    if int(response['status']) not in (200, 204):
        log.error("can't send metrics, code: %d", response['status'])
    SETTINGS['next_send_time'] = now() + SETTINGS['tick']

def need_send():
    if len(SETTINGS['buff']) > 0:
        return len(SETTINGS['buff']) > SETTINGS['buff_size'] or now() > SETTINGS['next_send_time']

def now():
    return time.time()
