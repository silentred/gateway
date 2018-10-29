#!/usr/local/bin/python
from base64 import b64encode
from hashlib import md5
from random import randint, choice
from string import ascii_letters, digits
from datetime import datetime
from locust import HttpLocust, TaskSet, task


def random_str(num):
    """ return random string """
    return ''.join(choice(ascii_letters + digits) for i in range(num))

def unix_timestamp():
    return int((datetime.utcnow() -  datetime(1970,1,1)).total_seconds())

def get_sign(**kwargs):
    secret = "test123"
    lines = "%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s" % (kwargs["key"], kwargs["method"],
    kwargs["type"], kwargs["path"], kwargs["query"], kwargs["timestamp"], kwargs["random"], secret)
    md5_obj = md5()
    md5_obj.update(lines)
    sign = b64encode(bytes(md5_obj.hexdigest()))
    return sign


class HelloTestSet(TaskSet):
    """ test set """
    min_wait = 100
    max_wait = 500
    host = 'localhost:8088'
    methods = ['GET', 'POST']
    services = ['hello', 'world', 'jason']

    @task
    def hello(self):
        svc = self.services[randint(0, len(self.services) - 1)]
        path = "/%s/api-%d" % (svc, randint(1, 5))
        method = self.methods[randint(0, len(self.methods) - 1)]
        app_key = "ios-4.1"
        con_type = "application/json"
        ts = unix_timestamp()
        nonce = random_str(16)
        sign = get_sign(key=app_key, method=method, type=con_type, query="", path=path, timestamp=ts, random=nonce)
        headers = {
            "Host": "%s.luoji.com" % svc,
            "Content-Type": con_type,
            "X-Sign": sign,
            "X-App-Key": app_key,
            "X-Nonce": nonce,
            "X-Timestamp": "%d" % ts,
        }
        self.client.request(method, path, headers=headers)

    def on_start(self):
        pass

class EntreeTest(HttpLocust):
    """ for entree gateway """
    task_set = HelloTestSet

#print unix_timestamp()