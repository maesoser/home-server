FROM debian

RUN apt update &&\
  apt upgrade -y &&\
  apt install -y python3-pip python3-pandas python3-numpy \
  python3-scipy python3-matplotlib python3-seaborn \
  libatlas-base-dev python-dev gfortran pkg-config libfreetype6-dev &&\
  pip3 install --no-cache jupyterlab zipline pyfolio empyrical alphalens


RUN mkdir /jupyter
COPY jupyter_base_config.py /jupyter/jupyter_base_config.py
COPY run.sh /jupyter/run.sh
ENTRYPOINT ["/usr/bin/bash", "/jupyter/run.sh"]

