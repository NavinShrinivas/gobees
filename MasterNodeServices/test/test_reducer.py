#!/usr/bin/env python3
import sys
import json

statefull_parent_node = "-1"
for line in sys.stdin:
    current_parent_node = line.strip().split(",")[0]
    current_outgoing_node = line.strip().split(",")[1]
    if statefull_parent_node == "-1":
        statefull_parent_node = current_parent_node
        print(statefull_parent_node,"\t","[",current_outgoing_node, sep="", end="")
    elif statefull_parent_node != current_parent_node : 
        print("]\n", sep = "", end="")
        statefull_parent_node = current_parent_node
        print(statefull_parent_node,"\t","[",current_outgoing_node, sep="", end="")
    else:
        print(", ",current_outgoing_node, sep = "", end="")
print("]", end="")

