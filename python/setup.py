import re
import os.path
from setuptools import setup

HERE = os.path.abspath(os.path.dirname(__file__))
README_PATH = os.path.join(HERE, 'README.md')

try:
    with open(README_PATH) as fd:
        README = fd.read()
except IOError:
    README = ''

INIT_PATH = os.path.join(HERE, 'tslogs/__init__.py')
with open(INIT_PATH) as fd:
    INIT_DATA = fd.read()
    VERSION = re.search(r"^__version__ = ['\"]([^'\"]+)['\"]", INIT_DATA, re.MULTILINE).group(1)

setup(
    name='tslogs',
    packages=['tslogs'],
    version=VERSION,
    author='Maxim Kremenev',
    author_email='ezo@kremenev.com',
    classifiers=[
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.2",
        "Programming Language :: Python :: 3.3",
        "Programming Language :: Python :: 3.4",
        "Programming Language :: Python :: 3.5",
        ],
    install_requires=[],
    tests_require=[],
    )
