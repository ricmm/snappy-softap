#!/usr/bin/env python3

import os
import os.path
import subprocess
import sys
import yaml

CONFIG_FILENAME = 'hostapd.conf'

creds = {"ssid": "", "passphrase": ""}

def set_config(config_file, config_yaml):
    """save config in a shell sourceable config_file"""
    global creds
    with open(config_file, 'r+') as f:
        try:
            creds = config_yaml['config'][os.environ['SNAP_NAME']]
        except (KeyError, TypeError):
            default = {'config': {
                          os.environ['SNAP_NAME']: creds } }
            print("You need to give a yaml file like: {}".format(default))
        save_shell_source(f, creds)


def save_shell_source(config_fd, creds):
    """Save in fd from creds dict"""
    config = config_fd.readlines()
    config_fd.seek(0)
    config_fd.truncate()
    for (key, value) in creds.items():
        for line,n in enumerate(config):
            print(config[line])
            if config[line].startswith(key):
                config[line] = "{}={}\n".format(key, value)

    for i in config:
        config_fd.write(i)


if __name__ == '__main__':
    config_file = os.path.join(os.environ['SNAP_APP_DATA_PATH'], CONFIG_FILENAME)

    config_yaml = yaml.load(sys.stdin)
    if config_yaml:
        set_config(config_file, config_yaml)
