FROM debian

RUN apt update && apt upgrade -y && apt install -y python3-pip python3-pandas python3-numpy python3-scipy python3-matplotlib python3-seaborn && pip3 install --no-cache jupyterlab

RUN mkdir /jupyter
 
ENTRYPOINT ["/usr/local/bin/jupyter", "notebook", "--config=/jupyter/config.py"]

