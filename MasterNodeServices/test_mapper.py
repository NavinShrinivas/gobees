#!/usr/bin/env python3

import sys


for line in sys.stdin:
    if line[0] == "#":
        pass
    else : 
        # '#' only occur in starting of file
        print(line.replace("\t"," "), end="")
        for i in sys.stdin:
            print(i.replace("\t", ' '), end="")
