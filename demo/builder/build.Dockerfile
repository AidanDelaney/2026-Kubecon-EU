FROM nvidia/cuda:13.1.1-cudnn-devel-ubuntu24.04

ARG cnb_uid=1001
ARG cnb_gid=1001
ENV CNB_USER_ID=${cnb_uid} \
    CNB_GROUP_ID=${cnb_gid}

RUN apt-get update && apt-get install -y --no-install-recommends \
      bash \
      ca-certificates \
      git \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid ${cnb_gid} cnb \
 && useradd --uid ${cnb_uid} --gid ${cnb_gid} -m -s /bin/bash cnb


USER ${cnb_uid}:${cnb_gid}
