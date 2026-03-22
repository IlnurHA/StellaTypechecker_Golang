#!/bin/sh
# Assuming you have the antlr JAR file in the same directory
java -Xmx500M -cp "./antlr-4.13.2-complete.jar:$CLASSPATH" org.antlr.v4.Tool -Dlanguage=Go -visitor -package parser *.g4