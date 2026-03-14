#!/bin/sh
# Assuming you have the antlr JAR file in the same directory
alias antlr4='java -Xmx500M -cp "./antlr-4.13.2-complete.jar:$CLASSPATH" org.antlr.v4.Tool'
antlr4 -Dlanguage=Go -visitor -package parser *.g4