#!/bin/bash
docker build . -t localhost:32000/informer
docker push localhost:32000/informer 