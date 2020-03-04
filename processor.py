#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = 'imwuxuxing@gmail.com'

import os, commands, time, sys, logging, json
import yaml
import click
import docker
from redis import Redis


confFile = "%s/%s" %(os.getcwd(), "conf/appconf.yaml")
pushImgScript = "%s/%s" %(os.getcwd(), "script/img.sh")
logging.basicConfig(level=logging.INFO, filename="/www/dnt-swr/logs/dnt-swr-processor.info.log")
rdq = "REDISQUEUE"

def AppConf():
    confFile = "%s/%s" %(os.getcwd(), "conf/appconf.yaml")
    confRes = None
    if os.path.exists(confFile):
        f = open(confFile, 'r')
        confRes = yaml.load(f, Loader=yaml.FullLoader)
    return confRes

def DockerLogin(username=None, password=None, baseUrl=None):
    client = docker.from_env()
    try:
        loginRes = client.login(username=username, password=password, registry=registry)
    except Exception,err:
        print err

class Opharbor(object):
    def __init__(self):
        pass
    def SyncToSwr(self, harborUri, swrUri, addr):
        if addr in AppConf()['Builder']:
            imgCommand = "/bin/bash %s %s %s %s" %(pushImgScript, addr, harborUri, swrUri)
            (status, output) = commands.getstatusoutput(imgCommand)
            #print(status, output)
            return {"status":status, "res":output, "other":"push from builder", "processor": "processor id %s"%(os.getpid()), "time": time.strftime("%d/%b/%Y %H:%M:%S", time.localtime())}
        else:
            imgCommand = "/bin/bash %s %s %s %s" %(pushImgScript, 'local', harborUri, swrUri)
            (status, output) = commands.getstatusoutput(imgCommand)
            #print(status, output)
            return {"status":status, "msg":output, "other": "push from local", "processor": "processor id %s"%(os.getpid()), "time": time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())}

    def DoAction(self, redis):
        try:
            keys = redis.keys()
            length = len(keys)
            if length != 0:
                data = redis.lpop(rdq)
                jdata = json.loads(data)
                harborUri = jdata["harborurl"]
                swrUri = jdata["swrurl"]
                addr = jdata["builder"]
                syncRes = self.SyncToSwr(harborUri, swrUri, addr)
                logging.info(syncRes)
            else:
                time.sleep(1)
        except Exception, e:
            print e
            sys.exit(1)

@click.command()
@click.option('-h', "--host", default="127.0.0.1", help="queue server's host")    
@click.option('-p', "--port", default="6379", help="queue server's port")
def run(**options):
    hhost = options['host']
    pport = options['port']
    print("%s: processor start, Queue Server info: -h %s -p %s") %(time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()), hhost, pport)
    redis = Redis(host=hhost, port=pport, db=0)
    opharbor = Opharbor()
    while True:
        opharbor.DoAction(redis)
    
if __name__ == '__main__':
    run()
