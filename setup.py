#!/usr/bin/env python3

from setuptools import setup

setup(
    name='config',
    description='config for soft-ap',
    author='Ricardo Mendoza <ricmm@canonical.com>',
    license='GPLv3',
    install_requires=[
        'pyyaml',
    ],
    scripts=[
        'config.py',
    ],
)
